package auth

import (
	"errors"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/suite"
)

// MockUserRepository is a mock implementation of UserRepository
type MockUserRepository struct {
	mock.Mock
}

func (m *MockUserRepository) GetUserByID(id string) (*User, error) {
	args := m.Called(id)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) GetUserByEmail(email string) (*User, error) {
	args := m.Called(email)
	if args.Get(0) == nil {
		return nil, args.Error(1)
	}
	return args.Get(0).(*User), args.Error(1)
}

func (m *MockUserRepository) CreateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) UpdateUser(user *User) error {
	args := m.Called(user)
	return args.Error(0)
}

func (m *MockUserRepository) DeleteUser(id string) error {
	args := m.Called(id)
	return args.Error(0)
}

// TestSuite for AuthService
type AuthServiceTestSuite struct {
	suite.Suite
	service          *AuthService
	mockRepo         *MockUserRepository
	testUser         *User
	validPassword    string
	invalidPassword  string
}

func (suite *AuthServiceTestSuite) SetupTest() {
	suite.mockRepo = new(MockUserRepository)
	suite.service = NewAuthService(suite.mockRepo, "test-secret-key")

	suite.testUser = &User{
		ID:        "test-user-123",
		Email:     "test@example.com",
		Name:      "Test User",
		Password:  "$2a$10$hashedpassword", // bcrypt hash
		Role:      "analyst",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	suite.validPassword = "correctpassword"
	suite.invalidPassword = "wrongpassword"
}

func (suite *AuthServiceTestSuite) TearDownTest() {
	suite.mockRepo.AssertExpectations(suite.T())
}

// TestAuthenticateUser_Success tests successful user authentication
func (suite *AuthServiceTestSuite) TestAuthenticateUser_Success() {
	// Arrange
	suite.mockRepo.On("GetUserByEmail", suite.testUser.Email).Return(suite.testUser, nil)

	// Act
	token, err := suite.service.AuthenticateUser(suite.testUser.Email, suite.validPassword)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), token)
}

// TestAuthenticateUser_InvalidPassword tests authentication with invalid password
func (suite *AuthServiceTestSuite) TestAuthenticateUser_InvalidPassword() {
	// Arrange
	suite.mockRepo.On("GetUserByEmail", suite.testUser.Email).Return(suite.testUser, nil)

	// Act
	token, err := suite.service.AuthenticateUser(suite.testUser.Email, suite.invalidPassword)

	// Assert
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Contains(suite.T(), err.Error(), "invalid credentials")
}

// TestAuthenticateUser_UserNotFound tests authentication when user doesn't exist
func (suite *AuthServiceTestSuite) TestAuthenticateUser_UserNotFound() {
	// Arrange
	suite.mockRepo.On("GetUserByEmail", "nonexistent@example.com").Return(nil, errors.New("user not found"))

	// Act
	token, err := suite.service.AuthenticateUser("nonexistent@example.com", suite.validPassword)

	// Assert
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), token)
	assert.Contains(suite.T(), err.Error(), "user not found")
}

// TestValidateToken_Success tests successful token validation
func (suite *AuthServiceTestSuite) TestValidateToken_Success() {
	// Arrange
	suite.mockRepo.On("GetUserByEmail", suite.testUser.Email).Return(suite.testUser, nil)
	token, _ := suite.service.AuthenticateUser(suite.testUser.Email, suite.validPassword)

	// Act
	claims, err := suite.service.ValidateToken(token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), claims)
	assert.Equal(suite.T(), suite.testUser.ID, claims.UserID)
	assert.Equal(suite.T(), suite.testUser.Email, claims.Email)
}

// TestValidateToken_InvalidToken tests validation of invalid token
func (suite *AuthServiceTestSuite) TestValidateToken_InvalidToken() {
	// Act
	claims, err := suite.service.ValidateToken("invalid.jwt.token")

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
}

// TestValidateToken_ExpiredToken tests validation of expired token
func (suite *AuthServiceTestSuite) TestValidateToken_ExpiredToken() {
	// Arrange - Create service with very short expiration for testing
	shortLivedService := NewAuthServiceWithExpiration(suite.mockRepo, "test-secret", time.Millisecond*1)
	suite.mockRepo.On("GetUserByEmail", suite.testUser.Email).Return(suite.testUser, nil)
	token, _ := shortLivedService.AuthenticateUser(suite.testUser.Email, suite.validPassword)

	// Wait for token to expire
	time.Sleep(time.Millisecond * 2)

	// Act
	claims, err := shortLivedService.ValidateToken(token)

	// Assert
	assert.Error(suite.T(), err)
	assert.Nil(suite.T(), claims)
	assert.Contains(suite.T(), err.Error(), "token is expired")
}

