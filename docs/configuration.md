# INSEC Configuration Guide

This guide covers all configuration options for INSEC components, including server, agent, and console settings.

## 📁 Configuration Files

### Server Configuration
Location: `/etc/insec/server.yaml` or `./server/config.yaml`

### Agent Configuration
Location: `/etc/insec/agent.yaml` or `./agent/config.yaml`

### Console Configuration
Location: `./ui/.env` or `./ui/config.json`

## ⚙️ Server Configuration

### Basic Server Settings
```yaml
server:
  # Server binding
  host: "0.0.0.0"
  port: 8080

  # TLS configuration
  tls:
    enabled: true
    cert_file: "/etc/ssl/certs/insec.crt"
    key_file: "/etc/ssl/private/insec.key"
    ca_cert: "/etc/ssl/certs/ca.crt"

  # CORS settings
  cors:
    allowed_origins: ["https://console.insec.com"]
    allowed_methods: ["GET", "POST", "PUT", "DELETE"]
    allowed_headers: ["Authorization", "Content-Type"]
    allow_credentials: true

  # Rate limiting
  rate_limit:
    enabled: true
    requests_per_minute: 1000
    burst_limit: 100
```

### Database Configuration
```yaml
database:
  # PostgreSQL connection
  host: "localhost"
  port: 5432
  database: "insec"
  username: "insec_user"
  password: "your_secure_password"
  ssl_mode: "require"

  # Connection pool
  max_connections: 20
  max_idle_connections: 5
  connection_timeout: "30s"

  # Migration settings
  auto_migrate: true
  migration_path: "./migrations"
```

### Redis Configuration
```yaml
redis:
  # Redis connection
  host: "localhost"
  port: 6379
  password: ""
  database: 0

  # Connection pool
  max_connections: 10
  min_idle_connections: 2

  # Key prefixes
  key_prefix: "insec:"
  session_ttl: "24h"
  cache_ttl: "1h"
```

### Security Configuration
```yaml
security:
  # JWT settings
  jwt:
    secret: "your-256-bit-secret"
    expiration: "24h"
    refresh_expiration: "168h"

  # mTLS settings
  mtls:
    enabled: true
    ca_cert: "/etc/ssl/certs/ca.crt"
    server_cert: "/etc/ssl/certs/server.crt"
    server_key: "/etc/ssl/private/server.key"
    client_ca_cert: "/etc/ssl/certs/client-ca.crt"

  # Password policy
  password_policy:
    min_length: 12
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_symbols: true

  # Session management
  session:
    timeout: "8h"
    max_concurrent_sessions: 5
```

### Event Processing Configuration
```yaml
events:
  # Ingestion settings
  ingestion:
    max_batch_size: 100
    max_queue_size: 10000
    processing_timeout: "30s"

  # Processing pipeline
  processing:
    workers: 4
    queue_size: 1000
    retry_attempts: 3
    retry_delay: "5s"

  # Storage settings
  storage:
    retention_days: 365
    compression: "gzip"
    encryption: true
```

### Alerting Configuration
```yaml
alerting:
  # Alert rules
  rules:
    max_alerts_per_hour: 100
    alert_cooldown: "5m"
    auto_close_after: "24h"

  # Notification channels
  notifications:
    email:
      enabled: true
      smtp_host: "smtp.gmail.com"
      smtp_port: 587
      username: "alerts@insec.com"
      password: "your_app_password"
    slack:
      enabled: true
      webhook_url: "https://hooks.slack.com/services/..."
    webhook:
      enabled: true
      url: "https://your-webhook.com/alerts"
      secret: "webhook_secret"
```

### UEBA Configuration
```yaml
ueba:
  # Baseline calculation
  baseline:
    window_days: 30
    update_interval: "1h"
    min_samples: 100

  # Risk scoring
  risk_scoring:
    enabled: true
    model_path: "./models/risk_model.pkl"
    features: ["login_time", "location", "device", "behavior"]

  # Anomaly detection
  anomaly_detection:
    algorithm: "isolation_forest"
    contamination: 0.1
    threshold: 0.8
```

## 🤖 Agent Configuration

