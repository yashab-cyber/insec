# Contributing to INSEC

Thank you for your interest in contributing to INSEC! This guide explains how to get started with development, coding standards, and contribution processes.

## 🚀 Getting Started

### Development Environment Setup

#### Prerequisites
- **Go**: 1.19+ for server development
- **Rust**: 1.70+ for agent development
- **Node.js**: 18+ for UI development
- **PostgreSQL**: 13+ for database
- **Redis**: 6+ for caching
- **Git**: Latest version

#### Clone and Setup
```bash
# Clone the repository
git clone https://github.com/yashab-cyber/insec.git
cd insec

# Set up Git hooks
cp scripts/pre-commit .git/hooks/
chmod +x .git/hooks/pre-commit

# Install dependencies
./scripts/setup-dev.sh

# Set up local database
./scripts/setup-db.sh

# Start development environment
./scripts/dev-start.sh
```

### Project Structure
```
insec/
├── agent/                 # Rust agent code
│   ├── src/
│   ├── Cargo.toml
│   └── config/
├── server/                # Go server code
│   ├── cmd/
│   ├── internal/
│   ├── pkg/
│   └── go.mod
├── ui/                    # React UI code
│   ├── src/
│   ├── public/
│   └── package.json
├── docs/                  # Documentation
├── scripts/               # Build and utility scripts
├── tests/                 # Test suites
└── docker/                # Docker configurations
```

## 💻 Development Workflow

