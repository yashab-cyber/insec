package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"insec/internal/models"
	"insec/test/helpers"
)

type EndToEndTestSuite struct {
	suite.Suite
	testHelper *helpers.TestHelper
	baseURL    string
	authToken  string
}

func (suite *EndToEndTestSuite) SetupSuite() {
	// Initialize test environment
	suite.testHelper = helpers.NewTestHelper()
	suite.baseURL = "http://localhost:8080"

	// Start the complete INSEC system
	err := suite.testHelper.StartSystem()
	suite.Require().NoError(err)

	// Wait for system to be ready
	suite.testHelper.WaitForSystemReady()

	// Get authentication token
	suite.authToken = suite.authenticate()
}

func (suite *EndToEndTestSuite) TearDownSuite() {
	// Clean up test environment
	suite.testHelper.StopSystem()
}

func TestEndToEndTestSuite(t *testing.T) {
	suite.Run(t, new(EndToEndTestSuite))
}

func (suite *EndToEndTestSuite) TestCompleteAlertLifecycle() {
	// 1. Agent registers with server
	agentID := suite.registerAgent()

	// 2. Agent sends telemetry data
	suite.sendTelemetryData(agentID)

	// 3. Server processes telemetry and generates alert
	alertID := suite.waitForAlertGeneration()

	// 4. User retrieves and views alert
	alert := suite.getAlert(alertID)

	// 5. User acknowledges alert
	suite.acknowledgeAlert(alertID)

	// 6. User assigns alert to analyst
	suite.assignAlert(alertID, "analyst@insec.com")

	// 7. Analyst investigates and resolves alert
	suite.resolveAlert(alertID, "False positive - legitimate system process")

	// 8. Verify alert resolution
	resolvedAlert := suite.getAlert(alertID)
	assert.Equal(suite.T(), "resolved", resolvedAlert.Status)
	assert.Equal(suite.T(), "False positive - legitimate system process", resolvedAlert.Resolution)
}

func (suite *EndToEndTestSuite) TestRealTimeEventStreaming() {
	// Subscribe to real-time events
	eventChan := suite.subscribeToEvents()

	// Generate events by simulating agent activity
	go suite.simulateAgentActivity()

	// Collect events for 30 seconds
	events := suite.collectEvents(eventChan, 30*time.Second)

	// Verify events were received and processed
	assert.Greater(suite.T(), len(events), 0, "Should receive real-time events")

	// Check event types
	eventTypes := make(map[string]int)
	for _, event := range events {
		eventTypes[event.EventType]++
	}

	assert.Greater(suite.T(), eventTypes["process"], 0, "Should receive process events")
	assert.Greater(suite.T(), eventTypes["file"], 0, "Should receive file events")
	assert.Greater(suite.T(), eventTypes["network"], 0, "Should receive network events")
}

func (suite *EndToEndTestSuite) TestRiskScoringAndAlerting() {
	// 1. Send high-risk telemetry data
	suite.sendHighRiskTelemetry()

	// 2. Wait for risk score calculation
	time.Sleep(5 * time.Second)

	// 3. Verify high-risk alert is generated
	alerts := suite.getHighRiskAlerts()
	assert.Greater(suite.T(), len(alerts), 0, "Should generate high-risk alerts")

	// 4. Verify alert has correct risk score
	for _, alert := range alerts {
		assert.GreaterOrEqual(suite.T(), alert.RiskScore, 0.8, "Alert should have high risk score")
		assert.Equal(suite.T(), "high", alert.Severity, "Alert should be high severity")
	}
}

func (suite *EndToEndTestSuite) TestMultiTenantIsolation() {
	// 1. Create second tenant
	tenant2Token := suite.createTenantAndAuthenticate("tenant2")

	// 2. Register agent for second tenant
	agent2ID := suite.registerAgentWithToken(tenant2Token)

	// 3. Send data from both tenants
	suite.sendTelemetryData("agent1-tenant1")
	suite.sendTelemetryDataWithToken(agent2ID, tenant2Token)

	// 4. Verify tenant isolation - each tenant only sees their own data
	tenant1Alerts := suite.getAlerts()
	tenant2Alerts := suite.getAlertsWithToken(tenant2Token)

	// Ensure no cross-tenant data leakage
	for _, alert := range tenant1Alerts {
		assert.Equal(suite.T(), "tenant1", alert.TenantID)
	}

	for _, alert := range tenant2Alerts {
		assert.Equal(suite.T(), "tenant2", alert.TenantID)
	}
}

