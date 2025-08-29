#[cfg(test)]
mod tests {
    use super::*;
    use std::env;
    use std::fs;
    use tempfile::TempDir;
    use serde_json;

    // Test data factories
    fn create_valid_config() -> Config {
        Config {
            server_url: "https://api.insec.com".to_string(),
            agent_id: "test-agent-123".to_string(),
            tenant_id: "test-tenant".to_string(),
            collection_interval: 30,
            max_batch_size: 100,
            tls_ca_cert: Some("/path/to/ca.crt".to_string()),
            tls_client_cert: Some("/path/to/client.crt".to_string()),
            tls_client_key: Some("/path/to/client.key".to_string()),
            log_level: "info".to_string(),
            enable_compression: true,
            retry_attempts: 3,
            retry_delay: 5,
            heartbeat_interval: 60,
        }
    }

    fn create_minimal_config() -> Config {
        Config {
            server_url: "https://api.insec.com".to_string(),
            agent_id: "test-agent-123".to_string(),
            tenant_id: "test-tenant".to_string(),
            collection_interval: 30,
            max_batch_size: 100,
            tls_ca_cert: None,
            tls_client_cert: None,
            tls_client_key: None,
            log_level: "info".to_string(),
            enable_compression: false,
            retry_attempts: 3,
            retry_delay: 5,
            heartbeat_interval: 60,
        }
    }

    #[test]
    fn test_config_validation_valid() {
        let config = create_valid_config();
        assert!(config.validate().is_ok());
    }

    #[test]
    fn test_config_validation_invalid_server_url() {
        let mut config = create_valid_config();
        config.server_url = "".to_string();
        assert!(config.validate().is_err());

        config.server_url = "not-a-url".to_string();
        assert!(config.validate().is_err());

        config.server_url = "ftp://api.insec.com".to_string();
        assert!(config.validate().is_err());
    }

    #[test]
    fn test_config_validation_invalid_agent_id() {
        let mut config = create_valid_config();
        config.agent_id = "".to_string();
        assert!(config.validate().is_err());

        config.agent_id = "invalid@agent".to_string();
        assert!(config.validate().is_err());
    }

    #[test]
    fn test_config_validation_invalid_collection_interval() {
        let mut config = create_valid_config();
        config.collection_interval = 0;
        assert!(config.validate().is_err());

        config.collection_interval = 3601; // > 1 hour
        assert!(config.validate().is_err());
    }

    #[test]
    fn test_config_validation_invalid_batch_size() {
        let mut config = create_valid_config();
        config.max_batch_size = 0;
        assert!(config.validate().is_err());

        config.max_batch_size = 10001; // > 10000
        assert!(config.validate().is_err());
    }

    #[test]
    fn test_config_validation_tls_certificates() {
        let mut config = create_valid_config();

        // If client cert is provided, client key must also be provided
        config.tls_client_cert = Some("/path/to/client.crt".to_string());
        config.tls_client_key = None;
        assert!(config.validate().is_err());

        // If client key is provided, client cert must also be provided
        config.tls_client_cert = None;
        config.tls_client_key = Some("/path/to/client.key".to_string());
        assert!(config.validate().is_err());

        // Both should be provided or both should be None
        config.tls_client_cert = Some("/path/to/client.crt".to_string());
        config.tls_client_key = Some("/path/to/client.key".to_string());
        assert!(config.validate().is_ok());
    }

    #[test]
    fn test_config_validation_retry_settings() {
        let mut config = create_valid_config();
        config.retry_attempts = 0;
        assert!(config.validate().is_err());

        config.retry_attempts = 11; // > 10
        assert!(config.validate().is_err());

        config.retry_attempts = 3;
        config.retry_delay = 0;
        assert!(config.validate().is_err());

        config.retry_delay = 301; // > 300
        assert!(config.validate().is_err());
    }

    #[test]
    fn test_config_validation_heartbeat_interval() {
        let mut config = create_valid_config();
        config.heartbeat_interval = 0;
        assert!(config.validate().is_err());

        config.heartbeat_interval = 3601; // > 1 hour
        assert!(config.validate().is_err());
    }

