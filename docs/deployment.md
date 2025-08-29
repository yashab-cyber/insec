# INSEC Deployment Guide

This guide covers production deployment strategies for INSEC, including Docker, Kubernetes, and traditional server deployments.

## 🏗️ Deployment Strategies

### Quick Start Deployment
For evaluation and development:
```bash
# Clone repository
git clone https://github.com/yashab-cyber/insec.git
cd insec

# Deploy with Docker Compose
docker-compose up -d

# Access console at http://localhost:3000
```

### Production Deployment Options
1. **Docker Containers** - Containerized deployment
2. **Kubernetes** - Orchestrated deployment
3. **Traditional** - Server-based deployment
4. **Hybrid** - Mixed deployment strategy

## 🐳 Docker Deployment

### Docker Compose (Single Node)
```yaml
version: '3.8'

services:
  # PostgreSQL Database
  postgres:
    image: postgres:15-alpine
    environment:
      POSTGRES_DB: insec
      POSTGRES_USER: insec
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
      - ./db/init.sql:/docker-entrypoint-initdb.d/init.sql
    networks:
      - insec-network
    restart: unless-stopped

  # Redis Cache
  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    networks:
      - insec-network
    restart: unless-stopped

  # INSEC Server
  server:
    image: yashab/insec-server:latest
    environment:
      - DATABASE_URL=postgres://insec:${DB_PASSWORD}@postgres:5432/insec
      - REDIS_URL=redis://redis:6379
      - JWT_SECRET=${JWT_SECRET}
    ports:
      - "8080:8080"
    depends_on:
      - postgres
      - redis
    volumes:
      - ./config/server.yaml:/app/config.yaml
      - ./ssl:/app/ssl
    networks:
      - insec-network
    restart: unless-stopped

  # INSEC Console
  console:
    image: yashab/insec-console:latest
    environment:
      - REACT_APP_API_URL=https://api.insec.com
    ports:
      - "3000:3000"
    depends_on:
      - server
    networks:
      - insec-network
    restart: unless-stopped

  # Nginx Reverse Proxy
  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./ssl:/etc/ssl/certs
    depends_on:
      - server
      - console
    networks:
      - insec-network
    restart: unless-stopped

volumes:
  postgres_data:
  redis_data:

networks:
  insec-network:
    driver: bridge
```

### Multi-Node Docker Swarm
```yaml
version: '3.8'

services:
  server:
    image: yashab/insec-server:latest
    deploy:
      replicas: 3
      restart_policy:
        condition: on-failure
      resources:
        limits:
          memory: 1G
        reservations:
          memory: 512M
    environment:
      - DATABASE_URL=postgres://insec:${DB_PASSWORD}@postgres:5432/insec
      - REDIS_URL=redis://redis:6379
    networks:
      - insec-network

  console:
    image: yashab/insec-console:latest
    deploy:
      replicas: 2
    networks:
      - insec-network

  postgres:
    image: postgres:15-alpine
    deploy:
      placement:
        constraints:
          - node.role == manager
    environment:
      POSTGRES_DB: insec
      POSTGRES_USER: insec
      POSTGRES_PASSWORD: ${DB_PASSWORD}
    volumes:
      - postgres_data:/var/lib/postgresql/data
    networks:
      - insec-network

  redis:
    image: redis:7-alpine
    deploy:
      replicas: 2
    networks:
      - insec-network
```

## ☸️ Kubernetes Deployment

### Prerequisites
```bash
# Install kubectl and helm
curl -LO "https://dl.k8s.io/release/$(curl -L -s https://dl.k8s.io/release/stable.txt)/bin/linux/amd64/kubectl"
chmod +x kubectl && sudo mv kubectl /usr/local/bin/

# Install Helm
curl https://get.helm.sh/helm-v3.12.0-linux-amd64.tar.gz -o helm.tar.gz
tar -zxvf helm.tar.gz && sudo mv linux-amd64/helm /usr/local/bin/
```

### Namespace Creation
```yaml
# namespace.yaml
apiVersion: v1
kind: Namespace
metadata:
  name: insec
  labels:
    name: insec
```