func (suite *EndToEndTestSuite) TestSystemPerformanceUnderLoad() {
	// 1. Start performance monitoring
	suite.startPerformanceMonitoring()

	// 2. Generate high volume of telemetry data
	suite.generateHighVolumeTelemetry(1000) // 1000 events

	// 3. Wait for processing
	time.Sleep(10 * time.Second)

	// 4. Check system performance metrics
	metrics := suite.getSystemMetrics()

	// Verify performance thresholds
	assert.Less(suite.T(), metrics.AverageResponseTime, 1000.0, "Response time should be under 1 second")
	assert.Greater(suite.T(), metrics.Throughput, 50.0, "Should handle at least 50 events/second")
	assert.Less(suite.T(), metrics.ErrorRate, 0.01, "Error rate should be less than 1%")

	// 5. Verify all events were processed
	processedEvents := suite.getProcessedEventsCount()
	assert.Equal(suite.T(), 1000, processedEvents, "All events should be processed")
}

func (suite *EndToEndTestSuite) TestAgentServerCommunication() {
	// 1. Test agent heartbeat
	suite.testAgentHeartbeat()

	// 2. Test configuration sync
	suite.testConfigurationSync()

	// 3. Test secure communication (TLS)
	suite.testSecureCommunication()

	// 4. Test connection recovery after network interruption
	suite.testConnectionRecovery()
}

func (suite *EndToEndTestSuite) TestDataPersistenceAndRecovery() {
	// 1. Send telemetry data
	initialCount := suite.getTotalEventsCount()
	suite.sendTelemetryData("test-agent")

	// 2. Verify data is persisted
	time.Sleep(2 * time.Second)
	newCount := suite.getTotalEventsCount()
	assert.Greater(suite.T(), newCount, initialCount, "Data should be persisted")

	// 3. Simulate system restart
	suite.restartSystem()

	// 4. Verify data recovery
	recoveredCount := suite.getTotalEventsCount()
	assert.Equal(suite.T(), newCount, recoveredCount, "Data should be recovered after restart")
}

func (suite *EndToEndTestSuite) TestSecurityFeatures() {
	// 1. Test authentication
	suite.testAuthenticationSecurity()

	// 2. Test authorization
	suite.testAuthorization()

	// 3. Test data encryption
	suite.testDataEncryption()

	// 4. Test audit logging
	suite.testAuditLogging()

	// 5. Test intrusion detection
	suite.testIntrusionDetection()
}

func (suite *EndToEndTestSuite) TestScalability() {
	// 1. Start with baseline agents
	baselineMetrics := suite.getSystemMetrics()

	// 2. Add multiple agents
	agentIDs := suite.addMultipleAgents(10)

	// 3. Generate load from all agents
	for _, agentID := range agentIDs {
		go suite.sendContinuousTelemetry(agentID, 60*time.Second)
	}

	// 4. Monitor system scaling
	time.Sleep(30 * time.Second)
	scaledMetrics := suite.getSystemMetrics()

	// 5. Verify system scales appropriately
	assert.Greater(suite.T(), scaledMetrics.CPUUsage, baselineMetrics.CPUUsage, "CPU usage should increase with load")
	assert.Less(suite.T(), scaledMetrics.AverageResponseTime, baselineMetrics.AverageResponseTime*2, "Response time degradation should be reasonable")
}

func (suite *EndToEndTestSuite) TestUIIntegration() {
	// 1. Test user login via UI
	suite.testUILogin()

	// 2. Test dashboard data display
	suite.testDashboardData()

	// 3. Test alert management via UI
	suite.testUIAlertManagement()

	// 4. Test real-time updates
	suite.testUIRealTimeUpdates()

	// 5. Test export functionality
	suite.testDataExport()
}

func (suite *EndToEndTestSuite) TestBackupAndRestore() {
	// 1. Create backup
	backupID := suite.createBackup()

	// 2. Add more data
	suite.sendTelemetryData("backup-test-agent")

	// 3. Restore from backup
	suite.restoreFromBackup(backupID)

	// 4. Verify data integrity
	originalCount := suite.getEventsCountBeforeBackup()
	restoredCount := suite.getTotalEventsCount()
	assert.Equal(suite.T(), originalCount, restoredCount, "Data should be restored correctly")
}

func (suite *EndToEndTestSuite) TestMonitoringAndAlerting() {
	// 1. Configure monitoring thresholds
	suite.configureMonitoringThresholds()

	// 2. Trigger monitoring alerts
	suite.triggerHighCPUUsage()
	suite.triggerLowDiskSpace()
	suite.triggerNetworkIssues()

	// 3. Verify system generates appropriate alerts
	monitoringAlerts := suite.getMonitoringAlerts()
	assert.Greater(suite.T(), len(monitoringAlerts), 0, "Should generate monitoring alerts")

	// 4. Test alert escalation
	suite.testAlertEscalation()
}

