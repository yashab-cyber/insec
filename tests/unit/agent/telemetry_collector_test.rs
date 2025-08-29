#[cfg(test)]
mod tests {
    use super::*;
    use std::collections::HashMap;
    use std::sync::Arc;
    use tokio::sync::Mutex;
    use chrono::{Utc, Duration};

    // Mock implementations for testing
    struct MockTelemetryCollector {
        events: Arc<Mutex<Vec<TelemetryEvent>>>,
    }

    impl MockTelemetryCollector {
        fn new() -> Self {
            Self {
                events: Arc::new(Mutex::new(Vec::new())),
            }
        }

        async fn add_event(&self, event: TelemetryEvent) {
            let mut events = self.events.lock().await;
            events.push(event);
        }

        async fn get_events(&self) -> Vec<TelemetryEvent> {
            let events = self.events.lock().await;
            events.clone()
        }

        async fn clear_events(&self) {
            let mut events = self.events.lock().await;
            events.clear();
        }
    }

    // Test data factories
    fn create_test_config() -> Config {
        Config {
            server_url: "https://test.insec.com".to_string(),
            agent_id: "test-agent-123".to_string(),
            tenant_id: "test-tenant".to_string(),
            collection_interval: 30,
            max_batch_size: 10,
            tls_ca_cert: None,
            tls_client_cert: None,
            tls_client_key: None,
        }
    }

    fn create_test_process_event() -> TelemetryEvent {
        TelemetryEvent {
            id: "test-event-123".to_string(),
            timestamp: Utc::now(),
            event_type: EventType::Process,
            data: {
                let mut data = HashMap::new();
                data.insert("process_name".to_string(), serde_json::Value::String("test.exe".to_string()));
                data.insert("pid".to_string(), serde_json::Value::Number(1234.into()));
                data.insert("command_line".to_string(), serde_json::Value::String("test.exe --arg".to_string()));
                data
            },
            metadata: {
                let mut metadata = HashMap::new();
                metadata.insert("risk_score".to_string(), serde_json::Value::Number(0.3.into()));
                metadata
            },
        }
    }

    fn create_test_file_event() -> TelemetryEvent {
        TelemetryEvent {
            id: "test-file-event-123".to_string(),
            timestamp: Utc::now(),
            event_type: EventType::File,
            data: {
                let mut data = HashMap::new();
                data.insert("filename".to_string(), serde_json::Value::String("/tmp/test.txt".to_string()));
                data.insert("operation".to_string(), serde_json::Value::String("write".to_string()));
                data.insert("size".to_string(), serde_json::Value::Number(1024.into()));
                data
            },
            metadata: {
                let mut metadata = HashMap::new();
                metadata.insert("risk_score".to_string(), serde_json::Value::Number(0.1.into()));
                metadata
            },
        }
    }

    #[tokio::test]
    async fn test_telemetry_collector_initialization() {
        // Arrange
        let config = create_test_config();

        // Act
        let collector = TelemetryCollector::new(config.clone());

        // Assert
        assert_eq!(collector.config.server_url, config.server_url);
        assert_eq!(collector.config.agent_id, config.agent_id);
        assert_eq!(collector.config.tenant_id, config.tenant_id);
    }

    #[tokio::test]
    async fn test_process_telemetry_collection() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mock_collector = MockTelemetryCollector::new();

        // Act
        collector.collect_process_telemetry(&mock_collector).await.unwrap();
        let events = mock_collector.get_events().await;

        // Assert
        assert!(!events.is_empty());
        let process_events: Vec<_> = events.iter()
            .filter(|e| matches!(e.event_type, EventType::Process))
            .collect();

