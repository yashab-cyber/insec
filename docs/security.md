# INSEC Security Guide

This comprehensive guide covers security features, best practices, and hardening procedures for INSEC deployments.

## 🔒 Security Architecture

### Core Security Principles

INSEC implements a defense-in-depth security model with multiple layers of protection:

1. **Network Security**: TLS encryption, firewall rules, network segmentation
2. **Authentication**: Multi-factor authentication, JWT tokens, certificate-based auth
3. **Authorization**: Role-based access control, fine-grained permissions
4. **Data Protection**: Encryption at rest, data masking, retention policies
5. **Monitoring**: Comprehensive audit logging, real-time alerting
6. **Compliance**: GDPR, HIPAA, SOC 2 compliance features

## 🛡️ Authentication & Authorization

### Multi-Factor Authentication (MFA)

#### Setup MFA for Users
```bash
# Enable MFA for a user
curl -X POST https://api.insec.com/v1/users/{user_id}/mfa/enable \
  -H "Authorization: Bearer {token}" \
  -H "Content-Type: application/json" \
  -d '{
    "method": "totp",
    "phone": "+1234567890"
  }'
```

#### MFA Configuration
```yaml
# server.yaml
security:
  mfa:
    enabled: true
    required: true
    methods:
      - totp
      - sms
      - email
    grace_period: "7d"
```

### Role-Based Access Control (RBAC)

#### Built-in Roles
- **Super Admin**: Full system access
- **Security Admin**: Security policy management
- **Analyst**: Event investigation and reporting
- **Operator**: Alert management and response
- **Auditor**: Read-only access to all data
- **Agent**: Limited API access for telemetry

#### Custom Roles
```json
{
  "name": "Custom Analyst",
  "permissions": [
    "events:read",
    "alerts:read",
    "alerts:update",
    "reports:create"
  ],
  "restrictions": {
    "max_query_time": "30d",
    "allowed_ips": ["192.168.1.0/24"]
  }
}
```

### Certificate-Based Authentication

#### Agent Certificate Generation
```bash
# Create agent certificate
openssl genrsa -out agent.key 2048
openssl req -new -key agent.key -out agent.csr \
  -subj "/CN=agent-001/O=INSEC/C=US"

# Sign certificate
openssl x509 -req -in agent.csr \
  -CA ca.crt -CAkey ca.key \
  -out agent.crt -days 365

# Convert to PKCS#12
openssl pkcs12 -export -out agent.p12 \
  -inkey agent.key -in agent.crt -certfile ca.crt
```

## 🔐 Data Protection

### Encryption at Rest

#### Database Encryption
```sql
-- Enable encryption for PostgreSQL
CREATE EXTENSION pgcrypto;

-- Encrypt sensitive data
UPDATE users
SET encrypted_data = pgp_sym_encrypt('sensitive_data', 'encryption_key')
WHERE id = 1;
```

#### File System Encryption
```bash
# Encrypt sensitive files
gpg --encrypt --recipient security@insec.com sensitive_file.txt

# Decrypt when needed
gpg --decrypt sensitive_file.txt.gpg > sensitive_file.txt
```

### Data Masking & Redaction

#### Configuration
```yaml
data_protection:
  masking:
    enabled: true
    rules:
      - pattern: "\\b\\d{4}-\\d{4}-\\d{4}-\\d{4}\\b"
        replacement: "XXXX-XXXX-XXXX-XXXX"
      - pattern: "\\b[A-Za-z0-9._%+-]+@[A-Za-z0-9.-]+\\.[A-Z|a-z]{2,}\\b"
        replacement: "user@domain.com"
```

#### Custom Masking Rules
```json
{
  "name": "credit_card_masking",
  "pattern": "\\b\\d{13,19}\\b",
  "replacement": "****************",
  "context": ["payment_data", "transaction_logs"]
}
```

### Data Retention Policies

#### Retention Configuration
```yaml
retention:
  policies:
    - name: "security_events"
      retention_days: 365
      compression: true
      archive_to: "s3://insec-archive/security-events/"
    - name: "audit_logs"
      retention_days: 2555  # 7 years
      encryption: true
      immutable: true
    - name: "debug_logs"
      retention_days: 30
      auto_delete: true
```

## 🌐 Network Security

### TLS Configuration

#### Strong TLS Settings
```nginx
# nginx.conf
server {
    listen 443 ssl http2;
    server_name api.insec.com;

    # SSL Configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512;
    ssl_prefer_server_ciphers off;
    ssl_session_cache shared:SSL:10m;
    ssl_session_timeout 10m;

    # HSTS
    add_header Strict-Transport-Security "max-age=31536000; includeSubDomains" always;

    # Certificate pinning
    add_header Public-Key-Pins 'pin-sha256="..."; max-age=2592000' always;
}
```

#### Certificate Management
```bash
# Check certificate expiration
openssl x509 -in certificate.crt -text -noout | grep "Not After"

# Renew certificate
certbot renew --cert-name api.insec.com

# Validate certificate chain
openssl verify -CAfile ca.crt certificate.crt
```

### Firewall Configuration