### PostgreSQL Deployment
```yaml
# postgres-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: insec-postgres
  namespace: insec
spec:
  replicas: 1
  selector:
    matchLabels:
      app: insec-postgres
  template:
    metadata:
      labels:
        app: insec-postgres
    spec:
      containers:
      - name: postgres
        image: postgres:15-alpine
        env:
        - name: POSTGRES_DB
          value: "insec"
        - name: POSTGRES_USER
          value: "insec"
        - name: POSTGRES_PASSWORD
          valueFrom:
            secretKeyRef:
              name: insec-secrets
              key: db-password
        ports:
        - containerPort: 5432
        volumeMounts:
        - name: postgres-storage
          mountPath: /var/lib/postgresql/data
      volumes:
      - name: postgres-storage
        persistentVolumeClaim:
          claimName: postgres-pvc

---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: postgres-pvc
  namespace: insec
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 50Gi
```

### INSEC Server Deployment
```yaml
# server-deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: insec-server
  namespace: insec
spec:
  replicas: 3
  selector:
    matchLabels:
      app: insec-server
  template:
    metadata:
      labels:
        app: insec-server
    spec:
      containers:
      - name: server
        image: yashab/insec-server:latest
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: insec-secrets
              key: database-url
        - name: REDIS_URL
          value: "redis://insec-redis:6379"
        - name: JWT_SECRET
          valueFrom:
            secretKeyRef:
              name: insec-secrets
              key: jwt-secret
        ports:
        - containerPort: 8080
        resources:
          limits:
            memory: "1Gi"
            cpu: "500m"
          requests:
            memory: "512Mi"
            cpu: "250m"
        livenessProbe:
          httpGet:
            path: /health
            port: 8080
          initialDelaySeconds: 30
          periodSeconds: 10
        readinessProbe:
          httpGet:
            path: /health/ready
            port: 8080
          initialDelaySeconds: 5
          periodSeconds: 5
```

### Service Definitions
```yaml
# server-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: insec-server
  namespace: insec
spec:
  selector:
    app: insec-server
  ports:
  - port: 8080
    targetPort: 8080
  type: ClusterIP

---
# console-service.yaml
apiVersion: v1
kind: Service
metadata:
  name: insec-console
  namespace: insec
spec:
  selector:
    app: insec-console
  ports:
  - port: 3000
    targetPort: 3000
  type: ClusterIP
```

### Ingress Configuration
```yaml
# ingress.yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: insec-ingress
  namespace: insec
  annotations:
    nginx.ingress.kubernetes.io/ssl-redirect: "true"
    cert-manager.io/cluster-issuer: "letsencrypt-prod"
spec:
  tls:
  - hosts:
    - api.insec.com
    - console.insec.com
    secretName: insec-tls
  rules:
  - host: api.insec.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: insec-server
            port:
              number: 8080
  - host: console.insec.com
    http:
      paths:
      - path: /
        pathType: Prefix
        backend:
          service:
            name: insec-console
            port:
              number: 3000
```

### ConfigMaps and Secrets
```yaml
# configmap.yaml
apiVersion: v1
kind: ConfigMap
metadata:
  name: insec-config
  namespace: insec
data:
  server.yaml: |
    server:
      host: "0.0.0.0"
      port: 8080
    database:
      host: "insec-postgres"
      database: "insec"
    redis:
      host: "insec-redis"

---
# secrets.yaml
apiVersion: v1
kind: Secret
metadata:
  name: insec-secrets
  namespace: insec
type: Opaque
data:
  db-password: <base64-encoded-password>
  jwt-secret: <base64-encoded-jwt-secret>
  database-url: <base64-encoded-database-url>
```

### Horizontal Pod Autoscaler
```yaml
# hpa.yaml
apiVersion: autoscaling/v2
kind: HorizontalPodAutoscaler
metadata:
  name: insec-server-hpa
  namespace: insec
spec:
  scaleTargetRef:
    apiVersion: apps/v1
    kind: Deployment
    name: insec-server
  minReplicas: 3
  maxReplicas: 10
  metrics:
  - type: Resource
    resource:
      name: cpu
      target:
        type: Utilization
        averageUtilization: 70
  - type: Resource
    resource:
      name: memory
      target:
        type: Utilization
        averageUtilization: 80
```

## 🖥️ Traditional Server Deployment

### System Requirements
- **OS**: Ubuntu 20.04+ / CentOS 8+ / RHEL 8+
- **CPU**: 4+ cores
- **RAM**: 8GB+ minimum, 16GB+ recommended
- **Storage**: 100GB+ SSD
- **Network**: 1Gbps connection

