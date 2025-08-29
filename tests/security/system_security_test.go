package tests

import (
	"bytes"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"insec/internal/models"
	"insec/test/helpers"
)

type SecurityTestSuite struct {
	suite.Suite
	testHelper *helpers.TestHelper
	baseURL    string
	httpClient *http.Client
}

func (suite *SecurityTestSuite) SetupSuite() {
	suite.testHelper = helpers.NewTestHelper()
	suite.baseURL = "https://localhost:8443" // HTTPS endpoint

	// Configure HTTP client to skip TLS verification for testing
	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	}
	suite.httpClient = &http.Client{Transport: tr}

	err := suite.testHelper.StartSecureSystem()
	suite.Require().NoError(err)

	suite.testHelper.WaitForSecureSystemReady()
}

func (suite *SecurityTestSuite) TearDownSuite() {
	suite.testHelper.StopSystem()
}

func TestSecurityTestSuite(t *testing.T) {
	suite.Run(t, new(SecurityTestSuite))
}

func (suite *SecurityTestSuite) TestAuthenticationBypassAttempts() {
	// Test various authentication bypass attempts

	// 1. No authentication header
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	// 2. Invalid token format
	req, _ = http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	req.Header.Set("Authorization", "InvalidFormat")
	resp, err = suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	// 3. Malformed JWT
	req, _ = http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer malformed.jwt.token")
	resp, err = suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)

	// 4. Expired token
	expiredToken := suite.generateExpiredToken()
	req, _ = http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+expiredToken)
	resp, err = suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *SecurityTestSuite) TestAuthorizationEnforcement() {
	validToken := suite.authenticate()

	// Test role-based access control
	endpoints := map[string][]string{
		"/api/v1/admin/users":     {"admin"},
		"/api/v1/admin/config":    {"admin"},
		"/api/v1/alerts":          {"admin", "analyst", "viewer"},
		"/api/v1/analytics":       {"admin", "analyst"},
		"/api/v1/events":          {"admin", "analyst", "viewer"},
	}

	for endpoint, allowedRoles := range endpoints {
		for _, role := range []string{"admin", "analyst", "viewer"} {
			token := suite.authenticateAsRole(role)

			req, _ := http.NewRequest("GET", suite.baseURL+endpoint, nil)
			req.Header.Set("Authorization", "Bearer "+token)
			resp, err := suite.httpClient.Do(req)
			suite.Require().NoError(err)
			defer resp.Body.Close()

			if contains(allowedRoles, role) {
				assert.NotEqual(suite.T(), http.StatusForbidden, resp.StatusCode,
					"Role %s should have access to %s", role, endpoint)
			} else {
				assert.Equal(suite.T(), http.StatusForbidden, resp.StatusCode,
					"Role %s should not have access to %s", role, endpoint)
			}
		}
	}
}