#### iptables Rules
```bash
# Allow SSH only from specific IPs
iptables -A INPUT -p tcp -s 192.168.1.0/24 --dport 22 -j ACCEPT
iptables -A INPUT -p tcp --dport 22 -j DROP

# Allow HTTPS
iptables -A INPUT -p tcp --dport 443 -j ACCEPT

# Allow INSEC API
iptables -A INPUT -p tcp --dport 8080 -j ACCEPT

# Default deny
iptables -P INPUT DROP
iptables -P FORWARD DROP
iptables -P OUTPUT ACCEPT
```

#### UFW Rules
```bash
# Enable UFW
sudo ufw enable

# Allow specific ports
sudo ufw allow from 192.168.1.0/24 to any port 22
sudo ufw allow 443
sudo ufw allow 8080

# Rate limiting
sudo ufw limit ssh
```

### Network Segmentation

#### VLAN Configuration
```bash
# Create VLAN for INSEC servers
vconfig add eth0 100
ifconfig eth0.100 192.168.100.1 netmask 255.255.255.0

# Configure routing
ip route add 192.168.100.0/24 dev eth0.100
```

## 🏢 Agent Security

### Agent Hardening

#### Secure Agent Installation
```bash
# Install in protected directory
sudo mkdir -p /opt/insec/agent
sudo chown root:root /opt/insec/agent
sudo chmod 755 /opt/insec/agent

# Run agent as non-privileged user
sudo useradd --system --shell /bin/false insec-agent
sudo chown insec-agent:insec-agent /opt/insec/agent
```

#### Agent Self-Protection
```yaml
agent:
  security:
    self_protection:
      enabled: true
      prevent_deletion: true
      prevent_modification: true
      integrity_check: true
      tamper_detection: true
```

### Secure Communication

#### mTLS Configuration
```yaml
agent:
  tls:
    enabled: true
    ca_cert: "/etc/ssl/certs/ca.crt"
    client_cert: "/etc/ssl/certs/agent.crt"
    client_key: "/etc/ssl/private/agent.key"
    server_name_verification: true
    certificate_revocation_check: true
```

## 📊 Monitoring & Auditing

### Audit Logging

#### Audit Configuration
```yaml
audit:
  enabled: true
  log_level: "detailed"
  destinations:
    - file: "/var/log/insec/audit.log"
    - syslog: "local6"
    - remote: "https://audit.insec.com/v1/events"
  events:
    - authentication
    - authorization
    - data_access
    - configuration_changes
    - security_events
```

#### Audit Event Format
```json
{
  "timestamp": "2025-08-29T10:30:00Z",
  "event_type": "authentication",
  "user_id": "user_123",
  "ip_address": "192.168.1.100",
  "user_agent": "INSEC-Console/1.0",
  "action": "login",
  "result": "success",
  "details": {
    "mfa_used": true,
    "method": "totp"
  }
}
```

### Security Monitoring

#### SIEM Integration
```yaml
integrations:
  siem:
    splunk:
      enabled: true
      host: "splunk.company.com"
      port: 8088
      token: "your-splunk-token"
      index: "insec_security"
    elastic:
      enabled: true
      host: "elasticsearch.company.com"
      port: 9200
      index: "insec-audit"
```

#### Real-time Alerts
```yaml
alerts:
  security:
    - name: "Brute Force Attack"
      condition: "failed_login_attempts > 5 within 5m"
      severity: "high"
      actions: ["block_ip", "notify_security"]
    - name: "Privilege Escalation"
      condition: "user_role_changed AND old_role != 'admin'"
      severity: "critical"
      actions: ["alert_security", "revoke_session"]
```

## 🔑 Key Management

### Encryption Key Management

#### Key Rotation
```bash
# Rotate database encryption key
insec-cli keys rotate --type database --grace-period 24h

# Rotate JWT signing key
insec-cli keys rotate --type jwt --immediate

# List current keys
insec-cli keys list
```

#### Key Backup and Recovery
```bash
# Backup keys
insec-cli keys backup --output /secure/backup/keys.enc

# Restore keys
insec-cli keys restore --input /secure/backup/keys.enc

# Key recovery procedure
insec-cli keys recover --token recovery_token_123
```

### Hardware Security Modules (HSM)

#### HSM Integration
```yaml
hsm:
  enabled: true
  provider: "aws-cloudhsm"
  cluster_id: "cluster-12345"
  key_rotation:
    enabled: true
    interval: "90d"
  backup:
    enabled: true
    interval: "7d"
```

## 🛠️ Security Hardening

### Server Hardening

#### System Hardening Script
```bash
#!/bin/bash
# INSEC Server Hardening Script

# Disable unused services
systemctl disable avahi-daemon
systemctl disable cups

# Configure SSH
sed -i 's/#PermitRootLogin yes/PermitRootLogin no/' /etc/ssh/sshd_config
sed -i 's/#PasswordAuthentication yes/PasswordAuthentication no/' /etc/ssh/sshd_config
systemctl reload sshd

# Configure firewall
ufw enable
ufw allow ssh
ufw allow 443

# Install security updates
apt update && apt upgrade -y

# Configure auditd
systemctl enable auditd
auditctl -e 1
```