### Server Installation
```bash
# Update system
sudo apt update && sudo apt upgrade -y

# Install dependencies
sudo apt install -y postgresql postgresql-contrib redis-server nginx

# Install Go
wget https://go.dev/dl/go1.21.0.linux-amd64.tar.gz
sudo tar -C /usr/local -xzf go1.21.0.linux-amd64.tar.gz
export PATH=$PATH:/usr/local/go/bin

# Install Node.js
curl -fsSL https://deb.nodesource.com/setup_18.x | sudo -E bash -
sudo apt-get install -y nodejs

# Install Rust
curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
source ~/.cargo/env
```

### Database Setup
```bash
# Initialize PostgreSQL
sudo systemctl enable postgresql
sudo systemctl start postgresql

# Create database and user
sudo -u postgres psql
```

```sql
CREATE DATABASE insec;
CREATE USER insec WITH ENCRYPTED PASSWORD 'secure_password';
GRANT ALL PRIVILEGES ON DATABASE insec TO insec;
ALTER USER insec CREATEDB;
\q
```

### Application Deployment
```bash
# Create application user
sudo useradd -r -s /bin/false insec

# Create directories
sudo mkdir -p /opt/insec/{server,agent,console}
sudo mkdir -p /etc/insec
sudo mkdir -p /var/log/insec
sudo mkdir -p /var/lib/insec

# Set permissions
sudo chown -R insec:insec /opt/insec
sudo chown -R insec:insec /var/log/insec
sudo chown -R insec:insec /var/lib/insec
```

### Build and Install
```bash
# Build server
cd /opt/insec/server
git clone https://github.com/yashab-cyber/insec.git .
cd server
go mod download
go build -o insec-server .

# Build console
cd ../ui
npm install
npm run build

# Build agent
cd ../agent
cargo build --release
```

### Systemd Services
```ini
# /etc/systemd/system/insec-server.service
[Unit]
Description=INSEC Server
After=network.target postgresql.service redis.service

[Service]
Type=simple
User=insec
Group=insec
WorkingDirectory=/opt/insec/server
ExecStart=/opt/insec/server/insec-server -config /etc/insec/server.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

```ini
# /etc/systemd/system/insec-agent.service
[Unit]
Description=INSEC Agent
After=network.target

[Service]
Type=simple
User=root
WorkingDirectory=/opt/insec/agent
ExecStart=/opt/insec/agent/target/release/insec-agent -config /etc/insec/agent.yaml
Restart=always
RestartSec=5

[Install]
WantedBy=multi-user.target
```

### Nginx Configuration
```nginx
# /etc/nginx/sites-available/insec
server {
    listen 80;
    server_name api.insec.com console.insec.com;

    # Redirect to HTTPS
    return 301 https://$server_name$request_uri;
}

server {
    listen 443 ssl http2;
    server_name api.insec.com;

    ssl_certificate /etc/ssl/certs/insec.crt;
    ssl_certificate_key /etc/ssl/private/insec.key;

    location / {
        proxy_pass http://localhost:8080;
        proxy_set_header Host $host;
        proxy_set_header X-Real-IP $remote_addr;
        proxy_set_header X-Forwarded-For $proxy_add_x_forwarded_for;
        proxy_set_header X-Forwarded-Proto $scheme;
    }
}

server {
    listen 443 ssl http2;
    server_name console.insec.com;

    ssl_certificate /etc/ssl/certs/insec.crt;
    ssl_certificate_key /etc/ssl/private/insec.key;

    location / {
        root /opt/insec/console/build;
        try_files $uri $uri/ /index.html;
    }
}
```

## ☁️ Cloud Deployments

### AWS Deployment
```yaml
# CloudFormation template excerpt
Resources:
  INSECServer:
    Type: AWS::EC2::Instance
    Properties:
      ImageId: ami-0abcdef1234567890
      InstanceType: t3.medium
      SecurityGroups:
        - !Ref INSECSecurityGroup
      UserData:
        Fn::Base64: |
          #!/bin/bash
          yum update -y
          # Install INSEC...

  INSECDatabase:
    Type: AWS::RDS::DBInstance
    Properties:
      DBInstanceClass: db.t3.micro
      Engine: postgres
      MasterUsername: insec
      MasterUserPassword: !Ref DBPassword
