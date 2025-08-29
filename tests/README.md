# INSEC Test Suite

This directory contains comprehensive test suites for all INSEC components.

## 📁 Test Structure

```
tests/
├── unit/                    # Unit tests
│   ├── agent/              # Rust agent unit tests
│   ├── server/             # Go server unit tests
│   └── ui/                 # JavaScript UI unit tests
├── integration/            # Integration tests
│   ├── api/                # API integration tests
│   ├── database/           # Database integration tests
│   └── agent-server/       # Agent-server integration tests
├── e2e/                    # End-to-end tests
│   ├── scenarios/          # E2E test scenarios
│   └── fixtures/           # Test data fixtures
├── performance/            # Performance tests
│   ├── load/               # Load testing
│   └── stress/             # Stress testing
└── security/               # Security tests
    ├── penetration/        # Penetration testing
    └── vulnerability/      # Vulnerability scanning
```

## 🚀 Running Tests

### All Tests
```bash
# Run complete test suite
./scripts/test-all.sh

# Run with coverage
./scripts/test-coverage.sh

# Run in CI mode
./scripts/test-ci.sh
```

### Component-Specific Tests
```bash
# Agent tests (Rust)
cd tests/unit/agent && cargo test

# Server tests (Go)
cd tests/unit/server && go test ./...

# UI tests (JavaScript)
cd tests/unit/ui && npm test

# Integration tests
cd tests/integration && ./run-integration-tests.sh

# E2E tests
cd tests/e2e && ./run-e2e-tests.sh
```

## 🧪 Test Categories

### Unit Tests
- Test individual functions and methods
- Mock external dependencies
- Fast execution (< 1 second per test)
- High coverage target (80%+)

### Integration Tests
- Test component interactions
- Use real dependencies where possible
- Validate data flow between components
- Medium execution time (seconds to minutes)

### End-to-End Tests
- Test complete user workflows
- Use production-like environment
- Validate system behavior from user perspective
- Longer execution time (minutes)

### Performance Tests
- Load testing with multiple concurrent users
- Stress testing with high data volumes
- Benchmarking for performance regression detection
- Resource usage monitoring

### Security Tests
- Penetration testing scenarios
- Vulnerability scanning
- Security control validation
- Compliance testing

## 📋 Test Naming Conventions

### File Naming
```
# Unit tests
user_service_test.go
telemetry_collector_test.rs
UserService.test.js

# Integration tests
api_integration_test.go
database_integration_test.go

# E2E tests
user_registration_e2e_test.go
event_ingestion_e2e_test.go
```

### Test Function Naming
```go
// Go
func TestUserService_GetUser(t *testing.T)
func TestUserService_GetUser_WhenUserNotFound(t *testing.T)
func TestUserService_GetUser_WithInvalidID(t *testing.T)
```

```rust
// Rust
#[test]
fn test_telemetry_collection_success()
#[test]
fn test_telemetry_collection_with_network_error()
#[test]
fn test_telemetry_collection_with_invalid_config()
```

```javascript
// JavaScript
test('renders user data correctly')
test('handles loading state')
test('displays error message on failure')
```

## 🛠️ Test Utilities

### Test Helpers
```go
// test_helpers.go
package test

import (
    "database/sql"
    "github.com/stretchr/testify/suite"
)

type TestSuite struct {
    suite.Suite
    db *sql.DB
}

func (s *TestSuite) SetupTest() {
    // Setup test database
    s.db = setupTestDB()
}

func (s *TestSuite) TearDownTest() {
    // Clean up test data
    cleanupTestDB(s.db)
}
```

### Mock Services
```go
// mock_user_service.go
type MockUserService struct {
    mock.Mock
}

func (m *MockUserService) GetUser(id string) (*User, error) {
    args := m.Called(id)
    return args.Get(0).(*User), args.Error(1)
}
```

### Test Data Factories
```go
// factories.go
func NewTestUser() *User {
    return &User{
        ID:        "test-user-123",
        Email:     "test@example.com",
        Name:      "Test User",
        CreatedAt: time.Now(),
    }
}

func NewTestEvent() *Event {
    return &Event{
        ID:        "test-event-123",
        Type:      "process",
        Timestamp: time.Now(),
        Data: map[string]interface{}{
            "process_name": "test.exe",
            "pid": 1234,
        },
    }
}
```

## 📊 Test Coverage

### Coverage Targets
- **Unit Tests**: 80%+ coverage
- **Integration Tests**: 70%+ coverage
- **Critical Paths**: 90%+ coverage

### Coverage Reporting
```bash
# Generate coverage reports
go test -coverprofile=coverage.out ./...
go tool cover -html=coverage.out -o coverage.html

# Upload to coverage service
./scripts/upload-coverage.sh
```

## 🔄 Test Automation

