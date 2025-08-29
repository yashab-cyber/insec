# INSEC API Reference

This document provides comprehensive API documentation for INSEC's RESTful APIs.

## 🌐 Base URL

```
https://api.insec.com/v1
```

All API requests must use HTTPS and include proper authentication headers.

## 🔐 Authentication

INSEC uses JWT (JSON Web Tokens) for API authentication.

### Headers Required
```
Authorization: Bearer <jwt_token>
X-Tenant-ID: <tenant_id>
Content-Type: application/json
```

### Authentication Endpoints

#### Login
```http
POST /auth/login
```

**Request Body:**
```json
{
  "username": "admin@insec.com",
  "password": "secure_password",
  "mfa_code": "123456"
}
```

**Response:**
```json
{
  "token": "eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9...",
  "expires_at": "2025-08-29T15:30:00Z",
  "user": {
    "id": "user_123",
    "email": "admin@insec.com",
    "role": "admin"
  }
}
```

#### Refresh Token
```http
POST /auth/refresh
```

**Request Body:**
```json
{
  "refresh_token": "refresh_token_here"
}
```

## 📊 Events API

### Ingest Events
```http
POST /events
```

**Description:** Ingest security events from agents.

**Request Body:**
```json
[
  {
    "ts": "2025-08-29T10:30:00Z",
    "tenant_id": "tenant_123",
    "host_id": "host_456",
    "user": {
      "id": "user_789",
      "email": "user@company.com",
      "dept": "engineering"
    },
    "os": {
      "family": "linux",
      "version": "Ubuntu 22.04"
    },
    "event": {
      "type": "process",
      "id": "evt_001"
    },
    "proc": {
      "name": "curl",
      "cmd": ["curl", "https://example.com"],
      "ppid": 1234,
      "hash": "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"
    },
    "labels": ["suspicious", "network"],
    "risk_hints": ["unusual_command", "external_connection"],
    "agent": {
      "ver": "1.0.0",
      "mode": "enforce"
    }
  }
]
```

**Response:**
```json
{
  "status": "ok",
  "count": 1,
  "events_ingested": ["evt_001"]
}
```

### Query Events
```http
GET /events
```

**Query Parameters:**
- `start_time`: ISO8601 timestamp
- `end_time`: ISO8601 timestamp
- `event_type`: process|file|network|user
- `host_id`: specific host ID
- `user_id`: specific user ID
- `limit`: maximum results (default: 100, max: 1000)
- `offset`: pagination offset

**Example:**
```http
GET /events?event_type=process&start_time=2025-08-29T00:00:00Z&limit=50
```

**Response:**
```json
{
  "events": [
    {
      "id": "evt_001",
      "ts": "2025-08-29T10:30:00Z",
      "event_type": "process",
      "host_id": "host_456",
      "user_id": "user_789",
      "data": {
        "process_name": "curl",
        "command_line": "curl https://example.com"
      },
      "risk_score": 0.8,
      "labels": ["suspicious"]
    }
  ],
  "total": 1,
  "has_more": false
}
```

### Get Event Details
```http
GET /events/{event_id}
```

**Response:**
```json
{
  "id": "evt_001",
  "ts": "2025-08-29T10:30:00Z",
  "tenant_id": "tenant_123",
  "host_id": "host_456",
  "event_type": "process",
  "data": {
    "process": {
      "name": "curl",
      "pid": 1234,
      "cmdline": ["curl", "https://example.com"],
      "hash": "a665a45920422f9d417e4867efdc4fb8a04a1f3fff1fa07e998e86f7f7a27ae3"
    }
  },
  "metadata": {
    "risk_score": 0.8,
    "correlation_id": "corr_123",
    "labels": ["suspicious", "network"],
    "risk_hints": ["unusual_command", "external_connection"]
  },
  "agent": {
    "version": "1.0.0",
    "mode": "enforce"
  }
}
```

## 🚨 Alerts API

### Get Alerts
```http
GET /alerts
```

**Query Parameters:**
- `status`: open|closed|acknowledged
- `severity`: low|medium|high|critical
- `start_time`: ISO8601 timestamp
- `end_time`: ISO8601 timestamp
- `limit`: maximum results (default: 50)

**Response:**
```json
{
  "alerts": [
    {
      "id": "alert_001",
      "title": "Suspicious Process Execution",
      "description": "Process 'curl' executed with suspicious parameters",
      "severity": "high",
      "status": "open",
      "created_at": "2025-08-29T10:30:00Z",
      "updated_at": "2025-08-29T10:30:00Z",
      "events": ["evt_001", "evt_002"],
      "assignee": "analyst@company.com",
      "tags": ["process", "network"]
    }
  ],
  "total": 1
}
```

