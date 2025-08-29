# INSEC Troubleshooting Guide

This guide helps you diagnose and resolve common issues with INSEC deployments.

## 🔍 Diagnostic Tools

### System Information Collection
```bash
# Collect system information
insec-diagnostics collect --output diagnostic-report.tar.gz

# Include logs, configuration, and system metrics
# Upload to support for analysis
```

### Health Check Commands
```bash
# Check server health
curl -f https://api.insec.com/health

# Check agent connectivity
sudo insec-agent --test-connection

# Check database connectivity
psql -U insec -d insec -c "SELECT version();"

# Check Redis connectivity
redis-cli -h localhost ping
```

## 🚨 Common Issues and Solutions

### 1. Agent Won't Start

#### Symptoms
- Agent service fails to start
- Error: "Failed to bind to port"
- Logs show permission denied

#### Solutions

**Check Configuration File**
```bash
# Validate agent configuration
sudo insec-agent -config /etc/insec/agent.yaml -validate

# Check file permissions
ls -la /etc/insec/agent.yaml
sudo chown insec:insec /etc/insec/agent.yaml
```

**Check System Resources**
```bash
# Check available memory
free -h

# Check disk space
df -h /opt/insec

# Check running processes
ps aux | grep insec
```

**Check Network Connectivity**
```bash
# Test server connectivity
telnet api.insec.com 8080

# Check DNS resolution
nslookup api.insec.com

# Check firewall rules
sudo ufw status
sudo iptables -L
```

### 2. Server Connection Issues

#### Symptoms
- Agent can't connect to server
- TLS handshake failures
- Certificate validation errors

#### Solutions

**Certificate Issues**
```bash
# Check certificate validity
openssl x509 -in /etc/ssl/certs/agent.crt -text -noout | grep "Not After"

# Verify certificate chain
openssl verify -CAfile /etc/ssl/certs/ca.crt /etc/ssl/certs/agent.crt

# Check certificate permissions
ls -la /etc/ssl/certs/
ls -la /etc/ssl/private/
```

**Network Issues**
```bash
# Test basic connectivity
ping api.insec.com

# Check port availability
netstat -tlnp | grep 8080

# Test TLS connection
openssl s_client -connect api.insec.com:8080 -servername api.insec.com
```

**Server Configuration**
```yaml
# Check server TLS configuration
server:
  tls:
    enabled: true
    cert_file: "/etc/ssl/certs/server.crt"
    key_file: "/etc/ssl/private/server.key"
    ca_cert: "/etc/ssl/certs/ca.crt"
```

### 3. Database Connection Problems

#### Symptoms
- Server fails to start
- "Connection refused" errors
- Slow query performance

#### Solutions

**Connection Configuration**
```bash
# Test database connection
psql -h localhost -U insec -d insec

# Check PostgreSQL service
sudo systemctl status postgresql

# Check connection pool settings
# server.yaml
database:
  max_connections: 20
  max_idle_connections: 5
  connection_timeout: "30s"
```

**Performance Issues**
```sql
-- Check active connections
SELECT count(*) FROM pg_stat_activity;

-- Check slow queries
SELECT pid, now() - pg_stat_activity.query_start AS duration, query
FROM pg_stat_activity
WHERE state = 'active'
ORDER BY duration DESC;

-- Check table sizes
SELECT schemaname, tablename, pg_size_pretty(pg_total_relation_size(schemaname||'.'||tablename))
FROM pg_tables
ORDER BY pg_total_relation_size(schemaname||'.'||tablename) DESC;
```

**Database Maintenance**
```bash
# Vacuum database
psql -U insec -d insec -c "VACUUM ANALYZE;"

# Reindex tables
psql -U insec -d insec -c "REINDEX DATABASE insec;"

# Check for corruption
psql -U insec -d insec -c "SELECT * FROM pg_stat_database WHERE datname = 'insec';"
```

### 4. High Resource Usage

#### Symptoms
- High CPU usage
- High memory consumption
- Slow response times

#### Solutions

**CPU Usage**
```bash
# Check CPU usage by process
top -p $(pgrep insec)

# Profile application
go tool pprof http://localhost:8080/debug/pprof/profile

# Check for infinite loops or heavy computations
```

