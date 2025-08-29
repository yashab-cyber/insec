# INSEC Installation Guide

This guide provides step-by-step instructions for installing and setting up INSEC components.

## 📋 Prerequisites

### System Requirements

#### Minimum Requirements
- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 20GB free space
- **Network**: 100Mbps connection

#### Recommended Requirements
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Storage**: 100GB+ SSD
- **Network**: 1Gbps connection

### Operating System Support

#### Agent
- ✅ Windows 10/11 (64-bit)
- ✅ macOS 10.15+ (Intel/Apple Silicon)
- ✅ Ubuntu 18.04+ / CentOS 7+ / RHEL 7+
- ✅ Debian 10+

#### Server
- ✅ Ubuntu 20.04+ (recommended)
- ✅ CentOS 8+ / RHEL 8+
- ✅ Docker containers (any OS)

#### Console
- ✅ Modern web browsers (Chrome, Firefox, Safari, Edge)
- ✅ Node.js 18+ for development

### Dependencies

#### Required Software
- **Rust**: 1.70+ (for agent compilation)
- **Go**: 1.19+ (for server compilation)
- **Node.js**: 18+ (for UI development)
- **PostgreSQL**: 13+ (for data storage)
- **Redis**: 6+ (for caching)

#### Optional Dependencies
- **Docker**: 20+ (for containerized deployment)
- **Kubernetes**: 1.24+ (for orchestration)
- **Nginx**: For reverse proxy
- **Certbot**: For SSL certificates

## 🚀 Quick Installation

### Option 1: Docker Compose (Recommended)

```bash
# Clone the repository
git clone https://github.com/yashab-cyber/insec.git
cd insec

# Start all services
docker-compose up -d

# Access the console at http://localhost:3000
```

### Option 2: Manual Installation

```bash
# Clone the repository
git clone https://github.com/yashab-cyber/insec.git
cd insec

# Install dependencies
./scripts/install-dependencies.sh

# Build all components
./scripts/build.sh

# Start services
./scripts/start.sh
```

## 📦 Detailed Installation

### 1. Install System Dependencies

#### Ubuntu/Debian
```bash
sudo apt update
sudo apt install -y curl wget git build-essential postgresql redis-server nginx

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs
```

#### CentOS/RHEL
```bash
sudo yum update -y
sudo yum install -y curl wget git gcc postgresql-server redis nginx

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Node.js
curl -fsSL https://rpm.nodesource.com/setup_18.x | sudo bash -
sudo yum install -y nodejs
```

#### macOS
```bash
# Install Homebrew if not installed
/bin/bash -c "$(curl -fsSL https://raw.githubusercontent.com/Homebrew/install/HEAD/install.sh)"

# Install dependencies
brew install rust go node postgresql redis nginx git

# Start services
brew services start postgresql
brew services start redis
```

#### Windows
```powershell
# Install Chocolatey if not installed
Set-ExecutionPolicy Bypass -Scope Process -Force; [System.Net.ServicePointManager]::SecurityProtocol = [System.Net.ServicePointManager]::SecurityProtocol -bor 3072; iex ((New-Object System.Net.WebClient).DownloadString('https://chocolatey.org/install.ps1'))

# Install dependencies
choco install -y rust go nodejs postgresql redis-64 git

# Install PostgreSQL and Redis as services
# Follow the installation wizards
```

### 2. Database Setup

#### PostgreSQL Configuration
```bash
# Initialize PostgreSQL (if not already done)
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Create database and user
sudo -u postgres psql
```

```sql
CREATE DATABASE insec;
CREATE USER insec_user WITH ENCRYPTED PASSWORD 'your_secure_password';
GRANT ALL PRIVILEGES ON DATABASE insec TO insec_user;
\q
```

#### Redis Configuration
```bash
# Configure Redis
sudo systemctl enable redis
sudo systemctl start redis

# Test Redis connection
redis-cli ping
```

### 3. Build INSEC Components

#### Build Agent
```bash
cd agent
cargo build --release
```

#### Build Server
```bash
cd ../server
go mod download
go build -o insec-server .
```

#### Build UI
```bash
cd ../ui
npm install
npm run build
```

### 4. Configuration

#### Server Configuration
Create `/etc/insec/server.yaml`:

```yaml
server:
  host: "0.0.0.0"
  port: 8080
  tls:
    enabled: true
    cert_file: "/etc/ssl/certs/insec.crt"
    key_file: "/etc/ssl/private/insec.key"

database:
  host: "localhost"
  port: 5432
  database: "insec"
  username: "insec_user"
  password: "your_secure_password"

redis:
  host: "localhost"
  port: 6379
  password: ""

security:
  jwt_secret: "your-jwt-secret-key"
  mTls:
    ca_cert: "/etc/ssl/certs/ca.crt"
    server_cert: "/etc/ssl/certs/server.crt"
    server_key: "/etc/ssl/private/server.key"
```

#### Agent Configuration
Create `/etc/insec/agent.yaml`:

```yaml
server:
  url: "https://your-server:8080"
  tls:
    ca_cert: "/etc/ssl/certs/ca.crt"
    client_cert: "/etc/ssl/certs/agent.crt"
    client_key: "/etc/ssl/private/agent.key"

agent:
  id: "unique-agent-id"
  tenant_id: "your-tenant-id"
  collection_interval: 30
  max_queue_size: 1000

policies:
  enabled: true
  update_interval: 300

logging:
  level: "info"
  file: "/var/log/insec/agent.log"
```