### Create Alert
```http
POST /alerts
```

**Request Body:**
```json
{
  "title": "Custom Security Alert",
  "description": "Manually created alert for investigation",
  "severity": "medium",
  "events": ["evt_001"],
  "tags": ["manual", "investigation"]
}
```

### Update Alert
```http
PUT /alerts/{alert_id}
```

**Request Body:**
```json
{
  "status": "acknowledged",
  "assignee": "analyst@company.com",
  "notes": "Investigating suspicious activity"
}
```

### Execute Alert Action
```http
POST /alerts/{alert_id}/actions
```

**Request Body:**
```json
{
  "action": "isolate_host",
  "parameters": {
    "duration": 3600,
    "reason": "Suspicious activity detected"
  }
}
```

## 📋 Policies API

### Get Policies
```http
GET /policies
```

**Response:**
```json
{
  "policies": [
    {
      "id": "policy_001",
      "name": "Data Exfiltration Prevention",
      "description": "Prevent unauthorized data transfers",
      "enabled": true,
      "rules": [
        {
          "id": "rule_001",
          "condition": "file.size > 100MB AND file.destination NOT IN allowed_domains",
          "action": "block",
          "severity": "high"
        }
      ],
      "created_at": "2025-08-29T09:00:00Z",
      "updated_at": "2025-08-29T09:00:00Z"
    }
  ]
}
```

### Create Policy
```http
POST /policies
```

**Request Body:**
```json
{
  "name": "New Security Policy",
  "description": "Policy description",
  "enabled": true,
  "rules": [
    {
      "condition": "process.name == 'suspicious.exe'",
      "action": "alert",
      "severity": "high",
      "parameters": {
        "notify_channels": ["email", "slack"]
      }
    }
  ]
}
```

### Update Policy
```http
PUT /policies/{policy_id}
```

**Request Body:**
```json
{
  "enabled": false,
  "rules": [
    {
      "id": "rule_001",
      "action": "quarantine"
    }
  ]
}
```

### Delete Policy
```http
DELETE /policies/{policy_id}
```

## 👥 Users API

### Get Users
```http
GET /users
```

**Response:**
```json
{
  "users": [
    {
      "id": "user_123",
      "email": "user@company.com",
      "name": "John Doe",
      "role": "analyst",
      "department": "security",
      "last_login": "2025-08-29T08:30:00Z",
      "status": "active"
    }
  ]
}
```

### Create User
```http
POST /users
```

**Request Body:**
```json
{
  "email": "newuser@company.com",
  "name": "Jane Smith",
  "role": "analyst",
  "department": "security",
  "send_invite": true
}
```

### Update User
```http
PUT /users/{user_id}
```

**Request Body:**
```json
{
  "role": "admin",
  "department": "it"
}
```

## 🏢 Tenants API

### Get Tenants
```http
GET /tenants
```

**Response:**
```json
{
  "tenants": [
    {
      "id": "tenant_123",
      "name": "Acme Corp",
      "domain": "acme.com",
      "status": "active",
      "created_at": "2025-01-01T00:00:00Z",
      "settings": {
        "retention_days": 365,
        "max_agents": 1000,
        "features": ["ueba", "automated_response"]
      }
    }
  ]
}
```

### Create Tenant
```http
POST /tenants
```

**Request Body:**
```json
{
  "name": "New Company",
  "domain": "newcompany.com",
  "admin_email": "admin@newcompany.com",
  "settings": {
    "retention_days": 90,
    "max_agents": 100
  }
}
```

## 📊 Analytics API

### Get Dashboard Data
```http
GET /analytics/dashboard
```

**Query Parameters:**
- `time_range`: 1h|24h|7d|30d
- `group_by`: hour|day|week

**Response:**
```json
{
  "summary": {
    "total_events": 15420,
    "active_alerts": 12,
    "risk_score_avg": 0.3,
    "agents_online": 98
  },
  "charts": {
    "events_over_time": [
      {"timestamp": "2025-08-29T00:00:00Z", "count": 450},
      {"timestamp": "2025-08-29T01:00:00Z", "count": 380}
    ],
    "alerts_by_severity": {
      "low": 25,
      "medium": 15,
      "high": 8,
      "critical": 2
    }
  }
}
```

### Get Risk Analytics
```http
GET /analytics/risk
```

**Response:**
```json
{
  "user_risk_scores": [
    {"user_id": "user_123", "risk_score": 0.8, "factors": ["unusual_login", "data_access"]},
    {"user_id": "user_456", "risk_score": 0.2, "factors": []}
  ],
  "entity_risk_scores": [
    {"entity_id": "host_789", "risk_score": 0.6, "factors": ["suspicious_processes"]}
  ]
}
```