### 1. Choose an Issue
- Check [GitHub Issues](https://github.com/yashab-cyber/insec/issues) for open tasks
- Look for issues labeled `good first issue` or `help wanted`
- Comment on the issue to indicate you're working on it

### 2. Create a Branch
```bash
# Create and switch to a feature branch
git checkout -b feature/your-feature-name

# Or for bug fixes
git checkout -b fix/issue-number-description
```

### 3. Make Changes
- Follow the coding standards below
- Write tests for your changes
- Update documentation if needed
- Ensure all tests pass

### 4. Commit Changes
```bash
# Stage your changes
git add .

# Commit with descriptive message
git commit -m "feat: add user authentication endpoint

- Implement JWT-based authentication
- Add user registration and login
- Include password hashing with bcrypt
- Add input validation and error handling

Closes #123"
```

### 5. Push and Create Pull Request
```bash
# Push your branch
git push origin feature/your-feature-name

# Create a Pull Request on GitHub
# Include a clear description and link to the issue
```

## 📝 Coding Standards

### Go (Server) Standards

#### Code Style
```go
// Use gofmt for formatting
gofmt -w .

// Use golint for linting
golint ./...

// Follow these naming conventions
type UserService struct{}          // PascalCase for exported types
func (s *UserService) GetUser() {} // PascalCase for exported functions
userID string                     // camelCase for variables
const MaxRetries = 3              // PascalCase for constants
```

#### Project Structure
```
server/
├── cmd/                    # Main applications
│   └── server/
├── internal/              # Private application code
│   ├── auth/
│   ├── events/
│   └── users/
├── pkg/                   # Public library code
│   ├── api/
│   └── models/
└── go.mod
```

#### Error Handling
```go
// Use custom error types
type ValidationError struct {
    Field   string
    Message string
}

func (e ValidationError) Error() string {
    return fmt.Sprintf("validation error on field %s: %s", e.Field, e.Message)
}

// Wrap errors with context
if err != nil {
    return fmt.Errorf("failed to create user: %w", err)
}
```

### Rust (Agent) Standards

#### Code Style
```rust
// Use rustfmt for formatting
cargo fmt

// Use clippy for linting
cargo clippy

// Follow these naming conventions
struct UserService {}           // PascalCase for types
fn get_user() {}                // snake_case for functions
let user_id: String;            // snake_case for variables
const MAX_RETRIES: u32 = 3;     // SCREAMING_SNAKE_CASE for constants
```

#### Project Structure
```
agent/
├── src/
│   ├── main.rs              # Application entry point
│   ├── lib.rs               # Library code
│   ├── config.rs            # Configuration handling
│   ├── telemetry.rs         # Telemetry collection
│   └── agent.rs             # Core agent logic
├── tests/                   # Integration tests
└── Cargo.toml
```

#### Error Handling
```rust
// Use Result and custom error types
#[derive(Debug, thiserror::Error)]
pub enum AgentError {
    #[error("Configuration error: {0}")]
    Config(String),
    #[error("Network error: {0}")]
    Network(#[from] reqwest::Error),
}

// Use ? operator for error propagation
fn connect_to_server(config: &Config) -> Result<(), AgentError> {
    let client = reqwest::Client::new();
    client.post(&config.server_url)
        .json(&telemetry)
        .send()
        .await?;
    Ok(())
}
```

### JavaScript/React Standards

#### Code Style
```javascript
// Use ESLint and Prettier
npm run lint
npm run format

// Follow these naming conventions
const UserService = () => {}     // PascalCase for components
function getUser() {}           // camelCase for functions
const userId = '123';           // camelCase for variables
const MAX_RETRIES = 3;          // SCREAMING_SNAKE_CASE for constants
```

#### Project Structure
```
ui/
├── src/
│   ├── components/          # Reusable components
│   ├── pages/              # Page components
│   ├── hooks/              # Custom hooks
│   ├── services/           # API services
│   ├── utils/              # Utility functions
│   ├── types/              # TypeScript types
│   └── styles/             # Stylesheets
├── public/                 # Static assets
└── package.json
```

## 🧪 Testing

### Unit Tests

#### Go Unit Tests
```go
// user_service_test.go
package auth

import (
    "testing"
    "github.com/stretchr/testify/assert"
)

func TestUserService_GetUser(t *testing.T) {
    // Arrange
    service := NewUserService()
    userID := "user123"

    // Act
    user, err := service.GetUser(userID)

    // Assert
    assert.NoError(t, err)
    assert.NotNil(t, user)
    assert.Equal(t, userID, user.ID)
}
```

#### Rust Unit Tests
```rust
// lib.rs or specific test file
#[cfg(test)]
mod tests {
    use super::*;

    #[test]
    fn test_telemetry_collection() {
        // Arrange
        let config = Config::default();

        // Act
        let telemetry = collect_telemetry(&config);

        // Assert
        assert!(!telemetry.is_empty());
        assert!(telemetry.iter().all(|t| t.timestamp > 0));
    }
}
```

#### JavaScript Unit Tests
```javascript
// UserService.test.js
import { render, screen } from '@testing-library/react';
import UserService from './UserService';

test('renders user data', async () => {
  // Arrange
  const mockUser = { id: '1', name: 'John Doe' };

  // Act
  render(<UserService user={mockUser} />);

  // Assert
  expect(screen.getByText('John Doe')).toBeInTheDocument();
});
```

### Integration Tests
```bash
# Run all tests
./scripts/test.sh

# Run specific test suite
go test ./internal/auth/...
cargo test --package insec-agent
npm test -- --testPathPattern=auth

# Run with coverage
go test -cover ./...
cargo tarpaulin
npm test -- --coverage
```

### End-to-End Tests
```bash
# Start test environment
./scripts/test-e2e.sh

# Run E2E tests
npm run test:e2e

# Test scenarios
- User registration and login
- Event ingestion and processing
- Alert creation and notification
- Dashboard functionality
```

## 📚 Documentation

### Code Documentation

#### Go Documentation
```go
// Package auth provides authentication services for INSEC.
//
// It handles user authentication, authorization, and session management
// using JWT tokens and role-based access control.
package auth

// UserService handles user-related operations.
//
// It provides methods for user management including creation, retrieval,
// updating, and deletion of user accounts.
type UserService struct {
    db *sql.DB
}

// NewUserService creates a new instance of UserService.
func NewUserService(db *sql.DB) *UserService {
    return &UserService{db: db}
}

// GetUser retrieves a user by their ID.
//
// Returns the user if found, or an error if the user doesn't exist
// or if there's a database error.
func (s *UserService) GetUser(id string) (*User, error) {
    // Implementation...
}
```

#### Rust Documentation
```rust
/// Authentication module for INSEC agent.
///
/// This module handles secure communication with the INSEC server,
/// including TLS certificate validation and token management.
pub mod auth;

/// User authentication service.
///
/// Provides functionality for user login, token refresh, and
/// secure credential storage.
pub struct AuthService {
    config: Config,
    client: reqwest::Client,
}

impl AuthService {
    /// Creates a new authentication service instance.
    ///
    /// # Arguments
    ///
    /// * `config` - The authentication configuration
    ///
    /// # Returns
    ///
    /// A new `AuthService` instance
    pub fn new(config: Config) -> Self {
        Self {
            config,
            client: reqwest::Client::new(),
        }
    }

    /// Authenticates a user with the server.
    ///
    /// # Arguments
    ///
    /// * `username` - The user's username
    /// * `password` - The user's password
    ///
    /// # Returns
    ///
    /// A JWT token on success, or an error on failure
    pub async fn login(&self, username: &str, password: &str) -> Result<String, AuthError> {
        // Implementation...
    }
}
```

### API Documentation
```go
// API endpoint documentation
// @Summary Get user by ID
// @Description Retrieves a user by their unique identifier
// @Tags users
// @Accept json
// @Produce json
// @Param id path string true "User ID"
// @Success 200 {object} User
// @Failure 404 {object} ErrorResponse
// @Router /users/{id} [get]
// @Security BearerAuth
func (h *UserHandler) GetUser(c *gin.Context) {
    // Implementation...
}
```

## 🔒 Security Considerations

### Code Review Checklist
- [ ] No hardcoded secrets or credentials
- [ ] Input validation and sanitization
- [ ] SQL injection prevention
- [ ] XSS protection in web interfaces
- [ ] Proper error handling (no sensitive data leakage)
- [ ] Authentication and authorization checks
- [ ] Secure defaults for configuration

### Security Testing
```bash
# Run security scans
gosec ./...
cargo audit
npm audit

# Check for vulnerabilities
trivy fs .
safety check
```

## 📋 Pull Request Process

### PR Template
```markdown
## Description
Brief description of the changes made.

## Type of Change
- [ ] Bug fix (non-breaking change)
- [ ] New feature (non-breaking change)
- [ ] Breaking change
- [ ] Documentation update

## Testing
- [ ] Unit tests added/updated
- [ ] Integration tests added/updated
- [ ] E2E tests added/updated
- [ ] Manual testing completed

## Checklist
- [ ] Code follows style guidelines
- [ ] Documentation updated
- [ ] Security review completed
- [ ] Tests pass
- [ ] No breaking changes

## Related Issues
Closes #123
```

### Review Process
1. **Automated Checks**: CI/CD runs tests and linting
2. **Code Review**: At least one maintainer reviews the code
3. **Security Review**: Security team reviews for vulnerabilities
4. **Testing**: QA team validates the changes
5. **Approval**: Maintainers approve and merge

## 🎯 Contribution Guidelines

### Issue Reporting
- Use issue templates when available
- Provide clear, reproducible steps
- Include environment details and logs
- Check for existing issues first

### Feature Requests
- Describe the problem you're solving
- Explain your proposed solution
- Consider alternative approaches
- Include mockups or examples if applicable

### Code of Conduct
- Be respectful and inclusive
- Focus on constructive feedback
- Help newcomers learn and contribute
- Report violations to maintainers

## 🏆 Recognition

### Contributor Recognition
- Contributors are listed in CONTRIBUTORS.md
- Top contributors featured in release notes
- Special recognition for significant contributions
- Invitation to contributor meetings

### Rewards Program
- Bug bounty program for security issues
- Swag for first-time contributors
- Feature bounties for major enhancements
- Sponsored conference attendance

## 📞 Getting Help

### Communication Channels
- **GitHub Issues**: For bugs and feature requests
- **GitHub Discussions**: For questions and general discussion
- **Discord**: For real-time chat and community support
- **Email**: maintainers@insec.com for private matters

### Office Hours
- **Weekly Community Call**: Every Thursday at 3 PM UTC
- **Mentorship Program**: Pair with experienced contributors
- **Documentation Office Hours**: Help improve documentation

## 🙏 Thank You

Your contributions help make INSEC better for everyone in the cybersecurity community. Every contribution, from fixing a typo to implementing a major feature, is valuable and appreciated.

**Happy contributing! 🚀**

---

**Last updated:** August 29, 2025</content>
<parameter name="filePath">/workspaces/insec/docs/contributing.md
