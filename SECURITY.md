# Security Policy

## Supported Versions

We take security seriously and actively maintain security updates for the following versions:

| Version | Supported          |
| ------- | ------------------ |
| 1.0.x   | :white_check_mark: |
| < 1.0   | :x:                |

## Reporting a Vulnerability

If you discover a security vulnerability in INSEC, please help us by reporting it responsibly.

### How to Report

**Please do NOT report security vulnerabilities through public GitHub issues.**

Instead, please report security vulnerabilities by emailing:
- **security@insec.dev**

You should receive a response within 24 hours. If you don't, please follow up to ensure we received your report.

### What to Include

Please include the following information in your report:

- A clear description of the vulnerability
- Steps to reproduce the issue
- Potential impact of the vulnerability
- Any suggested fixes or mitigations

### Our Process

1. **Acknowledgment**: We'll acknowledge receipt of your report within 24 hours
2. **Investigation**: We'll investigate the issue and work on a fix
3. **Updates**: We'll provide regular updates on our progress
4. **Disclosure**: Once fixed, we'll coordinate disclosure with you
5. **Credit**: We'll credit you in our security advisory (if you wish)

### Security Updates

Security updates will be released as soon as possible after a fix is developed and tested. We'll announce security updates through:

- GitHub Security Advisories
- Release notes
- Our security mailing list (if you wish to subscribe)

## Security Best Practices

### For Contributors

- Never commit sensitive information (API keys, passwords, etc.)
- Use secure coding practices
- Run security scans on your code
- Keep dependencies updated
- Follow the principle of least privilege

### For Users

- Keep INSEC updated to the latest version
- Use strong authentication
- Configure network security properly
- Monitor logs for suspicious activity
- Follow security hardening guidelines

## Known Security Considerations

### Agent Deployment

- The agent requires elevated privileges for system monitoring
- Ensure proper access controls are in place
- Use secure communication channels (TLS 1.3+)
- Implement proper certificate validation

### Data Handling

- All telemetry data is encrypted in transit
- Sensitive data should be encrypted at rest
- Implement proper data retention policies
- Use secure deletion methods

### Network Security

- Use firewalls to restrict agent-server communication
- Implement network segmentation
- Use VPNs for remote deployments
- Monitor for unauthorized access attempts

## Contact

For security-related questions or concerns:
- Email: security@insec.dev
- PGP Key: Available upon request

## Recognition

We appreciate security researchers who help keep INSEC safe. With your permission, we'll acknowledge your contribution in our Hall of Fame.