    #[test]
    fn test_config_from_file() {
        let temp_dir = TempDir::new().unwrap();
        let config_path = temp_dir.path().join("config.json");

        let config_data = r#"
        {
            "server_url": "https://api.insec.com",
            "agent_id": "test-agent-123",
            "tenant_id": "test-tenant",
            "collection_interval": 30,
            "max_batch_size": 100,
            "tls_ca_cert": "/path/to/ca.crt",
            "tls_client_cert": "/path/to/client.crt",
            "tls_client_key": "/path/to/client.key",
            "log_level": "info",
            "enable_compression": true,
            "retry_attempts": 3,
            "retry_delay": 5,
            "heartbeat_interval": 60
        }
        "#;

        fs::write(&config_path, config_data).unwrap();

        let config = Config::from_file(config_path.to_str().unwrap()).unwrap();
        assert_eq!(config.server_url, "https://api.insec.com");
        assert_eq!(config.agent_id, "test-agent-123");
        assert_eq!(config.tenant_id, "test-tenant");
        assert_eq!(config.collection_interval, 30);
        assert_eq!(config.max_batch_size, 100);
        assert!(config.tls_ca_cert.is_some());
        assert!(config.tls_client_cert.is_some());
        assert!(config.tls_client_key.is_some());
        assert_eq!(config.log_level, "info");
        assert!(config.enable_compression);
        assert_eq!(config.retry_attempts, 3);
        assert_eq!(config.retry_delay, 5);
        assert_eq!(config.heartbeat_interval, 60);
    }

    #[test]
    fn test_config_from_file_invalid_json() {
        let temp_dir = TempDir::new().unwrap();
        let config_path = temp_dir.path().join("invalid_config.json");

        let invalid_config_data = r#"
        {
            "server_url": "https://api.insec.com",
            "agent_id": "test-agent-123",
            "invalid_json": ,
        }
        "#;

        fs::write(&config_path, invalid_config_data).unwrap();

        let result = Config::from_file(config_path.to_str().unwrap());
        assert!(result.is_err());
    }

    #[test]
    fn test_config_from_file_missing_required_fields() {
        let temp_dir = TempDir::new().unwrap();
        let config_path = temp_dir.path().join("incomplete_config.json");

        let incomplete_config_data = r#"
        {
            "server_url": "https://api.insec.com"
        }
        "#;

        fs::write(&config_path, incomplete_config_data).unwrap();

        let result = Config::from_file(config_path.to_str().unwrap());
        assert!(result.is_err());
    }

    #[test]
    fn test_config_from_env_vars() {
        // Set environment variables
        env::set_var("INSEC_SERVER_URL", "https://env-api.insec.com");
        env::set_var("INSEC_AGENT_ID", "env-agent-456");
        env::set_var("INSEC_TENANT_ID", "env-tenant");
        env::set_var("INSEC_COLLECTION_INTERVAL", "45");
        env::set_var("INSEC_MAX_BATCH_SIZE", "200");
        env::set_var("INSEC_TLS_CA_CERT", "/env/path/ca.crt");
        env::set_var("INSEC_TLS_CLIENT_CERT", "/env/path/client.crt");
        env::set_var("INSEC_TLS_CLIENT_KEY", "/env/path/client.key");
        env::set_var("INSEC_LOG_LEVEL", "debug");
        env::set_var("INSEC_ENABLE_COMPRESSION", "false");
        env::set_var("INSEC_RETRY_ATTEMPTS", "5");
        env::set_var("INSEC_RETRY_DELAY", "10");
        env::set_var("INSEC_HEARTBEAT_INTERVAL", "120");

        let config = Config::from_env().unwrap();

        assert_eq!(config.server_url, "https://env-api.insec.com");
        assert_eq!(config.agent_id, "env-agent-456");
        assert_eq!(config.tenant_id, "env-tenant");
        assert_eq!(config.collection_interval, 45);
        assert_eq!(config.max_batch_size, 200);
        assert_eq!(config.tls_ca_cert, Some("/env/path/ca.crt".to_string()));
        assert_eq!(config.tls_client_cert, Some("/env/path/client.crt".to_string()));
        assert_eq!(config.tls_client_key, Some("/env/path/client.key".to_string()));
        assert_eq!(config.log_level, "debug");
        assert!(!config.enable_compression);
        assert_eq!(config.retry_attempts, 5);
        assert_eq!(config.retry_delay, 10);
        assert_eq!(config.heartbeat_interval, 120);

        // Clean up environment variables
        env::remove_var("INSEC_SERVER_URL");
        env::remove_var("INSEC_AGENT_ID");
        env::remove_var("INSEC_TENANT_ID");
        env::remove_var("INSEC_COLLECTION_INTERVAL");
        env::remove_var("INSEC_MAX_BATCH_SIZE");
        env::remove_var("INSEC_TLS_CA_CERT");
        env::remove_var("INSEC_TLS_CLIENT_CERT");
        env::remove_var("INSEC_TLS_CLIENT_KEY");
        env::remove_var("INSEC_LOG_LEVEL");
        env::remove_var("INSEC_ENABLE_COMPRESSION");
        env::remove_var("INSEC_RETRY_ATTEMPTS");
        env::remove_var("INSEC_RETRY_DELAY");
        env::remove_var("INSEC_HEARTBEAT_INTERVAL");
    }

