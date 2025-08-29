# Contributing to INSEC

Thank you for your interest in contributing to INSEC! This document provides guidelines and information for contributors.

## Code of Conduct

This project follows a code of conduct to ensure a welcoming environment for all contributors. Please read and follow our [Code of Conduct](CODE_OF_CONDUCT.md).

## How to Contribute

### Reporting Issues

- Use the issue templates when creating new issues
- Provide detailed information about the problem
- Include steps to reproduce the issue
- Add relevant logs, screenshots, or other supporting information

### Contributing Code

1. Fork the repository
2. Create a feature branch from `main`
3. Make your changes following our coding standards
4. Add tests for new functionality
5. Ensure all tests pass
6. Update documentation as needed
7. Submit a pull request

### Pull Request Process

1. Update the README.md with details of changes if needed
2. Update the version numbers in any examples files if needed
3. The PR will be merged once you have the sign-off of at least one maintainer

## Development Setup

### Prerequisites

- Rust 1.70+
- Go 1.19+
- Node.js 18+
- Docker (optional)

### Building from Source

```bash
# Clone the repository
git clone https://github.com/yashab-cyber/insec.git
cd insec

# Build all components
./scripts/build.sh

# Run tests
cargo test  # Rust tests
go test ./...  # Go tests
npm test  # UI tests
```

### Running the System

```bash
# Start the server
cd server && ./insec-server

# Start the agent (in another terminal)
cd agent/insec-agent && cargo run

# Start the UI (in another terminal)
cd ui && npm start
```

## Coding Standards

### Rust (Agent)

- Follow the official Rust style guidelines
- Use `rustfmt` for code formatting
- Use `clippy` for linting
- Write comprehensive unit tests
- Document public APIs with rustdoc

### Go (Server)

- Follow standard Go formatting with `gofmt`
- Use `golint` for linting
- Follow Go naming conventions
- Write table-driven tests
- Document public functions

### TypeScript/React (UI)

- Use TypeScript for all new code
- Follow the Airbnb JavaScript style guide
- Use Prettier for code formatting
- Write unit tests with Jest
- Use ESLint for linting

## Testing

- Write unit tests for all new functionality
- Integration tests should cover end-to-end scenarios
- Performance tests for critical components
- Security tests for authentication and authorization

## Documentation

- Update README.md for any new features
- Document API changes in the server documentation
- Update configuration examples
- Add code comments for complex logic

## Security Considerations

- Never commit sensitive information
- Use secure coding practices
- Report security vulnerabilities through our security policy
- Encrypt sensitive data at rest and in transit

## Commit Messages

Use clear, descriptive commit messages:

```
feat: add network telemetry collection
fix: resolve memory leak in agent
docs: update installation instructions
```

## License

By contributing to INSEC, you agree that your contributions will be licensed under the same license as the project (Apache 2.0).