**Memory Usage**
```bash
# Check memory usage
ps aux --sort=-%mem | head -10

# Check for memory leaks
go tool pprof http://localhost:8080/debug/pprof/heap

# Configure memory limits
# server.yaml
performance:
  memory_limit: "1GB"
  gc_percentage: 100
```

**Disk Usage**
```bash
# Check disk usage
df -h

# Find large files
find /var/log/insec -type f -size +100M

# Configure log rotation
# /etc/logrotate.d/insec
/var/log/insec/*.log {
    daily
    rotate 30
    compress
    missingok
    notifempty
}
```

### 5. Authentication Issues

#### Symptoms
- Login failures
- Invalid token errors
- MFA problems

#### Solutions

**JWT Configuration**
```yaml
# Check JWT settings
security:
  jwt:
    secret: "your-256-bit-secret"
    expiration: "24h"
    issuer: "insec-server"

# Validate JWT secret length
echo -n "your-secret" | wc -c  # Should be 32+ bytes
```

**MFA Issues**
```bash
# Reset MFA for user
curl -X POST https://api.insec.com/v1/users/{user_id}/mfa/reset \
  -H "Authorization: Bearer {admin_token}"

# Check TOTP configuration
# Ensure system time is synchronized
timedatectl status
```

**Password Policies**
```yaml
# Check password requirements
security:
  password_policy:
    min_length: 12
    require_uppercase: true
    require_lowercase: true
    require_numbers: true
    require_symbols: true
```

### 6. Event Ingestion Problems

#### Symptoms
- Events not appearing in console
- High ingestion latency
- Event loss

#### Solutions

**Ingestion Pipeline**
```yaml
# Check ingestion settings
events:
  ingestion:
    max_batch_size: 100
    max_queue_size: 10000
    processing_timeout: "30s"
    retry_attempts: 3
```

**Queue Monitoring**
```bash
# Check Redis queues
redis-cli LLEN insec:events:queue

# Check processing workers
ps aux | grep "insec.*worker"

# Monitor ingestion rate
curl https://api.insec.com/metrics | grep ingestion
```

**Event Validation**
```bash
# Test event submission
curl -X POST https://api.insec.com/v1/events \
  -H "Authorization: Bearer {agent_token}" \
  -H "Content-Type: application/json" \
  -d '[{"ts": "2025-08-29T10:30:00Z", "event_type": "test"}]'

# Check event validation logs
tail -f /var/log/insec/server.log | grep "event.*validation"
```

### 7. Alert Configuration Issues

#### Symptoms
- Alerts not triggering
- False positive alerts
- Missing alert notifications

#### Solutions

**Rule Configuration**
```yaml
# Validate alert rules
alerting:
  rules:
    - name: "Suspicious Login"
      condition: "event.type == 'login' AND risk_score > 0.8"
      severity: "high"
      actions: ["email", "slack"]
      cooldown: "5m"
```

**Notification Settings**
```yaml
# Check notification channels
notifications:
  email:
    enabled: true
    smtp_host: "smtp.gmail.com"
    smtp_port: 587
    from: "alerts@insec.com"
  slack:
    enabled: true
    webhook_url: "https://hooks.slack.com/services/..."
```

**Test Alert System**
```bash
# Send test alert
curl -X POST https://api.insec.com/v1/alerts/test \
  -H "Authorization: Bearer {admin_token}" \
  -d '{"message": "Test alert"}'

# Check alert logs
tail -f /var/log/insec/alerts.log
```

## 📊 Performance Troubleshooting

### Slow Queries
```sql
-- Identify slow queries
SELECT query, calls, total_time, mean_time, rows
FROM pg_stat_statements
ORDER BY mean_time DESC
LIMIT 10;

-- Check query execution plan
EXPLAIN ANALYZE SELECT * FROM events WHERE ts > '2025-08-29';

-- Add missing indexes
CREATE INDEX idx_events_ts ON events(ts);
CREATE INDEX idx_events_type ON events(event_type);
```

### High Latency
```bash
# Check network latency
ping api.insec.com

# Monitor response times
curl -w "@curl-format.txt" -o /dev/null -s https://api.insec.com/health

# Check system load
uptime
iostat -x 1 5
```

### Memory Issues
```bash
# Check memory usage
free -h
vmstat 1 5

# Check swap usage
swapon -s

# Configure swap if needed
sudo fallocate -l 2G /swapfile
sudo chmod 600 /swapfile
sudo mkswap /swapfile
sudo swapon /swapfile
echo '/swapfile none swap sw 0 0' | sudo tee -a /etc/fstab
```