    #[test]
    fn test_config_from_env_vars_missing_required() {
        // Don't set required environment variables
        let result = Config::from_env();
        assert!(result.is_err());
    }

    #[test]
    fn test_config_from_env_vars_invalid_values() {
        // Set invalid environment variables
        env::set_var("INSEC_SERVER_URL", "https://api.insec.com");
        env::set_var("INSEC_AGENT_ID", "test-agent");
        env::set_var("INSEC_TENANT_ID", "test-tenant");
        env::set_var("INSEC_COLLECTION_INTERVAL", "0"); // Invalid: must be > 0
        env::set_var("INSEC_MAX_BATCH_SIZE", "100");

        let result = Config::from_env();
        assert!(result.is_err());

        // Clean up
        env::remove_var("INSEC_SERVER_URL");
        env::remove_var("INSEC_AGENT_ID");
        env::remove_var("INSEC_TENANT_ID");
        env::remove_var("INSEC_COLLECTION_INTERVAL");
        env::remove_var("INSEC_MAX_BATCH_SIZE");
    }

    #[test]
    fn test_config_merge_file_and_env() {
        let temp_dir = TempDir::new().unwrap();
        let config_path = temp_dir.path().join("base_config.json");

        let base_config_data = r#"
        {
            "server_url": "https://api.insec.com",
            "agent_id": "file-agent-123",
            "tenant_id": "file-tenant",
            "collection_interval": 30,
            "max_batch_size": 100,
            "log_level": "info",
            "enable_compression": true,
            "retry_attempts": 3,
            "retry_delay": 5,
            "heartbeat_interval": 60
        }
        "#;

        fs::write(&config_path, base_config_data).unwrap();

        // Set environment variables to override some values
        env::set_var("INSEC_AGENT_ID", "env-agent-456");
        env::set_var("INSEC_COLLECTION_INTERVAL", "45");
        env::set_var("INSEC_ENABLE_COMPRESSION", "false");

        let config = Config::from_file_with_env_override(config_path.to_str().unwrap()).unwrap();

        // Values from file
        assert_eq!(config.server_url, "https://api.insec.com");
        assert_eq!(config.tenant_id, "file-tenant");
        assert_eq!(config.max_batch_size, 100);

        // Values overridden by environment
        assert_eq!(config.agent_id, "env-agent-456");
        assert_eq!(config.collection_interval, 45);
        assert!(!config.enable_compression);

        // Clean up
        env::remove_var("INSEC_AGENT_ID");
        env::remove_var("INSEC_COLLECTION_INTERVAL");
        env::remove_var("INSEC_ENABLE_COMPRESSION");
    }

    #[test]
    fn test_config_default_values() {
        let config = Config::default();

        assert_eq!(config.collection_interval, 30);
        assert_eq!(config.max_batch_size, 100);
        assert_eq!(config.log_level, "info");
        assert!(!config.enable_compression);
        assert_eq!(config.retry_attempts, 3);
        assert_eq!(config.retry_delay, 5);
        assert_eq!(config.heartbeat_interval, 60);
        assert!(config.tls_ca_cert.is_none());
        assert!(config.tls_client_cert.is_none());
        assert!(config.tls_client_key.is_none());
    }