#### Kernel Security
```bash
# Enable security modules
echo "kernel.kptr_restrict = 1" >> /etc/sysctl.conf
echo "kernel.dmesg_restrict = 1" >> /etc/sysctl.conf
echo "net.core.bpf_jit_harden = 2" >> /etc/sysctl.conf

# Apply settings
sysctl -p
```

### Application Hardening

#### Security Headers
```nginx
# nginx.conf
add_header X-Frame-Options "SAMEORIGIN" always;
add_header X-Content-Type-Options "nosniff" always;
add_header X-XSS-Protection "1; mode=block" always;
add_header Referrer-Policy "strict-origin-when-cross-origin" always;
add_header Content-Security-Policy "default-src 'self'" always;
```

#### Rate Limiting
```yaml
rate_limiting:
  enabled: true
  global:
    requests_per_minute: 1000
    burst_limit: 100
  endpoints:
    "/auth/login":
      requests_per_minute: 5
      burst_limit: 2
    "/api/v1/events":
      requests_per_minute: 10000
      burst_limit: 1000
```

## 📋 Compliance

### GDPR Compliance

#### Data Subject Rights
```yaml
gdpr:
  enabled: true
  data_retention:
    personal_data: "2555d"  # 7 years
    logs: "365d"
  consent_management:
    enabled: true
    consent_types: ["analytics", "marketing"]
  data_portability:
    enabled: true
    export_formats: ["json", "csv"]
```

#### Data Processing Agreement
```json
{
  "processor": "INSEC",
  "controller": "Customer Company",
  "data_categories": ["personal_data", "usage_data"],
  "processing_purposes": ["security_monitoring", "threat_detection"],
  "retention_period": "7_years",
  "security_measures": ["encryption", "access_control", "audit_logging"]
}
```

### SOC 2 Compliance

#### Trust Criteria Implementation
```yaml
soc2:
  enabled: true
  controls:
    security:
      - access_control
      - encryption
      - monitoring
    availability:
      - redundancy
      - backup
      - disaster_recovery
    integrity:
      - data_validation
      - error_handling
      - audit_trails
    confidentiality:
      - data_classification
      - access_restrictions
      - encryption
    privacy:
      - consent_management
      - data_minimization
      - retention_policies
```

## 🚨 Incident Response

### Security Incident Response Plan

#### Incident Classification
- **Level 1**: Minor security event, no immediate threat
- **Level 2**: Potential security breach, investigation needed
- **Level 3**: Confirmed breach, immediate response required
- **Level 4**: Major breach, full incident response team activated

#### Response Procedures
```yaml
incident_response:
  level1:
    response_time: "4h"
    notification: ["security_team"]
    actions: ["investigate", "document"]
  level2:
    response_time: "1h"
    notification: ["security_team", "management"]
    actions: ["investigate", "contain", "notify_affected"]
  level3:
    response_time: "30m"
    notification: ["security_team", "management", "legal"]
    actions: ["contain", "eradicate", "recover", "notify_regulators"]
```

### Breach Notification
```json
{
  "incident_id": "INC-2025-001",
  "breach_date": "2025-08-29T10:30:00Z",
  "affected_users": 150,
  "data_compromised": ["emails", "names"],
  "risk_assessment": "medium",
  "containment_actions": ["password_reset", "account_lockout"],
  "notification_date": "2025-08-29T14:00:00Z"
}
```

## 🔍 Security Testing

### Penetration Testing

#### Automated Security Scanning
```bash
# Run vulnerability scan
nessus -T html -o report.html -i insec_targets.txt

# Web application scanning
owasp-zap -cmd -quickurl https://api.insec.com -quickout zap_report.html

# Container scanning
trivy image yashab/insec-server:latest
```

#### Manual Testing Checklist
- [ ] Authentication bypass attempts
- [ ] Authorization escalation
- [ ] SQL injection testing
- [ ] XSS vulnerability testing
- [ ] CSRF protection verification
- [ ] SSL/TLS configuration review
- [ ] Session management testing
- [ ] File upload security testing

### Security Audits

#### Internal Audit Procedure
```yaml
audit:
  schedule: "quarterly"
  scope:
    - access_controls
    - data_protection
    - network_security
    - incident_response
  deliverables:
    - audit_report
    - findings_remediation
    - compliance_status
```

## 📞 Security Support

### Security Contacts
- **Security Team**: security@insec.com
- **Emergency**: +1-800-INSEC-EM (46732)
- **PGP Key**: [Download](https://insec.com/security/pgp-key.asc)

### Vulnerability Disclosure
- **Responsible Disclosure**: security@insec.com
- **Bug Bounty**: [Program Details](https://insec.com/bug-bounty)
- **Response Time**: 48 hours for critical issues

### Security Updates
- **Security Advisories**: [RSS Feed](https://insec.com/security/advisories.xml)
- **Patch Management**: Automatic updates available
- **Emergency Patches**: Within 24 hours for critical vulnerabilities

---

**Last updated:** August 29, 2025

**🔒 This document contains sensitive security information. Handle with care.**</content>
<parameter name="filePath">/workspaces/insec/docs/security.md