## 🔍 Search API

### Search Events
```http
POST /search/events
```

**Request Body:**
```json
{
  "query": "process.name:powershell AND file.size:>100MB",
  "start_time": "2025-08-29T00:00:00Z",
  "end_time": "2025-08-29T23:59:59Z",
  "limit": 100,
  "sort": {
    "field": "timestamp",
    "order": "desc"
  }
}
```

**Response:**
```json
{
  "results": [
    {
      "id": "evt_001",
      "timestamp": "2025-08-29T10:30:00Z",
      "event_type": "file",
      "data": {
        "filename": "large_file.zip",
        "size": 150000000,
        "destination": "external_site.com"
      },
      "risk_score": 0.9
    }
  ],
  "total": 1,
  "took": 45
}
```

## 🏥 Health API

### System Health
```http
GET /health
```

**Response:**
```json
{
  "status": "healthy",
  "version": "1.0.0",
  "timestamp": "2025-08-29T10:30:00Z",
  "services": {
    "database": "healthy",
    "redis": "healthy",
    "ingestion": "healthy",
    "processing": "healthy"
  },
  "metrics": {
    "uptime": "24h 30m",
    "memory_usage": "2.1GB",
    "cpu_usage": "15%"
  }
}
```

### Readiness Probe
```http
GET /health/ready
```

**Response:**
```json
{
  "status": "ready",
  "checks": {
    "database": true,
    "redis": true,
    "dependencies": true
  }
}
```

## 📋 Rate Limits

- **Authenticated Requests**: 1000 per minute per user
- **Anonymous Requests**: 100 per hour per IP
- **Event Ingestion**: 10000 per minute per tenant
- **Search Requests**: 100 per minute per user

Rate limit headers are included in responses:
```
X-RateLimit-Limit: 1000
X-RateLimit-Remaining: 999
X-RateLimit-Reset: 1630342800
```

## 🚨 Error Responses

All errors follow a consistent format:

```json
{
  "error": {
    "code": "VALIDATION_ERROR",
    "message": "Invalid request parameters",
    "details": {
      "field": "email",
      "reason": "must be a valid email address"
    }
  },
  "request_id": "req_123456"
}
```

### Common Error Codes
- `VALIDATION_ERROR`: Invalid request data
- `AUTHENTICATION_ERROR`: Invalid or missing credentials
- `AUTHORIZATION_ERROR`: Insufficient permissions
- `NOT_FOUND`: Resource not found
- `RATE_LIMITED`: Too many requests
- `INTERNAL_ERROR`: Server error

## 📖 Pagination

List endpoints support pagination:

```json
{
  "data": [...],
  "pagination": {
    "page": 1,
    "limit": 50,
    "total": 1250,
    "has_more": true,
    "next_cursor": "cursor_123"
  }
}
```

Use `cursor` parameter for subsequent requests:
```http
GET /events?cursor=cursor_123&limit=50
```

## 🔗 Webhooks

Configure webhooks for real-time notifications:

```http
POST /webhooks
```

**Request Body:**
```json
{
  "url": "https://your-app.com/webhook",
  "events": ["alert.created", "alert.updated"],
  "secret": "webhook_secret",
  "enabled": true
}
```

Webhook payload:
```json
{
  "event": "alert.created",
  "data": {
    "alert": {
      "id": "alert_001",
      "title": "Suspicious Activity",
      "severity": "high"
    }
  },
  "timestamp": "2025-08-29T10:30:00Z",
  "signature": "sha256=..."
}
```

## 📚 SDKs and Libraries

### Official SDKs
- **Go SDK**: `go get github.com/yashab-cyber/insec-go`
- **Python SDK**: `pip install insec-python`
- **JavaScript SDK**: `npm install @insec/js-sdk`

### Community Libraries
- **Ruby**: `gem install insec-ruby`
- **Java**: Maven/Gradle packages available
- **.NET**: NuGet package available

## 🔄 API Versions

- **v1** (Current): Stable production API
- **v2** (Beta): Next-generation API with enhanced features

Use version headers for specific versions:
```
Accept: application/vnd.insec.v2+json
```

## 📞 Support

For API support:
- **Documentation**: This reference guide
- **Issues**: [GitHub Issues](https://github.com/yashab-cyber/insec/issues)
- **Discussions**: [GitHub Discussions](https://github.com/yashab-cyber/insec/discussions)
- **Email**: api-support@insec.com

---

**Last updated:** August 29, 2025</content>
<parameter name="filePath">/workspaces/insec/docs/api-reference.md