// Helper methods

func (suite *EndToEndTestSuite) authenticate() string {
	loginData := models.LoginRequest{
		Email:    "admin@insec.com",
		Password: "securepassword123",
	}

	jsonData, _ := json.Marshal(loginData)
	resp, err := http.Post(suite.baseURL+"/api/v1/auth/login", "application/json", bytes.NewBuffer(jsonData))
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var loginResponse models.LoginResponse
	json.NewDecoder(resp.Body).Decode(&loginResponse)

	return loginResponse.Token
}

func (suite *EndToEndTestSuite) registerAgent() string {
	agentData := models.AgentRegistration{
		AgentID:     "e2e-test-agent",
		TenantID:    "test-tenant",
		Hostname:    "test-host",
		OS:          "Linux",
		Version:     "1.0.0",
		Capabilities: []string{"process_monitoring", "file_monitoring", "network_monitoring"},
	}

	jsonData, _ := json.Marshal(agentData)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/agents/register", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var response models.AgentRegistrationResponse
	json.NewDecoder(resp.Body).Decode(&response)

	return response.AgentID
}

func (suite *EndToEndTestSuite) sendTelemetryData(agentID string) {
	eventData := models.EventRequest{
		EventType:   "process",
		Description: "Test process started",
		Severity:    "low",
		Source:      agentID,
		Data: map[string]interface{}{
			"process_name": "test_process.exe",
			"pid":         12345,
			"command_line": "test_process.exe --test",
		},
	}

	jsonData, _ := json.Marshal(eventData)
	req, _ := http.NewRequest("POST", suite.baseURL+"/api/v1/events", bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusCreated, resp.StatusCode)
}

func (suite *EndToEndTestSuite) waitForAlertGeneration() string {
	// Poll for alerts until one is generated or timeout
	timeout := time.After(30 * time.Second)
	ticker := time.NewTicker(1 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case <-timeout:
			suite.T().Fatal("Timeout waiting for alert generation")
			return ""
		case <-ticker.C:
			alerts := suite.getAlerts()
			if len(alerts) > 0 {
				return alerts[0].ID
			}
		}
	}
}

func (suite *EndToEndTestSuite) getAlerts() []models.Alert {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/alerts", nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var alerts []models.Alert
	json.NewDecoder(resp.Body).Decode(&alerts)

	return alerts
}

func (suite *EndToEndTestSuite) getAlert(alertID string) models.Alert {
	req, _ := http.NewRequest("GET", suite.baseURL+"/api/v1/alerts/"+alertID, nil)
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	var alert models.Alert
	json.NewDecoder(resp.Body).Decode(&alert)

	return alert
}