### Basic Agent Settings
```yaml
agent:
  # Agent identity
  id: "agent-001"
  tenant_id: "tenant-123"
  group: "production"

  # Server connection
  server:
    url: "https://api.insec.com"
    tls:
      ca_cert: "/etc/ssl/certs/ca.crt"
      client_cert: "/etc/ssl/certs/agent.crt"
      client_key: "/etc/ssl/private/agent.key"

  # Collection settings
  collection:
    interval: 30
    max_batch_size: 50
    queue_size: 1000
    retry_attempts: 3
```

### Telemetry Collection
```yaml
telemetry:
  # Process monitoring
  process:
    enabled: true
    include_system_processes: false
    hash_calculation: true
    command_line_capture: true

  # File system monitoring
  filesystem:
    enabled: true
    paths: ["/home", "/tmp", "/var/log"]
    exclude_patterns: ["*.log", "*.tmp"]
    hash_calculation: true

  # Network monitoring
  network:
    enabled: true
    capture_dns: true
    capture_http: false
    max_connections: 1000

  # User activity monitoring
  user_activity:
    enabled: true
    capture_logins: true
    capture_commands: true
    capture_file_access: true
```

### Policy Configuration
```yaml
policies:
  # Policy management
  enabled: true
  update_interval: "5m"
  cache_size: 100

  # Local enforcement
  enforcement:
    enabled: true
    block_mode: false  # Alert-only mode
    quarantine_path: "/var/quarantine"

  # Response actions
  responses:
    isolate_network: true
    kill_process: true
    quarantine_file: true
    alert_only: false
```

### Performance Configuration
```yaml
performance:
  # Resource limits
  memory_limit: "512MB"
  cpu_limit: "50%"

  # Collection optimization
  batch_size: 100
  compression: "gzip"
  encryption: true

  # Monitoring
  metrics:
    enabled: true
    interval: "60s"
    endpoint: "http://localhost:9090/metrics"
```

## 🖥️ Console Configuration

### Environment Variables
```bash
# API Configuration
REACT_APP_API_URL=https://api.insec.com
REACT_APP_API_VERSION=v1

# Authentication
REACT_APP_AUTH_URL=https://auth.insec.com
REACT_APP_CLIENT_ID=your-client-id

# Features
REACT_APP_ENABLE_UEBA=true
REACT_APP_ENABLE_AUTOMATION=true

# UI Settings
REACT_APP_THEME=dark
REACT_APP_TIMEZONE=UTC
```

### Runtime Configuration
```json
{
  "api": {
    "baseUrl": "https://api.insec.com",
    "timeout": 30000,
    "retries": 3
  },
  "auth": {
    "clientId": "your-client-id",
    "authority": "https://auth.insec.com",
    "redirectUri": "https://console.insec.com/callback"
  },
  "features": {
    "ueba": true,
    "automation": true,
    "reporting": true,
    "integrations": true
  },
  "ui": {
    "theme": "dark",
    "locale": "en-US",
    "timezone": "UTC",
    "dateFormat": "YYYY-MM-DD HH:mm:ss"
  }
}
```

## 🔧 Advanced Configuration

### High Availability Configuration
```yaml
ha:
  # Cluster settings
  cluster:
    enabled: true
    node_id: "node-01"
    peers: ["node-02", "node-03"]

  # Load balancing
  load_balancer:
    algorithm: "round_robin"
    health_check_interval: "30s"

  # Failover
  failover:
    enabled: true
    timeout: "30s"
    retry_attempts: 3
```

### Monitoring Configuration
```yaml
monitoring:
  # Metrics collection
  metrics:
    enabled: true
    interval: "15s"
    exporters: ["prometheus", "datadog"]

  # Logging
  logging:
    level: "info"
    format: "json"
    outputs: ["stdout", "file", "/var/log/insec/"]
    rotation: "daily"

  # Tracing
  tracing:
    enabled: true
    service_name: "insec-server"
    jaeger_endpoint: "http://jaeger:14268/api/traces"
```

### Integration Configuration
```yaml
integrations:
  # SIEM integration
  siem:
    splunk:
      enabled: true
      host: "splunk.company.com"
      port: 8088
      token: "your-splunk-token"
    elastic:
      enabled: true
      host: "elasticsearch.company.com"
      port: 9200
      api_key: "your-api-key"

  # Identity providers
  identity:
    okta:
      enabled: true
      domain: "company.okta.com"
      client_id: "your-client-id"
      client_secret: "your-client-secret"
    azure:
      enabled: true
      tenant_id: "your-tenant-id"
      client_id: "your-client-id"
      client_secret: "your-client-secret"

  # Ticketing systems
  ticketing:
    jira:
      enabled: true
      host: "company.atlassian.net"
      username: "insec@company.com"
      api_token: "your-api-token"
    servicenow:
      enabled: true
      instance: "company.service-now.com"
      username: "insec@company.com"
      password: "your-password"
```

