# # INSEC: Enterprise Insider-Threat Protection

![License](https://img.shields.io/badge/License-Apache%202.0-blue.svg)
![Build Status](https://img.shields.io/github/actions/workflow/status/yashab-cyber/insec/ci.yml)
![Version](https://img.shields.io/github/v/release/yashab-cyber/insec)
![Stars](https://img.shields.io/github/stars/yashab-cyber/insec)
![Forks](https://img.shields.io/github/forks/yashab-cyber/insec)
![Issues](https://img.shields.io/github/issues/yashab-cyber/insec)
![Contributors](https://img.shields.io/github/contributors/yashab-cyber/insec)
![Last Commit](https://img.shields.io/github/last-commit/yashab-cyber/insec)

**Tagline:** "Stop data walking out the door."

INSEC is a privacy-respectful, enterprise-grade insider-threat detection and response platform. It reduces data exfiltration, account misuse, policy violations, and sabotage by combining endpoint telemetry, UEBA (user & entity behavior analytics), policy controls, and automated response.

## 🚀 Quick Start

```bash
# Clone the repository
git clone https://github.com/yashab-cyber/insec.git
cd insec

# Build all components
./scripts/build.sh

# Start the server
cd server && go run main.go

# Start the UI (in another terminal)
cd ui && npm start

# Build and run the agent
cd agent && cargo build --release
./target/release/insec-agent
```

## 🏗️ Architecture

### Endpoint Agent (INSEC Agent)
- Cross-platform (Windows/macOS/Linux) using Rust.
- Collects telemetry, enforces policies, runs local detections, performs containment actions.
- Low resource usage: <2% CPU p95 idle, <200MB RAM.
- Auto-update, offline cache/queue, self-protection, signed binaries.

### Control Plane (INSEC Cloud/Server)
- Services: AuthN/Z (SAML/OIDC, SCIM), Policy Engine, Analytics/UEBA, Alerting, Orchestrator, API Gateway, Event Ingest, Storage.
- Multi-tenant, horizontally scalable, stateless services with message bus (NATS/Kafka).
- Encrypt data in transit (mTLS) and at rest (AES-256, envelope keys; per-tenant keys).

### Data Plane
- Hot path: event ingest → stream processing → rules engine → UEBA scores → alerting.
- Warm path: data lake for historical search, reporting, model training.

### Admin UI (INSEC Console)
- Web app (React/TypeScript) with RBAC: Org Admin, SecOps Analyst, Auditor, Read-Only.
- Dashboards, investigations, policy editor, search & analytics.

## 📁 Project Structure
- `agent/`: Rust project for endpoint agent.
- `server/`: Go project for control plane services.
- `ui/`: React TypeScript app for console.
- `docs/`: Documentation.
- `scripts/`: Build and deployment scripts.
- `tests/`: Test suites.

## 🛠️ Getting Started
1. Install dependencies: Rust, Go, Node.js.
2. For agent: `cd agent && cargo build`.
3. For server: `cd server && go build`.
4. For UI: `cd ui && npm start`.

## 🎯 Core Use Cases
- Data Exfiltration detection.
- Privilege Misuse.
- Account Compromise.
- Policy Violations.
- Lateral Movement & Recon.
- Insider Fraud/Sabotage.

## 🔒 Compliance & Privacy
- Per-policy masking/redaction.
- No keystrokes/no content by default.
- Configurable data retention.
- Region pinning & tenant KMS integration.

## 🔍 Detection & Analytics
- Rules Engine with deterministic rules.
- UEBA with baseline modeling.
- Correlation for narratives.
- False-positive controls.

## ⚡ Response & Orchestration
- Automations/Playbooks.
- Approval gates for high-impact actions.
- Forensics with artifact capture.

## 🔗 Integrations
- Identity & Device: Okta/Azure AD/Google.
- SIEM/SOAR: Splunk, Elastic, Sentinel.
- Ticketing/ChatOps: Jira/ServiceNow, Slack/Teams.
- Dev/Cloud: GitHub/GitLab, AWS/GCP/Azure.

## 🛡️ Security & Hardening
- Code-signing and notarization.
- mTLS with cert pinning.
- Agent self-protection.
- Supply chain security.

## 📊 Performance & Reliability
- <50ms event enqueue latency on-host.
- <5s end-to-end alerting p95.
- Auto-update with staged rollouts.

## 📦 Packaging & Deployment
- Windows: MSI with signed binaries.
- macOS: Notarized PKG.
- Linux: DEB/RPM + systemd units.

## 👁️ Observability & QA
- Metrics, tracing, structured logs.
- Unit, integration, load tests.
- Golden datasets for regression.

## 🌐 APIs
- Ingest: `/v1/events`.
- Query: `/v1/search`, `/v1/entities`.
- Alerts: `/v1/alerts`.
- Policies: `/v1/policies`.
- Webhooks with OAuth2.

## 🤝 Contributing

We welcome contributions! Please see our [Contributing Guide](CONTRIBUTING.md) for details.

## 📞 Support & Community

- 📧 **Email:** yashabalam707@gmail.com
- 💬 **Discord:** [ZehraSec Community Server](https://discord.gg/zehrasec)
- 📱 **WhatsApp:** [Business Channel](https://whatsapp.com/channel/0029Vaoa1GfKLaHlL0Kc8k1q)

## 💰 Support INSEC Development

Your donations help accelerate the development of advanced insider-threat protection tools. See [DONATE.md](DONATE.md) for donation options and funding goals, or [CRYPTO.md](CRYPTO.md) for cryptocurrency donations.

## 🌐 Connect with Us

**Official Channels:**
- 🌐 **Website:** [www.zehrasec.com](https://www.zehrasec.com)
- 📸 **Instagram:** [@_zehrasec](https://www.instagram.com/_zehrasec?igsh=bXM0cWl1ejdoNHM4)
- 📘 **Facebook:** [ZehraSec Official](https://www.facebook.com/profile.php?id=61575580721849)
- 🐦 **X (Twitter):** [@zehrasec](https://x.com/zehrasec?t=Tp9LOesZw2d2yTZLVo0_GA&s=08)
- 💼 **LinkedIn:** [ZehraSec Company](https://www.linkedin.com/company/zehrasec)

### 👨‍💻 Connect with Yashab Alam
- 💻 **GitHub:** [@yashab-cyber](https://github.com/yashab-cyber)
- 📸 **Instagram:** [@yashab.alam](https://www.instagram.com/yashab.alam)
- 💼 **LinkedIn:** [Yashab Alam](https://www.linkedin.com/in/yashab-alam)

## 📄 License

This project is licensed under the Apache License 2.0 - see the [LICENSE](LICENSE) file for details.

---

**Made with ❤️ by Yashab Alam and the ZehraSec team**

*Repository: [github.com/yashab-cyber/insec](https://github.com/yashab-cyber/insec)*NSEC: Enterprise Insider-Threat Protection

**Tagline:** “Stop data walking out the door.”

INSEC is a privacy-respectful, enterprise-grade insider-threat detection and response platform. It reduces data exfiltration, account misuse, policy violations, and sabotage by combining endpoint telemetry, UEBA (user & entity behavior analytics), policy controls, and automated response.

## Architecture

### Endpoint Agent (INSEC Agent)
- Cross-platform (Windows/macOS/Linux) using Rust.
- Collects telemetry, enforces policies, runs local detections, performs containment actions.
- Low resource usage: <2% CPU p95 idle, <200MB RAM.
- Auto-update, offline cache/queue, self-protection, signed binaries.

### Control Plane (INSEC Cloud/Server)
- Services: AuthN/Z (SAML/OIDC, SCIM), Policy Engine, Analytics/UEBA, Alerting, Orchestrator, API Gateway, Event Ingest, Storage.
- Multi-tenant, horizontally scalable, stateless services with message bus (NATS/Kafka).
- Encrypt data in transit (mTLS) and at rest (AES-256, envelope keys; per-tenant keys).

### Data Plane
- Hot path: event ingest → stream processing → rules engine → UEBA scores → alerting.
- Warm path: data lake for historical search, reporting, model training.

### Admin UI (INSEC Console)
- Web app (React/TypeScript) with RBAC: Org Admin, SecOps Analyst, Auditor, Read-Only.
- Dashboards, investigations, policy editor, search & analytics.

## Project Structure
- `agent/`: Rust project for endpoint agent.
- `server/`: Go project for control plane services.
- `ui/`: React TypeScript app for console.
- `docs/`: Documentation.
- `scripts/`: Build and deployment scripts.
- `tests/`: Test suites.

## Getting Started
1. Install dependencies: Rust, Go, Node.js.
2. For agent: `cd agent && cargo build`.
3. For server: `cd server && go build`.
4. For UI: `cd ui && npm start`.

## Core Use Cases
- Data Exfiltration detection.
- Privilege Misuse.
- Account Compromise.
- Policy Violations.
- Lateral Movement & Recon.
- Insider Fraud/Sabotage.

## Compliance & Privacy
- Per-policy masking/redaction.
- No keystrokes/no content by default.
- Configurable data retention.
- Region pinning & tenant KMS integration.

## Detection & Analytics
- Rules Engine with deterministic rules.
- UEBA with baseline modeling.
- Correlation for narratives.
- False-positive controls.

## Response & Orchestration
- Automations/Playbooks.
- Approval gates for high-impact actions.
- Forensics with artifact capture.

## Integrations
- Identity & Device: Okta/Azure AD/Google.
- SIEM/SOAR: Splunk, Elastic, Sentinel.
- Ticketing/ChatOps: Jira/ServiceNow, Slack/Teams.
- Dev/Cloud: GitHub/GitLab, AWS/GCP/Azure.

## Security & Hardening
- Code-signing and notarization.
- mTLS with cert pinning.
- Agent self-protection.
- Supply chain security.

## Performance & Reliability
- <50ms event enqueue latency on-host.
- <5s end-to-end alerting p95.
- Auto-update with staged rollouts.

## Packaging & Deployment
- Windows: MSI with signed binaries.
- macOS: Notarized PKG.
- Linux: DEB/RPM + systemd units.

## Observability & QA
- Metrics, tracing, structured logs.
- Unit, integration, load tests.
- Golden datasets for regression.

## APIs
- Ingest: `/v1/events`.
- Query: `/v1/search`, `/v1/entities`.
- Alerts: `/v1/alerts`.
- Policies: `/v1/policies`.
- Webhooks with OAuth2.

## Acceptance Criteria for v1
- Agents enroll and stream events with mTLS.
- Policies deploy in <5 minutes.
- Detect and alert on key scenarios.
- Automated responses: host isolation, USB block, ticket creation.
- RBAC in Console; audit log.
- SIEM integration.
- Performance targets met.
- Privacy controls implemented.