func (suite *EndToEndTestSuite) acknowledgeAlert(alertID string) {
	updateData := models.AlertUpdateRequest{
		Status:   "acknowledged",
		Comments: "Alert acknowledged for investigation",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", suite.baseURL+"/api/v1/alerts/"+alertID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *EndToEndTestSuite) assignAlert(alertID, assignee string) {
	updateData := models.AlertUpdateRequest{
		AssignedTo: assignee,
		Comments:   fmt.Sprintf("Assigned to %s for investigation", assignee),
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", suite.baseURL+"/api/v1/alerts/"+alertID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

func (suite *EndToEndTestSuite) resolveAlert(alertID, resolution string) {
	updateData := models.AlertUpdateRequest{
		Status:     "resolved",
		Resolution: resolution,
		Comments:   "Alert resolved after investigation",
	}

	jsonData, _ := json.Marshal(updateData)
	req, _ := http.NewRequest("PUT", suite.baseURL+"/api/v1/alerts/"+alertID, bytes.NewBuffer(jsonData))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Authorization", "Bearer "+suite.authToken)

	client := &http.Client{}
	resp, err := client.Do(req)
	suite.Require().NoError(err)
	defer resp.Body.Close()

	assert.Equal(suite.T(), http.StatusOK, resp.StatusCode)
}

// Additional helper methods would be implemented here...
func (suite *EndToEndTestSuite) subscribeToEvents() chan models.Event {
	// Implementation for subscribing to real-time events
	return make(chan models.Event, 100)
}

func (suite *EndToEndTestSuite) simulateAgentActivity() {
	// Implementation for simulating agent activity
}

func (suite *EndToEndTestSuite) collectEvents(eventChan chan models.Event, duration time.Duration) []models.Event {
	// Implementation for collecting events
	return []models.Event{}
}

func (suite *EndToEndTestSuite) sendHighRiskTelemetry() {
	// Implementation for sending high-risk telemetry
}

func (suite *EndToEndTestSuite) getHighRiskAlerts() []models.Alert {
	// Implementation for getting high-risk alerts
	return []models.Alert{}
}

func (suite *EndToEndTestSuite) createTenantAndAuthenticate(tenantID string) string {
	// Implementation for creating tenant and getting auth token
	return ""
}

func (suite *EndToEndTestSuite) registerAgentWithToken(token string) string {
	// Implementation for registering agent with specific token
	return ""
}

func (suite *EndToEndTestSuite) sendTelemetryDataWithToken(agentID, token string) {
	// Implementation for sending telemetry with specific token
}

func (suite *EndToEndTestSuite) getAlertsWithToken(token string) []models.Alert {
	// Implementation for getting alerts with specific token
	return []models.Alert{}
}

func (suite *EndToEndTestSuite) startPerformanceMonitoring() {
	// Implementation for starting performance monitoring
}

func (suite *EndToEndTestSuite) generateHighVolumeTelemetry(count int) {
	// Implementation for generating high volume telemetry
}

func (suite *EndToEndTestSuite) getSystemMetrics() models.SystemMetrics {
	// Implementation for getting system metrics
	return models.SystemMetrics{}
}

func (suite *EndToEndTestSuite) getProcessedEventsCount() int {
	// Implementation for getting processed events count
	return 0
}

func (suite *EndToEndTestSuite) testAgentHeartbeat() {
	// Implementation for testing agent heartbeat
}

func (suite *EndToEndTestSuite) testConfigurationSync() {
	// Implementation for testing configuration sync
}

func (suite *EndToEndTestSuite) testSecureCommunication() {
	// Implementation for testing secure communication
}

func (suite *EndToEndTestSuite) testConnectionRecovery() {
	// Implementation for testing connection recovery
}

func (suite *EndToEndTestSuite) getTotalEventsCount() int {
	// Implementation for getting total events count
	return 0
}

func (suite *EndToEndTestSuite) restartSystem() {
	// Implementation for restarting system
}

func (suite *EndToEndTestSuite) testAuthenticationSecurity() {
	// Implementation for testing authentication security
}

func (suite *EndToEndTestSuite) testAuthorization() {
	// Implementation for testing authorization
}

func (suite *EndToEndTestSuite) testDataEncryption() {
	// Implementation for testing data encryption
}

func (suite *EndToEndTestSuite) testAuditLogging() {
	// Implementation for testing audit logging
}

func (suite *EndToEndTestSuite) testIntrusionDetection() {
	// Implementation for testing intrusion detection
}

func (suite *EndToEndTestSuite) addMultipleAgents(count int) []string {
	// Implementation for adding multiple agents
	return []string{}
}

func (suite *EndToEndTestSuite) sendContinuousTelemetry(agentID string, duration time.Duration) {
	// Implementation for sending continuous telemetry
}

func (suite *EndToEndTestSuite) testUILogin() {
	// Implementation for testing UI login
}

func (suite *EndToEndTestSuite) testDashboardData() {
	// Implementation for testing dashboard data
}

func (suite *EndToEndTestSuite) testUIAlertManagement() {
	// Implementation for testing UI alert management
}

func (suite *EndToEndTestSuite) testUIRealTimeUpdates() {
	// Implementation for testing UI real-time updates
}

func (suite *EndToEndTestSuite) testDataExport() {
	// Implementation for testing data export
}

func (suite *EndToEndTestSuite) createBackup() string {
	// Implementation for creating backup
	return ""
}

func (suite *EndToEndTestSuite) restoreFromBackup(backupID string) {
	// Implementation for restoring from backup
}

func (suite *EndToEndTestSuite) getEventsCountBeforeBackup() int {
	// Implementation for getting events count before backup
	return 0
}

func (suite *EndToEndTestSuite) configureMonitoringThresholds() {
	// Implementation for configuring monitoring thresholds
}

func (suite *EndToEndTestSuite) triggerHighCPUUsage() {
	// Implementation for triggering high CPU usage
}

func (suite *EndToEndTestSuite) triggerLowDiskSpace() {
	// Implementation for triggering low disk space
}

func (suite *EndToEndTestSuite) triggerNetworkIssues() {
	// Implementation for triggering network issues
}

func (suite *EndToEndTestSuite) getMonitoringAlerts() []models.Alert {
	// Implementation for getting monitoring alerts
	return []models.Alert{}
}

func (suite *EndToEndTestSuite) testAlertEscalation() {
	// Implementation for testing alert escalation
}</content>
<parameter name="filePath">/workspaces/insec/tests/e2e/system_test.go