### CI/CD Integration
```yaml
# .github/workflows/test.yml
name: Tests
on: [push, pull_request]

jobs:
  test:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v3
      - name: Setup Go
        uses: actions/setup-go@v4
        with:
          go-version: '1.21'
      - name: Run tests
        run: make test
      - name: Upload coverage
        uses: codecov/codecov-action@v3
```

### Pre-commit Hooks
```bash
#!/bin/bash
# .git/hooks/pre-commit

# Run tests before commit
go test ./...
if [ $? -ne 0 ]; then
    echo "Tests failed. Please fix before committing."
    exit 1
fi

# Run linting
golint ./...
if [ $? -ne 0 ]; then
    echo "Linting failed. Please fix before committing."
    exit 1
fi
```

## 🎯 Test Scenarios

### Authentication Tests
- Valid login credentials
- Invalid login credentials
- MFA authentication
- Session management
- Password reset flow

### Event Processing Tests
- Valid event ingestion
- Invalid event rejection
- Event enrichment
- Event correlation
- Alert generation

### API Tests
- RESTful endpoint testing
- Authentication middleware
- Input validation
- Error handling
- Rate limiting

### UI Tests
- Component rendering
- User interactions
- Form validation
- Error states
- Responsive design

## 📈 Performance Benchmarks

### Benchmark Tests
```go
// benchmark_test.go
func BenchmarkUserService_GetUser(b *testing.B) {
    service := NewUserService()
    userID := "benchmark-user"

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        _, _ = service.GetUser(userID)
    }
}

func BenchmarkEventIngestion(b *testing.B) {
    ingestor := NewEventIngestor()

    b.ResetTimer()
    for i := 0; i < b.N; i++ {
        event := NewTestEvent()
        _ = ingestor.Ingest(event)
    }
}
```

### Load Testing
```bash
# Load test with Vegeta
echo "POST https://api.insec.com/v1/events" | \
vegeta attack -rate=100 -duration=30s | \
vegeta report

# Load test with k6
k6 run load-test.js
```

## 🔒 Security Testing

### Penetration Testing
```bash
# Automated security scanning
nikto -h https://api.insec.com

# SQL injection testing
sqlmap -u "https://api.insec.com/users?id=1" --batch

# XSS testing
xsstrike -u https://console.insec.com
```

### Vulnerability Scanning
```bash
# Container vulnerability scanning
trivy image yashab/insec-server:latest

# Dependency vulnerability scanning
safety check
npm audit

# SAST (Static Application Security Testing)
semgrep --config=auto .
```

## 📊 Test Reporting

### Test Results
```bash
# Generate test reports
go test -v -json ./... | go-test-report

# Generate HTML reports
allure generate allure-results --clean
allure open
```

### Coverage Reports
```bash
# Generate coverage badges
go test -coverprofile=coverage.out ./...
go tool cover -func=coverage.out

# Upload to coverage services
bash <(curl -s https://codecov.io/bash)
```

## 🐛 Debugging Tests

### Test Debugging Tips
```bash
# Run tests with verbose output
go test -v ./...

# Run specific test
go test -run TestUserService_GetUser ./...

# Debug with delve
dlv test ./internal/auth

# Run tests in IDE
# Use GoLand, VS Code, or Vim with debugging extensions
```

### Common Test Issues
1. **Race Conditions**: Use `-race` flag
2. **Flaky Tests**: Implement retry logic
3. **Slow Tests**: Optimize test setup/teardown
4. **Resource Leaks**: Use test cleanup functions

## 🤝 Contributing to Tests

### Adding New Tests
1. Follow naming conventions
2. Include test documentation
3. Add to appropriate test suite
4. Update CI configuration if needed
5. Ensure tests pass in CI

### Test Maintenance
- Keep tests up to date with code changes
- Remove obsolete tests
- Refactor tests for maintainability
- Monitor test execution time

## 📚 Resources

### Testing Frameworks
- **Go**: `testify`, `ginkgo`, `gomega`
- **Rust**: Built-in test framework
- **JavaScript**: Jest, React Testing Library

### Testing Tools
- **Mocking**: `mockery`, `gomock`
- **Load Testing**: `k6`, `vegeta`
- **Security Testing**: `owasp-zap`, `nikto`
- **Coverage**: `go-cover`, `istanbul`

### Learning Resources
- [Go Testing Documentation](https://golang.org/pkg/testing/)
- [Rust Testing Guide](https://doc.rust-lang.org/book/ch11-00-testing.html)
- [Jest Documentation](https://jestjs.io/docs/getting-started)
- [Testing Best Practices](https://martinfowler.com/bliki/TestPyramid.html)

---

**Last updated:** August 29, 2025</content>
<parameter name="filePath">/workspaces/insec/tests/README.md