    #[test]
    fn test_config_serialization() {
        let config = create_valid_config();

        // Serialize to JSON
        let serialized = serde_json::to_string_pretty(&config).unwrap();
        assert!(serialized.contains("server_url"));
        assert!(serialized.contains("agent_id"));
        assert!(serialized.contains("collection_interval"));

        // Deserialize back
        let deserialized: Config = serde_json::from_str(&serialized).unwrap();
        assert_eq!(config, deserialized);
    }

    #[test]
    fn test_config_display() {
        let config = create_minimal_config();
        let display = format!("{}", config);

        assert!(display.contains("server_url"));
        assert!(display.contains("agent_id"));
        assert!(display.contains("tenant_id"));
        assert!(display.contains("https://api.insec.com"));
        assert!(display.contains("test-agent-123"));
    }

    #[test]
    fn test_config_clone() {
        let config = create_valid_config();
        let cloned = config.clone();

        assert_eq!(config, cloned);
        assert_eq!(config.server_url, cloned.server_url);
        assert_eq!(config.agent_id, cloned.agent_id);
    }

    #[test]
    fn test_config_debug() {
        let config = create_minimal_config();
        let debug = format!("{:?}", config);

        assert!(debug.contains("Config"));
        assert!(debug.contains("server_url"));
        assert!(debug.contains("agent_id"));
    }

    #[test]
    fn test_config_partial_eq() {
        let config1 = create_valid_config();
        let config2 = create_valid_config();
        let mut config3 = create_valid_config();
        config3.agent_id = "different-agent".to_string();

        assert_eq!(config1, config2);
        assert_ne!(config1, config3);
    }

    #[test]
    fn test_config_file_not_found() {
        let result = Config::from_file("/nonexistent/config.json");
        assert!(result.is_err());
    }

    #[test]
    fn test_config_env_var_parsing_errors() {
        // Set invalid numeric values
        env::set_var("INSEC_SERVER_URL", "https://api.insec.com");
        env::set_var("INSEC_AGENT_ID", "test-agent");
        env::set_var("INSEC_TENANT_ID", "test-tenant");
        env::set_var("INSEC_COLLECTION_INTERVAL", "not-a-number");
        env::set_var("INSEC_MAX_BATCH_SIZE", "100");

        let result = Config::from_env();
        assert!(result.is_err());

        // Clean up
        env::remove_var("INSEC_SERVER_URL");
        env::remove_var("INSEC_AGENT_ID");
        env::remove_var("INSEC_TENANT_ID");
        env::remove_var("INSEC_COLLECTION_INTERVAL");
        env::remove_var("INSEC_MAX_BATCH_SIZE");
    }

    #[test]
    fn test_config_with_special_characters() {
        let mut config = create_valid_config();
        config.agent_id = "test_agent-123_special.chars".to_string();
        config.tenant_id = "tenant@domain.com".to_string();

        assert!(config.validate().is_ok());
    }

    #[test]
    fn test_config_large_values() {
        let mut config = create_valid_config();
        config.max_batch_size = 10000; // Maximum allowed
        config.collection_interval = 3600; // Maximum allowed
        config.retry_attempts = 10; // Maximum allowed
        config.retry_delay = 300; // Maximum allowed
        config.heartbeat_interval = 3600; // Maximum allowed

        assert!(config.validate().is_ok());
    }

    #[test]
    fn test_config_edge_case_urls() {
        let mut config = create_valid_config();

        // Test various valid URL formats
        let valid_urls = vec![
            "https://api.insec.com",
            "https://api.insec.com:8443",
            "https://api.insec.com/path",
            "https://subdomain.api.insec.com",
            "http://localhost:8080",
        ];

        for url in valid_urls {
            config.server_url = url.to_string();
            assert!(config.validate().is_ok(), "URL {} should be valid", url);
        }
    }
}</content>
<parameter name="filePath">/workspaces/insec/tests/unit/agent/config_test.rs