## 🔒 Security Best Practices

### Secret Management
```yaml
secrets:
  # Use external secret management
  provider: "vault"  # or "aws-secretsmanager", "gcp-secretmanager"
  endpoint: "https://vault.company.com"
  token: "${VAULT_TOKEN}"

  # Key rotation
  rotation:
    enabled: true
    interval: "30d"
    overlap: "24h"
```

### Network Security
```yaml
network:
  # Firewall rules
  firewall:
    allowed_ips: ["10.0.0.0/8", "172.16.0.0/12"]
    blocked_ports: [22, 3389]

  # VPN requirements
  vpn:
    required: true
    allowed_cidrs: ["10.0.0.0/8"]

  # Proxy settings
  proxy:
    http_proxy: "http://proxy.company.com:8080"
    https_proxy: "http://proxy.company.com:8080"
    no_proxy: "localhost,127.0.0.1"
```

## 📊 Performance Tuning

### Database Optimization
```yaml
database:
  # Query optimization
  query_cache_size: "256MB"
  temp_buffer_size: "64MB"

  # Index settings
  indexes:
    event_timestamp: true
    event_type: true
    host_id: true
    user_id: true

  # Maintenance
  maintenance:
    autovacuum: true
    analyze_threshold: 50
    vacuum_threshold: 100
```

### Caching Configuration
```yaml
cache:
  # Redis cache settings
  redis:
    ttl: "1h"
    max_memory: "512MB"
    eviction_policy: "allkeys-lru"

  # Application cache
  application:
    policy_cache_ttl: "5m"
    user_cache_ttl: "10m"
    config_cache_ttl: "1h"
```

## 🔄 Configuration Management

### Environment-Specific Configurations
```yaml
# config.development.yaml
environment: "development"
debug: true
log_level: "debug"

# config.staging.yaml
environment: "staging"
debug: false
log_level: "info"

# config.production.yaml
environment: "production"
debug: false
log_level: "warn"
```

### Configuration Validation
```yaml
validation:
  # Schema validation
  schema:
    enabled: true
    strict: true

  # Configuration testing
  testing:
    enabled: true
    test_endpoints: ["health", "ready"]
```

## 📝 Configuration Examples

### Minimal Configuration
```yaml
server:
  host: "0.0.0.0"
  port: 8080

database:
  host: "localhost"
  database: "insec"
  username: "insec"

agent:
  server_url: "http://localhost:8080"
  collection_interval: 60
```

### Production Configuration
```yaml
server:
  host: "0.0.0.0"
  port: 8443
  tls:
    enabled: true
    cert_file: "/etc/ssl/certs/insec.crt"
    key_file: "/etc/ssl/private/insec.key"

database:
  host: "prod-db.company.com"
  port: 5432
  database: "insec_prod"
  username: "insec_prod"
  ssl_mode: "require"
  max_connections: 50

security:
  jwt_secret: "${JWT_SECRET}"
  mtls_enabled: true

monitoring:
  metrics_enabled: true
  logging_level: "info"
```

## 🔍 Configuration Validation

### Validate Configuration
```bash
# Server configuration validation
./insec-server -config /etc/insec/server.yaml -validate

# Agent configuration validation
./insec-agent -config /etc/insec/agent.yaml -validate
```

### Configuration Testing
```bash
# Test database connection
./insec-server -config /etc/insec/server.yaml -test-db

# Test agent connectivity
./insec-agent -config /etc/insec/agent.yaml -test-connection
```

## 📞 Support

For configuration assistance:
- **Documentation**: This configuration guide
- **Examples**: Check `config/examples/` directory
- **Community**: [GitHub Discussions](https://github.com/yashab-cyber/insec/discussions)
- **Support**: config-support@insec.com

---

**Last updated:** August 29, 2025</content>
<parameter name="filePath">/workspaces/insec/docs/configuration.md
