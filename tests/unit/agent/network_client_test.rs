#[cfg(test)]
mod tests {
    use super::*;
    use std::sync::Arc;
    use tokio::sync::Mutex;
    use std::collections::HashMap;
    use chrono::Utc;
    use mockito::{mock, Matcher};

    // Mock implementations for testing
    struct MockEventStore {
        events: Arc<Mutex<Vec<TelemetryEvent>>>,
    }

    impl MockEventStore {
        fn new() -> Self {
            Self {
                events: Arc::new(Mutex::new(Vec::new())),
            }
        }

        async fn add_event(&self, event: TelemetryEvent) {
            let mut events = self.events.lock().await;
            events.push(event);
        }

        async fn get_pending_events(&self) -> Vec<TelemetryEvent> {
            let events = self.events.lock().await;
            events.clone()
        }

        async fn mark_events_sent(&self, event_ids: Vec<String>) {
            let mut events = self.events.lock().await;
            events.retain(|e| !event_ids.contains(&e.id));
        }

        async fn clear(&self) {
            let mut events = self.events.lock().await;
            events.clear();
        }
    }

    // Test data factories
    fn create_test_config() -> Config {
        Config {
            server_url: "https://api.insec.com".to_string(),
            agent_id: "test-agent-123".to_string(),
            tenant_id: "test-tenant".to_string(),
            collection_interval: 30,
            max_batch_size: 10,
            tls_ca_cert: None,
            tls_client_cert: None,
            tls_client_key: None,
            log_level: "info".to_string(),
            enable_compression: false,
            retry_attempts: 3,
            retry_delay: 1,
            heartbeat_interval: 60,
        }
    }

    fn create_test_events(count: usize) -> Vec<TelemetryEvent> {
        (0..count)
            .map(|i| TelemetryEvent {
                id: format!("test-event-{}", i),
                timestamp: Utc::now(),
                event_type: EventType::Process,
                data: {
                    let mut data = HashMap::new();
                    data.insert("process_name".to_string(), serde_json::Value::String(format!("test{}.exe", i)));
                    data.insert("pid".to_string(), serde_json::Value::Number(i as i64));
                    data
                },
                metadata: {
                    let mut metadata = HashMap::new();
                    metadata.insert("risk_score".to_string(), serde_json::Value::Number(0.5.into()));
                    metadata
                },
            })
            .collect()
    }

    #[tokio::test]
    async fn test_http_client_initialization() {
        let config = create_test_config();
        let client = HttpClient::new(config.clone());

        assert_eq!(client.config.server_url, config.server_url);
        assert_eq!(client.config.agent_id, config.agent_id);
    }