// TestRefreshToken_Success tests successful token refresh
func (suite *AuthServiceTestSuite) TestRefreshToken_Success() {
	// Arrange
	suite.mockRepo.On("GetUserByEmail", suite.testUser.Email).Return(suite.testUser, nil)
	token, _ := suite.service.AuthenticateUser(suite.testUser.Email, suite.validPassword)

	// Act
	newToken, err := suite.service.RefreshToken(token)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), newToken)
	assert.NotEqual(suite.T(), token, newToken) // Should be a different token
}

// TestRefreshToken_InvalidToken tests refresh with invalid token
func (suite *AuthServiceTestSuite) TestRefreshToken_InvalidToken() {
	// Act
	newToken, err := suite.service.RefreshToken("invalid.jwt.token")

	// Assert
	assert.Error(suite.T(), err)
	assert.Empty(suite.T(), newToken)
}

// TestHashPassword tests password hashing
func (suite *AuthServiceTestSuite) TestHashPassword() {
	// Act
	hashedPassword, err := suite.service.HashPassword(suite.validPassword)

	// Assert
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), hashedPassword)
	assert.NotEqual(suite.T(), suite.validPassword, hashedPassword)
	assert.True(suite.T(), len(hashedPassword) > len(suite.validPassword))
}

// TestVerifyPassword_Success tests successful password verification
func (suite *AuthServiceTestSuite) TestVerifyPassword_Success() {
	// Arrange
	hashedPassword, _ := suite.service.HashPassword(suite.validPassword)

	// Act
	isValid := suite.service.VerifyPassword(suite.validPassword, hashedPassword)

	// Assert
	assert.True(suite.T(), isValid)
}

// TestVerifyPassword_InvalidPassword tests password verification with wrong password
func (suite *AuthServiceTestSuite) TestVerifyPassword_InvalidPassword() {
	// Arrange
	hashedPassword, _ := suite.service.HashPassword(suite.validPassword)

	// Act
	isValid := suite.service.VerifyPassword(suite.invalidPassword, hashedPassword)

	// Assert
	assert.False(suite.T(), isValid)
}

// TestGenerateSecureToken tests secure token generation
func (suite *AuthServiceTestSuite) TestGenerateSecureToken() {
	// Act
	token1 := suite.service.GenerateSecureToken()
	token2 := suite.service.GenerateSecureToken()

	// Assert
	assert.NotEmpty(suite.T(), token1)
	assert.NotEmpty(suite.T(), token2)
	assert.NotEqual(suite.T(), token1, token2) // Should be unique
	assert.True(suite.T(), len(token1) >= 32)  // Should be sufficiently long
}

// TestConcurrentAuthentication tests concurrent authentication requests
func (suite *AuthServiceTestSuite) TestConcurrentAuthentication() {
	// Arrange
	suite.mockRepo.On("GetUserByEmail", suite.testUser.Email).Return(suite.testUser, nil).Maybe()

	// Act - Run multiple authentication requests concurrently
	done := make(chan bool, 10)
	for i := 0; i < 10; i++ {
		go func() {
			_, err := suite.service.AuthenticateUser(suite.testUser.Email, suite.validPassword)
			assert.NoError(suite.T(), err)
			done <- true
		}()
	}

	// Assert - Wait for all goroutines to complete
	for i := 0; i < 10; i++ {
		<-done
	}
}

// Run the test suite
func TestAuthServiceTestSuite(t *testing.T) {
	suite.Run(t, new(AuthServiceTestSuite))
}

// Benchmark tests
func BenchmarkAuthenticateUser(b *testing.B) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "benchmark-secret")

	user := &User{
		ID:       "bench-user",
		Email:    "bench@example.com",
		Password: "$2a$10$hashedpassword",
	}

	mockRepo.On("GetUserByEmail", user.Email).Return(user, nil)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.AuthenticateUser(user.Email, "password")
	}
}

func BenchmarkValidateToken(b *testing.B) {
	mockRepo := new(MockUserRepository)
	service := NewAuthService(mockRepo, "benchmark-secret")

	user := &User{
		ID:    "bench-user",
		Email: "bench@example.com",
	}

	mockRepo.On("GetUserByEmail", user.Email).Return(user, nil)
	token, _ := service.AuthenticateUser(user.Email, "password")

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		_, _ = service.ValidateToken(token)
	}
}</content>
<parameter name="filePath">/workspaces/insec/tests/unit/server/auth_service_test.go