### 5. SSL/TLS Setup

#### Generate Self-Signed Certificates (Development)
```bash
# Create CA
openssl genrsa -out ca.key 4096
openssl req -new -x509 -days 365 -key ca.key -sha256 -out ca.crt

# Create server certificate
openssl genrsa -out server.key 2048
openssl req -subj "/CN=localhost" -new -key server.key -out server.csr
openssl x509 -req -days 365 -in server.csr -CA ca.crt -CAkey ca.key -out server.crt

# Create agent certificate
openssl genrsa -out agent.key 2048
openssl req -subj "/CN=agent" -new -key agent.key -out agent.csr
openssl x509 -req -days 365 -in agent.csr -CA ca.crt -CAkey ca.key -out agent.crt
```

#### Production SSL (Let's Encrypt)
```bash
# Install Certbot
sudo apt install certbot

# Get certificate
sudo certbot certonly --standalone -d your-domain.com

# Configure automatic renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### 6. Service Setup

#### Systemd Service for Server
Create `/etc/systemd/system/insec-server.service`:

```ini
[Unit]
Description=INSEC Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=insec
Group=insec
ExecStart=/usr/local/bin/insec-server -config /etc/insec/server.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Systemd Service for Agent
Create `/etc/systemd/system/insec-agent.service`:

```ini
[Unit]
Description=INSEC Agent
After=network.target

[Service]
Type=simple
User=root
ExecStart=/usr/local/bin/insec-agent -config /etc/insec/agent.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

#### Enable and Start Services
```bash
# Reload systemd
sudo systemctl daemon-reload

# Enable services
sudo systemctl enable insec-server
sudo systemctl enable insec-agent

# Start services
sudo systemctl start insec-server
sudo systemctl start insec-agent

# Check status
sudo systemctl status insec-server
sudo systemctl status insec-agent
```

### 7. Nginx Reverse Proxy (Optional)

Create `/etc/nginx/sites-available/insec`:

```nginx
server {
    listen 80;
    server_name your-domain.com;

    # Redirect HTTP to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name your-domain.com;

    ssl_certificate /etc/letsencrypt/live/your-domain.com/fullchain.pem;
    ssl_certificate_key /etc/letsencrypt/live/your-domain.com/privkey.pem;

    # SSL configuration
    ssl_protocols TLSv1.2 TLSv1.3;
    ssl_ciphers ECDHE-RSA-AES256-GCM-SHA512:DHE-RSA-AES256-GCM-SHA512:ECDHE-RSA-AES256-GCM-SHA384;
    ssl_prefer_server_ciphers off;

    # API endpoints
    location /api/ {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }

    # Web console
    location / {
        proxy_pass http://localhost:3000;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}
```

Enable the site:
```bash
sudo ln -s /etc/nginx/sites-available/insec /etc/nginx/sites-enabled/
sudo nginx -t
sudo systemctl reload nginx
```

## 🔍 Verification

### Check Services
```bash
# Check server
curl -k https://localhost:8080/health

# Check agent
sudo systemctl status insec-agent

# Check logs
sudo journalctl -u insec-server -f
sudo journalctl -u insec-agent -f
```

### Test Installation
```bash
# Test API
curl -k https://localhost:8080/api/v1/events

# Test console
curl -k https://localhost:3000
```

## 🐳 Docker Deployment

### Docker Compose Setup
Create `docker-compose.yml`:

```yaml
version: '3.8'

services:
  postgres:
    image: postgres:15
    environment:
      POSTGRES_DB: insec
      POSTGRES_USER: insec
      POSTGRES_PASSWORD: password
    volumes:
      - postgres_data:/var/lib/postgresql/data

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data

  server:
    build: ./server
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    environment:
      - DATABASE_URL=postgres://insec:password@postgres:5432/insec
      - REDIS_URL=redis://redis:6379

  ui:
    build: ./ui
    ports:
      - "3000:3000"
    depends_on:
      - server

volumes:
  postgres_data:
  redis_data:
```

### Build and Run
```bash
docker-compose up -d
```

## 🔧 Troubleshooting

### Common Issues

#### Agent Won't Start
```bash
# Check configuration
sudo /usr/local/bin/insec-agent -config /etc/insec/agent.yaml -validate

# Check permissions
ls -la /etc/insec/
ls -la /var/log/insec/
```

#### Server Connection Issues
```bash
# Test network connectivity
telnet localhost 8080

# Check firewall
sudo ufw status
sudo iptables -L
```

#### Database Connection Issues
```bash
# Test database connection
psql -h localhost -U insec_user -d insec

# Check PostgreSQL logs
sudo tail -f /var/log/postgresql/postgresql-*.log
```

## 📞 Support

If you encounter issues during installation:

1. Check the [Troubleshooting Guide](troubleshooting.md)
2. Review logs in `/var/log/insec/`
3. Open an issue on [GitHub](https://github.com/yashab-cyber/insec/issues)
4. Contact support at yashabalam707@gmail.com

## 🎉 Next Steps

Once installation is complete:

1. **Configure Policies**: Set up your security policies
2. **Deploy Agents**: Install agents on endpoints
3. **Set up Monitoring**: Configure alerting and dashboards
4. **Test Detection**: Verify threat detection capabilities

Congratulations! INSEC is now installed and ready to protect your organization.</content>
<parameter name="filePath">/workspaces/insec/docs/installation.md