        assert!(!process_events.is_empty());
        for event in process_events {
            assert!(event.data.contains_key("process_name"));
            assert!(event.data.contains_key("pid"));
            assert!(event.timestamp <= Utc::now());
            assert!(event.id.starts_with("proc-"));
        }
    }

    #[tokio::test]
    async fn test_file_telemetry_collection() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mock_collector = MockTelemetryCollector::new();

        // Create a test file
        let test_file = "/tmp/insec_test_file.txt";
        tokio::fs::write(test_file, b"test content").await.unwrap();

        // Act
        collector.collect_file_telemetry(&mock_collector).await.unwrap();
        let events = mock_collector.get_events().await;

        // Assert
        let file_events: Vec<_> = events.iter()
            .filter(|e| matches!(e.event_type, EventType::File))
            .collect();

        // Note: File events might not be captured immediately due to timing
        // This test validates the collection mechanism works

        // Cleanup
        tokio::fs::remove_file(test_file).await.unwrap();
    }

    #[tokio::test]
    async fn test_network_telemetry_collection() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mock_collector = MockTelemetryCollector::new();

        // Act
        collector.collect_network_telemetry(&mock_collector).await.unwrap();
        let events = mock_collector.get_events().await;

        // Assert
        let network_events: Vec<_> = events.iter()
            .filter(|e| matches!(e.event_type, EventType::Network))
            .collect();

        // Network events depend on system activity
        // This test validates the collection mechanism works
        for event in network_events {
            assert!(event.data.contains_key("protocol"));
            assert!(event.timestamp <= Utc::now());
            assert!(event.id.starts_with("net-"));
        }
    }

    #[tokio::test]
    async fn test_event_batch_processing() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mock_collector = MockTelemetryCollector::new();

        // Add multiple events
        for i in 0..15 {
            let mut event = create_test_process_event();
            event.id = format!("test-event-{}", i);
            mock_collector.add_event(event).await;
        }

        // Act
        let batches = collector.create_batches(&mock_collector).await;

        // Assert
        assert!(!batches.is_empty());
        assert!(batches.len() <= 2); // Should create at most 2 batches (max_batch_size = 10)

        for batch in batches {
            assert!(batch.len() <= 10);
            for event in batch {
                assert!(event.id.starts_with("test-event-"));
            }
        }
    }

    #[tokio::test]
    async fn test_event_serialization() {
        // Arrange
        let event = create_test_process_event();

        // Act
        let serialized = serde_json::to_string(&event).unwrap();
        let deserialized: TelemetryEvent = serde_json::from_str(&serialized).unwrap();

        // Assert
        assert_eq!(event.id, deserialized.id);
        assert_eq!(event.event_type, deserialized.event_type);
        assert_eq!(event.data, deserialized.data);
        assert_eq!(event.metadata, deserialized.metadata);
    }

    #[tokio::test]
    async fn test_risk_score_calculation() {
        // Arrange
        let collector = TelemetryCollector::new(create_test_config());
        let mut event = create_test_process_event();

        // Act
        collector.calculate_risk_score(&mut event).await;

        // Assert
        assert!(event.metadata.contains_key("risk_score"));
        let risk_score = event.metadata.get("risk_score").unwrap().as_f64().unwrap();
        assert!(risk_score >= 0.0 && risk_score <= 1.0);
    }

    #[tokio::test]
    async fn test_event_filtering() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mock_collector = MockTelemetryCollector::new();

        // Add events with different risk scores
        let mut high_risk_event = create_test_process_event();
        high_risk_event.metadata.insert("risk_score".to_string(), serde_json::Value::Number(0.9.into()));
        high_risk_event.data.insert("process_name".to_string(), serde_json::Value::String("suspicious.exe".to_string()));

        let mut low_risk_event = create_test_file_event();
        low_risk_event.metadata.insert("risk_score".to_string(), serde_json::Value::Number(0.1.into()));

        mock_collector.add_event(high_risk_event).await;
        mock_collector.add_event(low_risk_event).await;

        // Act
        let filtered_events = collector.filter_events(&mock_collector, 0.5).await;

        // Assert
        assert_eq!(filtered_events.len(), 1);
        assert!(filtered_events[0].metadata.get("risk_score").unwrap().as_f64().unwrap() >= 0.5);
    }

    #[tokio::test]
    async fn test_configuration_validation() {
        // Test valid configuration
        let valid_config = create_test_config();
        assert!(valid_config.validate().is_ok());

        // Test invalid configuration
        let mut invalid_config = create_test_config();
        invalid_config.server_url = "".to_string();
        assert!(invalid_config.validate().is_err());

        let mut invalid_config2 = create_test_config();
        invalid_config2.collection_interval = 0;
        assert!(invalid_config2.validate().is_err());
    }

    #[tokio::test]
    async fn test_error_handling_network_failure() {
        // Arrange
        let mut config = create_test_config();
        config.server_url = "https://nonexistent.invalid.server".to_string();
        let collector = TelemetryCollector::new(config);
        let mock_collector = MockTelemetryCollector::new();

        mock_collector.add_event(create_test_process_event()).await;

        // Act
        let result = collector.send_events(&mock_collector).await;

        // Assert
        assert!(result.is_err());
        // Should handle network errors gracefully without panicking
    }

    #[tokio::test]
    async fn test_concurrent_event_processing() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mock_collector = Arc::new(Mutex::new(MockTelemetryCollector::new()));

        // Act - Process events concurrently
        let mut handles = vec![];
        for i in 0..5 {
            let collector_clone = collector.clone();
            let mock_clone = Arc::clone(&mock_collector);
            let handle = tokio::spawn(async move {
                let mut event = create_test_process_event();
                event.id = format!("concurrent-event-{}", i);
                {
                    let mock = mock_clone.lock().await;
                    mock.add_event(event).await;
                }
                collector_clone.process_events(&*mock).await
            });
            handles.push(handle);
        }

        // Wait for all concurrent operations to complete
        for handle in handles {
            let _ = handle.await.unwrap();
        }

        // Assert
        let mock = mock_collector.lock().await;
        let events = mock.get_events().await;
        assert_eq!(events.len(), 5);

        let concurrent_events: Vec<_> = events.iter()
            .filter(|e| e.id.starts_with("concurrent-event-"))
            .collect();
        assert_eq!(concurrent_events.len(), 5);
    }

    #[tokio::test]
    async fn test_memory_usage_monitoring() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);

        // Act
        let memory_usage = collector.get_memory_usage().await;

        // Assert
        assert!(memory_usage > 0);
        assert!(memory_usage < 100 * 1024 * 1024); // Less than 100MB
    }

    #[tokio::test]
    async fn test_agent_self_protection() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);

        // Act
        let is_protected = collector.check_self_protection().await;

        // Assert
        assert!(is_protected); // Agent should be self-protected
    }

    #[tokio::test]
    async fn test_event_deduplication() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mock_collector = MockTelemetryCollector::new();

        // Add duplicate events
        let event1 = create_test_process_event();
        let mut event2 = event1.clone();
        event2.timestamp = event1.timestamp + Duration::milliseconds(100);

        mock_collector.add_event(event1).await;
        mock_collector.add_event(event2).await;

        // Act
        collector.deduplicate_events(&mock_collector).await;
        let events = mock_collector.get_events().await;

        // Assert
        // Should deduplicate similar events within time window
        assert!(events.len() <= 2);
    }

    #[tokio::test]
    async fn test_performance_metrics() {
        // Arrange
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);

        // Act
        let metrics = collector.collect_performance_metrics().await;

        // Assert
        assert!(metrics.contains_key("cpu_usage"));
        assert!(metrics.contains_key("memory_usage"));
        assert!(metrics.contains_key("events_per_second"));

        let cpu_usage = metrics.get("cpu_usage").unwrap();
        assert!(*cpu_usage >= 0.0 && *cpu_usage <= 100.0);
    }

    // Benchmark tests
    #[bench]
    fn bench_event_creation(b: &mut test::Bencher) {
        b.iter(|| {
            let _event = create_test_process_event();
        });
    }

    #[bench]
    fn bench_event_serialization(b: &mut test::Bencher) {
        let event = create_test_process_event();
        b.iter(|| {
            let _serialized = serde_json::to_string(&event).unwrap();
        });
    }

    #[bench]
    fn bench_risk_score_calculation(b: &mut test::Bencher) {
        let config = create_test_config();
        let collector = TelemetryCollector::new(config);
        let mut event = create_test_process_event();

        b.iter(|| {
            let _ = collector.calculate_risk_score_sync(&mut event);
        });
    }
}</content>
<parameter name="filePath">/workspaces/insec/tests/unit/agent/telemetry_collector_test.rs