## 🔧 Configuration Issues

### YAML Syntax Errors
```bash
# Validate YAML syntax
python3 -c "import yaml; yaml.safe_load(open('/etc/insec/server.yaml'))"

# Check for common issues
grep -n " " /etc/insec/server.yaml  # Mixed spaces/tabs
grep -n "\t" /etc/insec/server.yaml  # Tabs in YAML
```

### Environment Variables
```bash
# Check environment variables
env | grep INSEC

# Validate variable substitution
# server.yaml
database:
  password: "${DB_PASSWORD}"
```

### File Permissions
```bash
# Check configuration file permissions
ls -la /etc/insec/

# Fix permissions
sudo chown -R insec:insec /etc/insec/
sudo chmod 600 /etc/insec/*.key
sudo chmod 644 /etc/insec/*.crt
sudo chmod 644 /etc/insec/*.yaml
```

## 🌐 Network Troubleshooting

### DNS Issues
```bash
# Check DNS resolution
dig api.insec.com

# Check /etc/hosts
cat /etc/hosts

# Test with IP address
curl -f https://192.168.1.100/health
```

### SSL/TLS Problems
```bash
# Check certificate
openssl s_client -connect api.insec.com:443 -servername api.insec.com

# Verify certificate chain
openssl verify -CAfile /etc/ssl/certs/ca.crt /etc/ssl/certs/server.crt

# Check certificate expiration
openssl x509 -in /etc/ssl/certs/server.crt -text | grep "Not After"
```

### Firewall Issues
```bash
# Check firewall status
sudo ufw status
sudo iptables -L

# Allow INSEC ports
sudo ufw allow 8080
sudo ufw allow 443

# Check SELinux/AppArmor
sudo getenforce
sudo apparmor_status | grep insec
```

## 📝 Log Analysis

### Server Logs
```bash
# View recent logs
tail -f /var/log/insec/server.log

# Search for errors
grep "ERROR" /var/log/insec/server.log | tail -20

# Check log rotation
ls -la /var/log/insec/server.log*
```

### Agent Logs
```bash
# View agent logs
sudo tail -f /var/log/insec/agent.log

# Check agent status
sudo systemctl status insec-agent

# Restart agent with verbose logging
sudo systemctl stop insec-agent
sudo insec-agent -config /etc/insec/agent.yaml -verbose
```

### Database Logs
```bash
# PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log

# Check slow query log
sudo tail -f /var/log/postgresql/postgresql-*-slow.log
```

## 🔄 Service Management

### Systemd Issues
```bash
# Check service status
sudo systemctl status insec-server
sudo systemctl status insec-agent

# View service logs
sudo journalctl -u insec-server -f
sudo journalctl -u insec-agent -f

# Restart services
sudo systemctl restart insec-server
sudo systemctl restart insec-agent
```

### Docker Issues
```bash
# Check container status
docker ps -a | grep insec

# View container logs
docker logs insec-server
docker logs insec-agent

# Restart containers
docker-compose restart
```

### Kubernetes Issues
```bash
# Check pod status
kubectl get pods -n insec

# View pod logs
kubectl logs -n insec deployment/insec-server

# Check service endpoints
kubectl get endpoints -n insec

# Restart deployment
kubectl rollout restart deployment/insec-server -n insec
```

## 📞 Getting Help

### Support Resources
1. **Documentation**: Check this troubleshooting guide
2. **Logs**: Collect diagnostic information
3. **Community**: [GitHub Discussions](https://github.com/yashab-cyber/insec/discussions)
4. **Issues**: [GitHub Issues](https://github.com/yashab-cyber/insec/issues)

### Diagnostic Information to Collect
```bash
# System information
uname -a
lsb_release -a
free -h
df -h

# INSEC version
insec-server --version
insec-agent --version

# Configuration files
cat /etc/insec/server.yaml
cat /etc/insec/agent.yaml

# Recent logs
tail -100 /var/log/insec/server.log
tail -100 /var/log/insec/agent.log
```

### Emergency Contacts
- **Critical Issues**: +1-800-INSEC-EM (46732)
- **Email Support**: support@insec.com
- **Security Issues**: security@insec.com

---

**Last updated:** August 29, 2025</content>
<parameter name="filePath">/workspaces/insec/docs/troubleshooting.md