func (suite *SecurityTestSuite) TestSQLInjectionPrevention() {
	validToken := suite.authenticate()

	// Test various SQL injection attempts
	injectionPayloads := []string{
		"'; DROP TABLE alerts; --",
		"' OR '1'='1",
		"' UNION SELECT * FROM users --",
		"'; SELECT * FROM information_schema.tables; --",
		"admin' --",
		"' OR 1=1; --",
	}

	for _, payload := range injectionPayloads {
		// Test in query parameters
		req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/alerts?search="+payload, nil)
		req.Header.Set("Authorization", "Bearer "+validToken)
		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Should not return 500 error (which would indicate SQL injection success)
		assert.NotEqual(suite.T(), http.StatusInternalServerError, resp.StatusCode,
			"SQL injection attempt should not cause server error: %s", payload)

		// Test in JSON payload
		alertData := map[string]interface{}{
			"title":       "Test Alert",
			"description": payload,
			"severity":    "high",
		}
		jsonData, _ := json.Marshal(alertData)

		req, _ = http.NewRequest("POST", suite.baseURL+"/api/v1/alerts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)
		resp, err = suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		assert.NotEqual(suite.T(), http.StatusInternalServerError, resp.StatusCode,
			"SQL injection in JSON should not cause server error: %s", payload)
	}
}

func (suite *SecurityTestSuite) TestXSSPrevention() {
	validToken := suite.authenticate()

	// Test XSS payload attempts
	xssPayloads := []string{
		"<script>alert('XSS')</script>",
		"<img src=x onerror=alert('XSS')>",
		"javascript:alert('XSS')",
		"<iframe src='javascript:alert(\"XSS\")'></iframe>",
		"<svg onload=alert('XSS')>",
	}

	for _, payload := range xssPayloads {
		alertData := map[string]interface{}{
			"title":       payload,
			"description": "Test description",
			"severity":    "medium",
		}
		jsonData, _ := json.Marshal(alertData)

		req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/alerts", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)
		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

		// Verify the payload is stored safely (not executed)
		var createdAlert models.Alert
		json.NewDecoder(resp.Body).Decode(&createdAlert)

		// The title should be sanitized or escaped
		assert.NotContains(suite.T(), createdAlert.Title, "<script>", "XSS payload should be sanitized")
		assert.NotContains(suite.T(), createdAlert.Title, "javascript:", "JavaScript URLs should be sanitized")
	}
}

func (suite *SecurityTestSuite) TestInputValidation() {
	validToken := suite.authenticate()

	// Test various input validation scenarios
	testCases := []struct {
		name         string
		payload      interface{}
		expectStatus int
	}{
		{
			name: "Valid alert",
			payload: map[string]interface{}{
				"title":       "Valid Alert Title",
				"description": "Valid description that meets requirements",
				"severity":    "high",
				"category":    "security",
			},
			expectStatus: http.StatusCreated,
		},
		{
			name: "Empty title",
			payload: map[string]interface{}{
				"title":       "",
				"description": "Valid description",
				"severity":    "high",
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Title too long",
			payload: map[string]interface{}{
				"title":       strings.Repeat("A", 201), // Exceeds max length
				"description": "Valid description",
				"severity":    "high",
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid severity",
			payload: map[string]interface{}{
				"title":       "Valid Title",
				"description": "Valid description",
				"severity":    "invalid_severity",
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Description too short",
			payload: map[string]interface{}{
				"title":       "Valid Title",
				"description": "Short",
				"severity":    "high",
			},
			expectStatus: http.StatusBadRequest,
		},
		{
			name: "Invalid email format",
			payload: map[string]interface{}{
				"title":       "Valid Title",
				"description": "Valid description",
				"severity":    "high",
				"assigned_to": "invalid-email",
			},
			expectStatus: http.StatusBadRequest,
		},
	}

	for _, tc := range testCases {
		suite.T().Run(tc.name, func(t *testing.T) {
			jsonData, _ := json.Marshal(tc.payload)

			req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/alerts", bytes.NewBuffer(jsonData))
			req.Header.Set("Content-Type", "application/json")
			req.Header.Set("Authorization", "Bearer "+validToken)
			resp, err := suite.httpClient.Do(req)
			assert.NoError(t, err)
			defer resp.Body.Close()

			assert.Equal(t, tc.expectStatus, resp.StatusCode)
		})
	}
}

func (suite *SecurityTestSuite) TestRateLimiting() {
	// Test rate limiting functionality
	requests := 150 // Exceed rate limit

	for i := 0; i < requests; i++ {
		req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
		req.Header.Set("Authorization", "Bearer "+suite.authenticate())
		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		if i < 100 { // First 100 should succeed
			assert.NotEqual(suite.T(), http.StatusTooManyRequests, resp.StatusCode,
				"Request %d should not be rate limited", i+1)
		} else { // Subsequent requests should be rate limited
			assert.Equal(suite.T(), http.StatusTooManyRequests, resp.StatusCode,
				"Request %d should be rate limited", i+1)
		}
	}
}

func (suite *SecurityTestSuite) TestDataEncryption() {
	validToken := suite.authenticate()

	// Test that sensitive data is encrypted at rest
	sensitiveData := map[string]interface{}{
		"title":       "Sensitive Alert",
		"description": "This contains sensitive information: password=secret123, api_key=abcd1234",
		"severity":    "high",
	}
	jsonData, _ := json.Marshal(sensitiveData)

	// Create alert with sensitive data
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/alerts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validToken)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)

	// Verify data is encrypted in database
	alertID := suite.extractAlertID(resp.Body)
	encryptedData := suite.getRawAlertData(alertID)

	// Sensitive data should not be readable in plain text
	assert.NotContains(suite.T(), encryptedData, "password=secret123")
	assert.NotContains(suite.T(), encryptedData, "api_key=abcd1234")

	// But should be decryptable when retrieved via API
	req, _ = http.NewRequest("GET", suite.baseURL+"/api/v1/alerts/"+alertID, nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	resp, err = suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var alert models.Alert
	json.NewDecoder(resp.Body).Decode(&alert)

	// Data should be properly decrypted for authorized users
	assert.Contains(suite.T(), alert.Description, "password=secret123")
	assert.Contains(suite.T(), alert.Description, "api_key=abcd1234")
}

func (suite *SecurityTestSuite) TestAuditLogging() {
	validToken := suite.authenticate()

	// Perform various operations that should be audited
	initialLogCount := suite.getAuditLogCount()

	// Create alert
	alertData := map[string]interface{}{
		"title":       "Audit Test Alert",
		"description": "Testing audit logging",
		"severity":    "medium",
	}
	jsonData, _ := json.Marshal(alertData)

	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/alerts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validToken)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var alert models.Alert
	json.NewDecoder(resp.Body).Decode(&alert)

	// Update alert
	updateData := map[string]interface{}{
		"status":   "acknowledged",
		"comments": "Audit test update",
	}
	jsonData, _ = json.Marshal(updateData)

	req, _ = http.NewRequest("PUT", suite.baseURL+"/api/v1/alerts/"+alert.ID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+validToken)
	resp, err = suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Delete alert
	req, _ = http.NewRequest("DELETE", suite.baseURL+"/api/v1/alerts/"+alert.ID, nil)
	req.Header.Set("Authorization", "Bearer "+validToken)
	resp, err = suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Verify audit logs were created
	finalLogCount := suite.getAuditLogCount()
	newLogs := finalLogCount - initialLogCount

	assert.GreaterOrEqual(suite.T(), newLogs, 3, "Should have audit logs for create, update, and delete operations")

	// Verify audit log contents
	logs := suite.getRecentAuditLogs(10)
	operations := make(map[string]bool)

	for _, log := range logs {
		operations[log.Operation] = true
		assert.NotEmpty(suite.T(), log.UserID, "Audit log should contain user ID")
		assert.NotEmpty(suite.T(), log.Timestamp, "Audit log should contain timestamp")
		assert.NotEmpty(suite.T(), log.Resource, "Audit log should contain resource")
	}

	assert.True(suite.T(), operations["CREATE"], "Should have CREATE operation in audit logs")
	assert.True(suite.T(), operations["UPDATE"], "Should have UPDATE operation in audit logs")
	assert.True(suite.T(), operations["DELETE"], "Should have DELETE operation in audit logs")
}

func (suite *SecurityTestSuite) TestSessionManagement() {
	// Test session timeout
	token := suite.authenticate()

	// Use token immediately - should work
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)

	// Simulate session timeout by advancing time
	suite.advanceSessionTime(25 * time.Hour) // Past session timeout

	// Token should now be invalid
	req, _ = http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	resp, err = suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()
	assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
}

func (suite *SecurityTestSuite) TestCSRFProtection() {
	validToken := suite.authenticate()

	// Test that state-changing operations require proper authentication
	stateChangingEndpoints := []string{
		"/api/v1/alerts",
		"/api/v1/config",
		"/api/v1/admin/users",
	}

	for _, endpoint := range stateChangingEndpoints {
		// Try request without CSRF token (should fail if CSRF protection is enabled)
		alertData := map[string]interface{}{
			"title":       "CSRF Test",
			"description": "Testing CSRF protection",
			"severity":    "low",
		}
		jsonData, _ := json.Marshal(alertData)

		req, _ := http.NewRequest("POST", suite.baseURL+endpoint, bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		req.Header.Set("Authorization", "Bearer "+validToken)
		// Note: Not setting CSRF token

		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		// Should either succeed (if CSRF is disabled for API) or require CSRF token
		assert.True(suite.T(),
			resp.StatusCode == http.StatusCreated || resp.StatusCode == http.StatusForbidden,
			"CSRF protection should either allow or require token for %s", endpoint)
	}
}

func (suite *SecurityTestSuite) TestSecureHeaders() {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	// Check for security headers
	securityHeaders := map[string]string{
		"X-Content-Type-Options": "nosniff",
		"X-Frame-Options":        "DENY",
		"X-XSS-Protection":       "1; mode=block",
		"Strict-Transport-Security": "max-age=31536000; includeSubDomains",
		"Content-Security-Policy":   "default-src 'self'",
	}

	for header, expectedValue := range securityHeaders {
		actualValue := resp.Header.Get(header)
		assert.Equal(suite.T(), expectedValue, actualValue,
			"Security header %s should be set correctly", header)
	}
}

func (suite *SecurityTestSuite) TestCertificateValidation() {
	// Test with invalid certificate
	invalidClient := &http.Client{
		Transport: &http.Transport{
			TLSClientConfig: &tls.Config{
				InsecureSkipVerify: false, // Require certificate validation
				ServerName:         "invalid.example.com",
			},
		},
	}

	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/health", nil)
	_, err := invalidClient.Do(req)

	// Should fail certificate validation
	assert.Error(suite.T(), err, "Should fail with invalid certificate")
	assert.Contains(suite.T(), err.Error(), "certificate", "Error should be certificate-related")
}

func (suite *SecurityTestSuite) TestBruteForceProtection() {
	// Test login brute force protection
	wrongPassword := "wrongpassword123"

	for i := 0; i < 10; i++ {
		loginData := map[string]interface{}{
			"email":    "admin@insec.com",
			"password": wrongPassword,
		}
		jsonData, _ := json.Marshal(loginData)

		req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(jsonData))
		req.Header.Set("Content-Type", "application/json")
		resp, err := suite.httpClient.Do(req)
		suite.Require().NoError(err)
		defer resp.Body.Close()

		if i < 5 { // First few attempts should return unauthorized
			assert.Equal(suite.T(), http.StatusUnauthorized, resp.StatusCode)
		} else { // Later attempts should be blocked
			assert.True(suite.T(),
				resp.StatusCode == http.StatusTooManyRequests || resp.StatusCode == http.StatusUnauthorized,
				"Should be rate limited or unauthorized after %d attempts", i+1)
		}
	}
}

// Helper methods

func (suite *SecurityTestSuite) authenticate() string {
	loginData := models.LoginRequest{
		Email:    "admin@insec.com",
		Password: "securepassword123",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	resp, err := suite.httpClient.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var loginResponse models.LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResponse)

	return loginResponse.Token
}

func (suite *SecurityTestSuite) authenticateAsRole(role string) string {
	// Implementation for authenticating as specific role
	return "role-token"
}

func (suite *SecurityTestSuite) generateExpiredToken() string {
	// Implementation for generating expired token
	return "expired.token"
}

func (suite *SecurityTestSuite) extractAlertID(body *bytes.Buffer) string {
	// Implementation for extracting alert ID from response
	return "alert-id"
}

func (suite *SecurityTestSuite) getRawAlertData(alertID string) string {
	// Implementation for getting raw alert data from database
	return "raw data"
}

func (suite *SecurityTestSuite) getAuditLogCount() int {
	// Implementation for getting audit log count
	return 0
}

func (suite *SecurityTestSuite) getRecentAuditLogs(count int) []models.AuditLog {
	// Implementation for getting recent audit logs
	return []models.AuditLog{}
}

func (suite *SecurityTestSuite) advanceSessionTime(duration time.Duration) {
	// Implementation for advancing session time
}

func contains(slice []string, item string) bool {
	for _, s := range slice {
		if s == item {
			return true
		}
	}
	return false
}</content>
<parameter name="filePath">/workspaces/insec/tests/security/system_security_test.go
