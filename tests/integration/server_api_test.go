package tests

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"insec/internal/auth"
	"insec/internal/models"
	"insec/internal/server"
)

type ServerIntegrationTestSuite struct {
	suite.Suite
	router *gin.Engine
	auth   *auth.AuthService
}

func (suite *ServerIntegrationTestSuite) SetupTest() {
	gin.SetMode(gin.TestMode)
	suite.router = server.SetupRouter()
	suite.auth = auth.NewAuthService()
}

func (suite *ServerIntegrationTestSuite) TearDownTest() {
	// Clean up test data
}

func TestServerIntegrationTestSuite(t *testing.T) {
	suite.Run(t, new(ServerIntegrationTestSuite))
}

// Test authentication endpoints
func (suite *ServerIntegrationTestSuite) TestAuthenticationFlow() {
	// Test login endpoint
	loginData := models.LoginRequest{
		Email:    "admin@insec.com",
		Password: "securepassword123",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var loginResponse models.LoginResponse
	err := json.Unmarshal(w.Body.Bytes(), &loginResponse)
	assert.NoError(suite.T(), err)
	assert.NotEmpty(suite.T(), loginResponse.Token)
	assert.NotEmpty(suite.T(), loginResponse.User.ID)
	assert.Equal(suite.T(), "admin@insec.com", loginResponse.User.Email)
}

func (suite *ServerIntegrationTestSuite) TestAuthenticationMiddleware() {
	// First, get a valid token
	token := suite.getValidToken()

	// Test protected endpoint without token
	req, _ := http.NewRequest("GET", "/api/v1/alerts", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)

	// Test protected endpoint with valid token
	req, _ = http.NewRequest("GET", "/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

func (suite *ServerIntegrationTestSuite) TestInvalidToken() {
	// Test with invalid token
	req, _ := http.NewRequest("GET", "/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer invalid.token.here")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

func (suite *ServerIntegrationTestSuite) TestExpiredToken() {
	// Test with expired token (this would require mocking time or using a pre-expired token)
	req, _ := http.NewRequest("GET", "/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer expired.token.here")
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusUnauthorized, w.Code)
}

// Test alerts endpoints
func (suite *ServerIntegrationTestSuite) TestCreateAlert() {
	token := suite.getValidToken()

	alertData := models.AlertRequest{
		Title:       "Test Alert",
		Description: "This is a test alert for integration testing",
		Severity:    "high",
		Category:    "security",
		Source:      "integration-test",
	}

	jsonData, _ := json.Marshal(alertData)
	req, _ := http.NewRequest("POST", "/api/v1/alerts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var alert models.Alert
	err := json.Unmarshal(w.Body.Bytes(), &alert)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "Test Alert", alert.Title)
	assert.Equal(suite.T(), "high", alert.Severity)
	assert.NotEmpty(suite.T(), alert.ID)
}

func (suite *ServerIntegrationTestSuite) TestGetAlerts() {
	token := suite.getValidToken()

	req, _ := http.NewRequest("GET", "/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var alerts []models.Alert
	err := json.Unmarshal(w.Body.Bytes(), &alerts)
	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), []models.Alert{}, alerts)
}

func (suite *ServerIntegrationTestSuite) TestGetAlertByID() {
	token := suite.getValidToken()

	// First create an alert
	alertID := suite.createTestAlert(token)

	// Then retrieve it
	req, _ := http.NewRequest("GET", "/api/v1/alerts/"+alertID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var alert models.Alert
	err := json.Unmarshal(w.Body.Bytes(), &alert)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), alertID, alert.ID)
}

func (suite *ServerIntegrationTestSuite) TestUpdateAlert() {
	token := suite.getValidToken()
	alertID := suite.createTestAlert(token)

	updateData := models.AlertUpdateRequest{
		Status:      "acknowledged",
		AssignedTo:  "analyst@insec.com",
		Priority:    "high",
		Comments:    "Alert acknowledged and assigned to analyst",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", "/api/v1/alerts/"+alertID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var updatedAlert models.Alert
	err := json.Unmarshal(w.Body.Bytes(), &updatedAlert)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "acknowledged", updatedAlert.Status)
	assert.Equal(suite.T(), "analyst@insec.com", updatedAlert.AssignedTo)
}

func (suite *ServerIntegrationTestSuite) TestDeleteAlert() {
	token := suite.getValidToken()
	alertID := suite.createTestAlert(token)

	req, _ := http.NewRequest("DELETE", "/api/v1/alerts/"+alertID, nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNoContent, w.Code)

	// Verify alert is deleted
	req, _ = http.NewRequest("GET", "/api/v1/alerts/"+alertID, nil)
	req.Header.Set("Authorization", "Bearer "+token)
	w = httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)
	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

// Test events endpoints
func (suite *ServerIntegrationTestSuite) TestCreateEvent() {
	token := suite.getValidToken()

	eventData := models.EventRequest{
		EventType:   "process",
		Description: "Process started",
		Severity:    "medium",
		Source:      "agent-123",
		Data: map[string]interface{}{
			"process_name": "test.exe",
			"pid":         1234,
		},
	}

	jsonData, _ := json.Marshal(eventData)
	req, _ := http.NewRequest("POST", "/api/v1/events", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var event models.Event
	err := json.Unmarshal(w.Body.Bytes(), &event)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "process", event.EventType)
	assert.Equal(suite.T(), "test.exe", event.Data["process_name"])
}

func (suite *ServerIntegrationTestSuite) TestGetEventsWithFilters() {
	token := suite.getValidToken()

	// Test with query parameters
	req, _ := http.NewRequest("GET", "/api/v1/events?event_type=process&severity=high&limit=10", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var events []models.Event
	err := json.Unmarshal(w.Body.Bytes(), &events)
	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), []models.Event{}, events)
}