```

### Azure Deployment
```json
{
  "type": "Microsoft.Compute/virtualMachines",
  "apiVersion": "2022-08-01",
  "name": "insec-server",
  "location": "[resourceGroup().location]",
  "properties": {
    "hardwareProfile": {
      "vmSize": "Standard_B2s"
    },
    "storageProfile": {
      "imageReference": {
        "publisher": "Canonical",
        "offer": "UbuntuServer",
        "sku": "18.04-LTS",
        "version": "latest"
      }
    }
  }
}
```

### GCP Deployment
```yaml
# deployment.yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: insec-server
spec:
  replicas: 3
  selector:
    matchLabels:
      app: insec-server
  template:
    metadata:
      labels:
        app: insec-server
    spec:
      containers:
      - name: server
        image: gcr.io/project-id/insec-server:latest
        ports:
        - containerPort: 8080
        env:
        - name: DATABASE_URL
          valueFrom:
            secretKeyRef:
              name: insec-secrets
              key: database-url
```

## 🔧 Post-Deployment Configuration

### SSL Certificate Setup
```bash
# Install Certbot
sudo apt install certbot python3-certbot-nginx

# Get certificates
sudo certbot --nginx -d api.insec.com -d console.insec.com

# Set up auto-renewal
sudo crontab -e
# Add: 0 12 * * * /usr/bin/certbot renew --quiet
```

### Monitoring Setup
```bash
# Install Prometheus and Grafana
sudo apt install prometheus grafana

# Configure INSEC metrics
curl -X POST http://localhost:9090/-/reload

# Access Grafana at http://localhost:3001
```

### Backup Configuration
```bash
# Database backup script
#!/bin/bash
DATE=$(date +%Y%m%d_%H%M%S)
pg_dump -U insec -h localhost insec > /backup/insec_$DATE.sql

# Schedule daily backups
crontab -e
# Add: 0 2 * * * /path/to/backup-script.sh
```

## 📊 Scaling Strategies

### Vertical Scaling
```bash
# Increase server resources
sudo systemctl stop insec-server
# Update instance type or add more CPU/RAM
sudo systemctl start insec-server
```

### Horizontal Scaling
```bash
# Add more server instances
sudo systemctl enable insec-server@2
sudo systemctl start insec-server@2

# Update load balancer configuration
```

### Database Scaling
```bash
# Enable read replicas
# Configure connection pooling
# Implement sharding if needed
```

## 🔍 Health Checks and Monitoring

### Application Health Checks
```bash
# Server health
curl -f http://localhost:8080/health

# Database connectivity
psql -U insec -d insec -c "SELECT 1"

# Redis connectivity
redis-cli ping
```

### Monitoring Dashboards
- **System Metrics**: CPU, memory, disk, network
- **Application Metrics**: Request rates, error rates, latency
- **Business Metrics**: Events processed, alerts generated
- **Security Metrics**: Failed authentications, suspicious activities

## 🚨 Backup and Recovery

### Backup Strategy
```bash
# Database backup
pg_dump -U insec -h localhost insec > backup.sql

# Configuration backup
tar -czf config_backup.tar.gz /etc/insec/

# Application data backup
tar -czf data_backup.tar.gz /var/lib/insec/
```

### Recovery Procedures
```bash
# Database recovery
psql -U insec -d insec < backup.sql

# Configuration recovery
tar -xzf config_backup.tar.gz -C /

# Application recovery
tar -xzf data_backup.tar.gz -C /
```

## 📞 Support and Troubleshooting

### Common Issues
1. **Service won't start**: Check logs in `/var/log/insec/`
2. **Database connection failed**: Verify credentials and network
3. **Agent can't connect**: Check firewall and SSL certificates
4. **High resource usage**: Monitor with `top` and `htop`

### Support Resources
- **Documentation**: [docs.insec.com](https://docs.insec.com)
- **Community**: [GitHub Discussions](https://github.com/yashab-cyber/insec/discussions)
- **Issues**: [GitHub Issues](https://github.com/yashab-cyber/insec/issues)
- **Email**: deployment-support@insec.com

---

**Last updated:** August 29, 2025</content>
<parameter name="filePath">/workspaces/insec/docs/deployment.md