    #[tokio::test]
    async fn test_successful_event_batch_send() {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(3);

        // Mock successful server response
        let _mock = mock("POST", "/api/v1/events")
            .match_header("Authorization", "Bearer test-token")
            .match_header("Content-Type", "application/json")
            .match_header("X-Agent-ID", "test-agent-123")
            .match_header("X-Tenant-ID", "test-tenant")
            .with_status(200)
            .with_body(r#"{"status": "success", "events_processed": 3}"#)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_ok());

        let response = result.unwrap();
        assert_eq!(response.status, 200);
    }

    #[tokio::test]
    async fn test_event_batch_send_server_error() {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(2);

        // Mock server error response
        let _mock = mock("POST", "/api/v1/events")
            .with_status(500)
            .with_body(r#"{"error": "Internal server error"}"#)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_err());

        let error = result.unwrap_err();
        assert!(error.to_string().contains("500"));
    }

    #[tokio::test]
    async fn test_event_batch_send_network_error() {
        let mut config = create_test_config();
        config.server_url = "https://nonexistent.invalid.server".to_string();
        let client = HttpClient::new(config);
        let events = create_test_events(1);

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_event_batch_send_timeout() {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(1);

        // Mock timeout response
        let _mock = mock("POST", "/api/v1/events")
            .with_status(408)
            .with_body(r#"{"error": "Request timeout"}"#)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_compressed_event_batch_send() {
        let mut config = create_test_config();
        config.enable_compression = true;
        let client = HttpClient::new(config);
        let events = create_test_events(5);

        // Mock successful compressed response
        let _mock = mock("POST", "/api/v1/events")
            .match_header("Content-Encoding", "gzip")
            .with_status(200)
            .with_body(r#"{"status": "success", "events_processed": 5}"#)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_large_event_batch_chunking() {
        let mut config = create_test_config();
        config.max_batch_size = 2;
        let client = HttpClient::new(config);
        let events = create_test_events(5); // More than max_batch_size

        let mut call_count = 0;
        let _mock = mock("POST", "/api/v1/events")
            .with_status(200)
            .with_body(move |_| {
                call_count += 1;
                r#"{"status": "success"}"#.to_string()
            })
            .expect(3) // Should be called 3 times (5 events / 2 batch size, rounded up)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_retry_mechanism_success_on_retry() {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(1);

        let mut attempt_count = 0;
        let _mock = mock("POST", "/api/v1/events")
            .with_status(move |_| {
                attempt_count += 1;
                if attempt_count == 1 {
                    500 // Fail on first attempt
                } else {
                    200 // Succeed on second attempt
                }
            })
            .with_body(move |_| {
                if attempt_count == 1 {
                    r#"{"error": "Temporary failure"}"#.to_string()
                } else {
                    r#"{"status": "success"}"#.to_string()
                }
            })
            .create();

        let result = client.send_event_batch_with_retry(&events, "test-token").await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_retry_mechanism_exhausted() {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(1);

        // Mock persistent failure
        let _mock = mock("POST", "/api/v1/events")
            .with_status(500)
            .with_body(r#"{"error": "Persistent failure"}"#)
            .expect(4) // Initial + 3 retries
            .create();

        let result = client.send_event_batch_with_retry(&events, "test-token").await;
        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_authentication_token_refresh() {
        let config = create_test_config();
        let client = HttpClient::new(config);

        // Mock token refresh endpoint
        let _token_mock = mock("POST", "/api/v1/auth/refresh")
            .with_status(200)
            .with_body(r#"{"token": "new-refreshed-token", "expires_in": 3600}"#)
            .create();

        // Mock events endpoint that requires fresh token
        let _events_mock = mock("POST", "/api/v1/events")
            .match_header("Authorization", "Bearer new-refreshed-token")
            .with_status(200)
            .with_body(r#"{"status": "success"}"#)
            .create();

        let events = create_test_events(1);
        let result = client.send_event_batch(&events, "expired-token").await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_heartbeat_send() {
        let config = create_test_config();
        let client = HttpClient::new(config);

        let heartbeat = HeartbeatData {
            agent_id: "test-agent-123".to_string(),
            timestamp: Utc::now(),
            status: "healthy".to_string(),
            version: "1.0.0".to_string(),
            uptime_seconds: 3600,
            memory_usage_mb: 50.5,
            cpu_usage_percent: 15.2,
        };

        // Mock heartbeat endpoint
        let _mock = mock("POST", "/api/v1/heartbeat")
            .match_header("Content-Type", "application/json")
            .with_status(200)
            .with_body(r#"{"status": "acknowledged"}"#)
            .create();

        let result = client.send_heartbeat(&heartbeat, "test-token").await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_agent_registration() {
        let config = create_test_config();
        let client = HttpClient::new(config);

        let registration = AgentRegistration {
            agent_id: "test-agent-123".to_string(),
            tenant_id: "test-tenant".to_string(),
            hostname: "test-host".to_string(),
            os: "Linux".to_string(),
            version: "1.0.0".to_string(),
            capabilities: vec!["process_monitoring".to_string(), "file_monitoring".to_string()],
        };

        // Mock registration endpoint
        let _mock = mock("POST", "/api/v1/agents/register")
            .with_status(201)
            .with_body(r#"{"status": "registered", "agent_id": "test-agent-123"}"#)
            .create();

        let result = client.register_agent(&registration).await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_configuration_sync() {
        let config = create_test_config();
        let client = HttpClient::new(config);

        // Mock config sync endpoint
        let _mock = mock("GET", "/api/v1/agents/config")
            .match_header("X-Agent-ID", "test-agent-123")
            .with_status(200)
            .with_body(r#"
            {
                "collection_interval": 45,
                "max_batch_size": 150,
                "enable_compression": true,
                "rules": [
                    {"name": "suspicious_process", "enabled": true},
                    {"name": "file_access", "enabled": false}
                ]
            }
            "#)
            .create();

        let result = client.sync_configuration("test-token").await;
        assert!(result.is_ok());

        let sync_config = result.unwrap();
        assert_eq!(sync_config.collection_interval, 45);
        assert_eq!(sync_config.max_batch_size, 150);
        assert!(sync_config.enable_compression);
        assert_eq!(sync_config.rules.len(), 2);
    }

    #[tokio::test]
    async fn test_rate_limiting() {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(1);

        // Mock rate limit response
        let _mock = mock("POST", "/api/v1/events")
            .with_status(429)
            .with_header("Retry-After", "60")
            .with_body(r#"{"error": "Rate limit exceeded"}"#)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_err());

        let error = result.unwrap_err();
        assert!(error.to_string().contains("429"));
    }

    #[tokio::test]
    async fn test_concurrent_requests() {
        let config = create_test_config();
        let client = Arc::new(HttpClient::new(config));

        let mut handles = vec![];

        // Mock concurrent requests
        let _mock = mock("POST", "/api/v1/events")
            .with_status(200)
            .with_body(r#"{"status": "success"}"#)
            .expect(5)
            .create();

        // Spawn concurrent requests
        for i in 0..5 {
            let client_clone = Arc::clone(&client);
            let events = create_test_events(1);
            let handle = tokio::spawn(async move {
                let result = client_clone.send_event_batch(&events, "test-token").await;
                assert!(result.is_ok());
            });
            handles.push(handle);
        }

        // Wait for all requests to complete
        for handle in handles {
            let _ = handle.await.unwrap();
        }
    }

    #[tokio::test]
    async fn test_request_metrics_collection() {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(2);

        // Mock successful response
        let _mock = mock("POST", "/api/v1/events")
            .with_status(200)
            .with_body(r#"{"status": "success"}"#)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_ok());

        // Check that metrics were collected
        let metrics = client.get_request_metrics().await;
        assert!(metrics.total_requests > 0);
        assert!(metrics.successful_requests > 0);
        assert_eq!(metrics.failed_requests, 0);
    }

    #[tokio::test]
    async fn test_connection_pooling() {
        let config = create_test_config();
        let client = HttpClient::new(config);

        // Send multiple requests to test connection reuse
        let _mock = mock("POST", "/api/v1/events")
            .with_status(200)
            .with_body(r#"{"status": "success"}"#)
            .expect(3)
            .create();

        for _ in 0..3 {
            let events = create_test_events(1);
            let result = client.send_event_batch(&events, "test-token").await;
            assert!(result.is_ok());
        }

        // Verify connection was reused (implementation detail)
        let metrics = client.get_connection_metrics().await;
        assert!(metrics.connections_created <= 2); // Should reuse connections
    }

    #[tokio::test]
    async fn test_tls_configuration() {
        let mut config = create_test_config();
        config.tls_ca_cert = Some("/path/to/ca.crt".to_string());
        config.tls_client_cert = Some("/path/to/client.crt".to_string());
        config.tls_client_key = Some("/path/to/client.key".to_string());

        let client = HttpClient::new(config);

        // Mock HTTPS endpoint
        let _mock = mock("POST", "/api/v1/events")
            .with_status(200)
            .with_body(r#"{"status": "success"}"#)
            .create();

        let events = create_test_events(1);
        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_ok());
    }

    #[tokio::test]
    async fn test_request_timeout_handling() {
        let config = create_test_config();
        let client = HttpClient::new(config);

        // Mock slow endpoint that times out
        let _mock = mock("POST", "/api/v1/events")
            .with_status(200)
            .with_body(delayed_response)
            .create();

        async fn delayed_response() -> String {
            tokio::time::sleep(tokio::time::Duration::from_secs(30)).await;
            r#"{"status": "success"}"#.to_string()
        }

        let events = create_test_events(1);
        let result = client.send_event_batch(&events, "test-token").await;
        // Should timeout and return error
        assert!(result.is_err());
    }

    #[tokio::test]
    async fn test_payload_size_limits() {
        let config = create_test_config();
        let client = HttpClient::new(config);

        // Create very large events to test size limits
        let large_event = TelemetryEvent {
            id: "large-event".to_string(),
            timestamp: Utc::now(),
            event_type: EventType::File,
            data: {
                let mut data = HashMap::new();
                data.insert("large_data".to_string(), serde_json::Value::String("x".repeat(1024 * 1024))); // 1MB string
                data
            },
            metadata: HashMap::new(),
        };

        let events = vec![large_event];

        // Mock response
        let _mock = mock("POST", "/api/v1/events")
            .with_status(413)
            .with_body(r#"{"error": "Payload too large"}"#)
            .create();

        let result = client.send_event_batch(&events, "test-token").await;
        assert!(result.is_err());
    }

    // Benchmark tests
    #[bench]
    fn bench_event_serialization(b: &mut test::Bencher) {
        let events = create_test_events(10);
        b.iter(|| {
            let _serialized = serde_json::to_string(&events).unwrap();
        });
    }

    #[bench]
    fn bench_request_preparation(b: &mut test::Bencher) {
        let config = create_test_config();
        let client = HttpClient::new(config);
        let events = create_test_events(5);

        b.iter(|| {
            let _request = client.prepare_request(&events, "test-token");
        });
    }
}</content>
<parameter name="filePath">/workspaces/insec/tests/unit/agent/network_client_test.rs