// Test analytics endpoints
func (suite *ServerIntegrationTestSuite) TestGetAnalytics() {
	token := suite.getValidToken()

	req, _ := http.NewRequest("GET", "/api/v1/analytics/summary", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var analytics models.AnalyticsSummary
	err := json.Unmarshal(w.Body.Bytes(), &analytics)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), analytics.TotalAlerts)
	assert.NotNil(suite.T(), analytics.TotalEvents)
}

func (suite *ServerIntegrationTestSuite) TestGetRiskMetrics() {
	token := suite.getValidToken()

	req, _ := http.NewRequest("GET", "/api/v1/analytics/risk-metrics", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var metrics models.RiskMetrics
	err := json.Unmarshal(w.Body.Bytes(), &metrics)
	assert.NoError(suite.T(), err)
	assert.IsType(suite.T(), models.RiskMetrics{}, metrics)
}

// Test agent endpoints
func (suite *ServerIntegrationTestSuite) TestAgentRegistration() {
	agentData := models.AgentRegistration{
		AgentID:     "test-agent-123",
		TenantID:    "test-tenant",
		Hostname:    "test-host",
		OS:          "Linux",
		Version:     "1.0.0",
		Capabilities: []string{"process_monitoring", "file_monitoring"},
	}

	jsonData, _ := json.Marshal(agentData)
	req, _ := http.NewRequest("POST", "/api/v1/agents/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusCreated, w.Code)

	var response models.AgentRegistrationResponse
	err := json.Unmarshal(w.Body.Bytes(), &response)
	assert.NoError(suite.T(), err)
	assert.Equal(suite.T(), "test-agent-123", response.AgentID)
	assert.NotEmpty(suite.T(), response.Token)
}

func (suite *ServerIntegrationTestSuite) TestAgentHeartbeat() {
	token := suite.getValidToken()

	heartbeatData := models.AgentHeartbeat{
		AgentID:        "test-agent-123",
		Timestamp:      time.Now(),
		Status:         "healthy",
		Version:        "1.0.0",
		UptimeSeconds:  3600,
		MemoryUsageMB:  150.5,
		CPUUsagePercent: 25.3,
	}

	jsonData, _ := json.Marshal(heartbeatData)
	req, _ := http.NewRequest("POST", "/api/v1/agents/heartbeat", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Test configuration endpoints
func (suite *ServerIntegrationTestSuite) TestGetConfiguration() {
	token := suite.getValidToken()

	req, _ := http.NewRequest("GET", "/api/v1/config", nil)
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)

	var config models.Configuration
	err := json.Unmarshal(w.Body.Bytes(), &config)
	assert.NoError(suite.T(), err)
	assert.NotNil(suite.T(), config)
}

func (suite *ServerIntegrationTestSuite) TestUpdateConfiguration() {
	token := suite.getValidToken()

	configUpdate := models.ConfigurationUpdate{
		CollectionInterval: 60,
		MaxBatchSize:       200,
		EnableCompression:  true,
		LogLevel:          "debug",
	}

	jsonData, _ := json.Marshal(configUpdate)
	req, _ := http.NewRequest("PUT", "/api/v1/config", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusOK, w.Code)
}

// Test error handling
func (suite *ServerIntegrationTestSuite) TestInvalidJSON() {
	token := suite.getValidToken()

	req, _ := http.NewRequest("POST", "/api/v1/alerts", bytes.NewBuffer([]byte("invalid json")))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusBadRequest, w.Code)
}

func (suite *ServerIntegrationTestSuite) TestInvalidRoute() {
	req, _ := http.NewRequest("GET", "/api/v1/nonexistent", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusNotFound, w.Code)
}

func (suite *ServerIntegrationTestSuite) TestMethodNotAllowed() {
	req, _ := http.NewRequest("PATCH", "/api/v1/alerts", nil)
	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	assert.Equal(suite.T(), http.StatusMethodNotAllowed, w.Code)
}

// Test rate limiting
func (suite *ServerIntegrationTestSuite) TestRateLimiting() {
	token := suite.getValidToken()

	// Make multiple rapid requests
	for i := 0; i < 150; i++ {
		req, _ := http.NewRequest("GET", "/api/v1/alerts", nil)
		req.Header.Set("Authorization", "Bearer "+token)
		w := httptest.NewRecorder()
		suite.router.ServeHTTP(w, req)

		if i < 100 { // First 100 should succeed
			assert.Equal(suite.T(), http.StatusOK, w.Code)
		} else if i >= 100 { // Subsequent requests should be rate limited
			assert.Equal(suite.T(), http.StatusTooManyRequests, w.Code)
		}
	}
}

// Helper methods
func (suite *ServerIntegrationTestSuite) getValidToken() string {
	loginData := models.LoginRequest{
		Email:    "admin@insec.com",
		Password: "securepassword123",
	}

	jsonData, _ := json.Marshal(loginData)
	req, _ := http.NewRequest("POST", "/api/v1/auth/login", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var loginResponse models.LoginResponse
	json.Unmarshal(w.Body.Bytes(), &loginResponse)

	return loginResponse.Token
}

func (suite *ServerIntegrationTestSuite) createTestAlert(token string) string {
	alertData := models.AlertRequest{
		Title:       "Integration Test Alert",
		Description: "Alert created during integration testing",
		Severity:    "medium",
		Category:    "test",
		Source:      "integration-test",
	}

	jsonData, _ := json.Marshal(alertData)
	req, _ := http.NewRequest("POST", "/api/v1/alerts", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+token)

	w := httptest.NewRecorder()
	suite.router.ServeHTTP(w, req)

	var alert models.Alert
	json.Unmarshal(w.Body.Bytes(), &alert)

	return alert.ID
}</content>
<parameter name="filePath">/workspaces/insec/tests/integration/server_api_test.go
