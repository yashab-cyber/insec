# Support

## Getting Help

If you need help with INSEC, here are the best ways to get assistance:

## Documentation

- [README.md](../README.md) - Main project documentation
- [CONTRIBUTING.md](../CONTRIBUTING.md) - How to contribute to the project
- [SECURITY.md](../SECURITY.md) - Security-related information

## Community Support

### GitHub Discussions
Use GitHub Discussions for:
- General questions about INSEC
- Feature requests and ideas
- Sharing best practices
- Community announcements

[Join the discussion →](https://github.com/yashab-cyber/insec/discussions)

### GitHub Issues
Use GitHub Issues for:
- Bug reports
- Feature requests
- Security vulnerabilities (see SECURITY.md)

[Report an issue →](https://github.com/yashab-cyber/insec/issues)

## Enterprise Support

For enterprise customers, we offer:
- Priority support
- Custom integrations
- Training and consulting
- SLA guarantees

Contact us at enterprise@insec.dev for more information.

## Troubleshooting

### Common Issues

#### Agent Won't Start
- Ensure you have the necessary permissions
- Check system requirements (Rust 1.70+)
- Verify the configuration file

#### Server Connection Issues
- Check network connectivity
- Verify server is running on the correct port
- Review firewall settings

#### UI Not Loading
- Ensure Node.js 18+ is installed
- Check that the build completed successfully
- Verify port 3000 is available

### Debug Mode

To run components in debug mode:

```bash
# Agent with debug logging
cd agent/insec-agent
RUST_LOG=debug cargo run

# Server with debug logging
cd server
GIN_MODE=debug ./insec-server

# UI with debug mode
cd ui
npm start
```

### Logs

Logs are typically located in:
- Agent: `~/.insec/agent.log`
- Server: `/var/log/insec/server.log`
- UI: Browser developer console

## System Requirements

### Minimum Requirements
- **OS**: Linux, Windows, macOS
- **CPU**: 2 cores
- **RAM**: 4GB
- **Storage**: 1GB free space

### Recommended Requirements
- **OS**: Ubuntu 20.04+, CentOS 8+, Windows Server 2019+
- **CPU**: 4+ cores
- **RAM**: 8GB+
- **Storage**: 10GB+ free space

### Component-Specific Requirements

#### Agent
- Rust 1.70+
- System monitoring permissions
- Network access to server

#### Server
- Go 1.19+
- Database (PostgreSQL/MySQL recommended)
- TLS certificate for production

#### UI
- Node.js 18+
- Modern web browser
- 1920x1080 minimum resolution

## Performance Tuning

### Agent Optimization
- Adjust polling intervals in configuration
- Limit monitored processes if needed
- Use efficient serialization formats

### Server Optimization
- Configure connection pooling
- Set appropriate timeouts
- Use load balancing for high traffic

### UI Optimization
- Enable gzip compression
- Use CDN for static assets
- Implement caching strategies

## Contact Information

- **General Support**: support@insec.dev
- **Enterprise Support**: enterprise@insec.dev
- **Security Issues**: security@insec.dev
- **Business Inquiries**: business@insec.dev

## Response Times

- **Community Support**: Within 48 hours
- **Enterprise Support**: Within 4 hours
- **Security Issues**: Within 24 hours

## Service Level Agreements

Enterprise customers receive:
- 99.9% uptime guarantee
- 1-hour critical issue response
- 24/7 phone support
- Dedicated technical account manager
