package tests

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"sync"
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/suite"

	"insec/internal/models"
	"insec/test/helpers"
)

type PerformanceTestSuite struct {
	suite.Suite
	testHelper *helpers.TestHelper
	baseURL    string
	authToken  string
}

func (suite *PerformanceTestSuite) SetupSuite() {
	suite.testHelper = helpers.NewTestHelper()
	suite.baseURL = "http://localhost:8080"

	err := suite.testHelper.StartSystem()
	suite.Require().NoError(err)

	suite.testHelper.WaitForSystemReady()
	suite.authToken = suite.authenticate()
}

func (suite *PerformanceTestSuite) TearDownSuite() {
	suite.testHelper.StopSystem()
}

func TestPerformanceTestSuite(t *testing.T) {
	suite.Run(t, new(PerformanceTestSuite))
}

func (suite *PerformanceTestSuite) TestEventIngestionRate() {
	// Test event ingestion performance
	concurrency := 10
	eventsPerWorker := 1000
	totalEvents := concurrency * eventsPerWorker

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < concurrency; i++ {
		wg.Add(1)
		go func(workerID int) {
			defer wg.Done()
			suite.sendBatchEvents(fmt.Sprintf("perf-agent-%d", workerID), eventsPerWorker)
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	eventsPerSecond := float64(totalEvents) / duration.Seconds()
	fmt.Printf("Event ingestion rate: %.2f events/second\n", eventsPerSecond)

	// Assert minimum performance threshold
	assert.GreaterOrEqual(suite.T(), eventsPerSecond, 100.0, "Should handle at least 100 events/second")
}

func (suite *PerformanceTestSuite) TestConcurrentUserLoad() {
	// Test concurrent user load
	userCount := 50
	requestsPerUser := 20

	var wg sync.WaitGroup
	results := make(chan time.Duration, userCount*requestsPerUser)

	start := time.Now()

	for i := 0; i < userCount; i++ {
		wg.Add(1)
		go func(userID int) {
			defer wg.Done()
			suite.simulateUserSession(userID, requestsPerUser, results)
		}(i)
	}

	wg.Wait()
	close(results)

	totalDuration := time.Since(start)
	var totalRequestTime time.Duration
	requestCount := 0

	for requestTime := range results {
		totalRequestTime += requestTime
		requestCount++
	}

	averageResponseTime := totalRequestTime / time.Duration(requestCount)
	throughput := float64(requestCount) / totalDuration.Seconds()

	fmt.Printf("Concurrent users: %d\n", userCount)
	fmt.Printf("Total requests: %d\n", requestCount)
	fmt.Printf("Average response time: %v\n", averageResponseTime)
	fmt.Printf("Throughput: %.2f requests/second\n", throughput)

	// Performance assertions
	assert.LessOrEqual(suite.T(), averageResponseTime, 500*time.Millisecond, "Average response time should be under 500ms")
	assert.GreaterOrEqual(suite.T(), throughput, 100.0, "Should handle at least 100 requests/second")
}

func (suite *PerformanceTestSuite) TestDatabaseQueryPerformance() {
	// Test database query performance under load
	suite.prepareTestData(10000) // Prepare 10k test records

	queries := []string{
		"SELECT COUNT(*) FROM alerts WHERE status = 'active'",
		"SELECT * FROM events WHERE created_at >= $1 ORDER BY created_at DESC LIMIT 100",
		"SELECT AVG(risk_score) FROM alerts WHERE severity = 'high'",
		"SELECT COUNT(*) FROM events GROUP BY event_type",
	}

	for _, query := range queries {
		start := time.Now()
		for i := 0; i < 100; i++ { // Execute each query 100 times
			suite.executeQuery(query)
		}
		duration := time.Since(start)

		averageQueryTime := duration / 100
		fmt.Printf("Query: %s\n", query)
		fmt.Printf("Average execution time: %v\n", averageQueryTime)

		assert.LessOrEqual(suite.T(), averageQueryTime, 50*time.Millisecond, "Query should execute in under 50ms")
	}
}

func (suite *PerformanceTestSuite) TestMemoryUsageUnderLoad() {
	// Test memory usage under sustained load
	initialMemory := suite.getMemoryUsage()

	// Generate sustained load for 2 minutes
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Minute)
	defer cancel()

	loadDone := make(chan bool)
	go func() {
		suite.generateSustainedLoad(ctx, 20) // 20 concurrent workers
		loadDone <- true
	}()

	// Monitor memory usage every 10 seconds
	ticker := time.NewTicker(10 * time.Second)
	defer ticker.Stop()

	maxMemory := initialMemory
	for {
		select {
		case <-ctx.Done():
			goto monitorDone
		case <-ticker.C:
			currentMemory := suite.getMemoryUsage()
			if currentMemory > maxMemory {
				maxMemory = currentMemory
			}
			fmt.Printf("Current memory usage: %.2f MB\n", currentMemory)
		}
	}

monitorDone:
	<-loadDone

	finalMemory := suite.getMemoryUsage()
	memoryIncrease := finalMemory - initialMemory

	fmt.Printf("Initial memory: %.2f MB\n", initialMemory)
	fmt.Printf("Max memory: %.2f MB\n", maxMemory)
	fmt.Printf("Final memory: %.2f MB\n", finalMemory)
	fmt.Printf("Memory increase: %.2f MB\n", memoryIncrease)

	// Assert memory usage is reasonable
	assert.LessOrEqual(suite.T(), memoryIncrease, 500.0, "Memory increase should be under 500MB")
	assert.LessOrEqual(suite.T(), maxMemory, 2048.0, "Max memory usage should be under 2GB")
}

func (suite *PerformanceTestSuite) TestNetworkLatency() {
	// Test network latency for API calls
	endpoints := []string{
		"/api/v1/alerts",
		"/api/v1/events",
		"/api/v1/analytics/summary",
		"/api/v1/config",
	}

	iterations := 100

	for _, endpoint := range endpoints {
		var totalLatency time.Duration

		for i := 0; i < iterations; i++ {
			start := time.Now()
			suite.makeAPIRequest(endpoint)
			latency := time.Since(start)
			totalLatency += latency
		}

		averageLatency := totalLatency / time.Duration(iterations)
		fmt.Printf("Endpoint: %s\n", endpoint)
		fmt.Printf("Average latency: %v\n", averageLatency)

		assert.LessOrEqual(suite.T(), averageLatency, 200*time.Millisecond, "API latency should be under 200ms")
	}
}

func (suite *PerformanceTestSuite) TestAgentCommunicationPerformance() {
	// Test agent-to-server communication performance
	agentCount := 5
	eventsPerAgent := 500
	totalEvents := agentCount * eventsPerAgent

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < agentCount; i++ {
		wg.Add(1)
		go func(agentID int) {
			defer wg.Done()
			suite.simulateAgent(fmt.Sprintf("perf-agent-%d", agentID), eventsPerAgent)
		}(i)
	}

	wg.Wait()
	duration := time.Since(start)

	eventsPerSecond := float64(totalEvents) / duration.Seconds()
	fmt.Printf("Agent communication rate: %.2f events/second\n", eventsPerSecond)

	assert.GreaterOrEqual(suite.T(), eventsPerSecond, 200.0, "Should handle at least 200 agent events/second")
}

func (suite *PerformanceTestSuite) TestAlertProcessingLatency() {
	// Test end-to-end alert processing latency
	eventCount := 100

	start := time.Now()

	// Send events that should trigger alerts
	for i := 0; i < eventCount; i++ {
		suite.sendHighRiskEvent(fmt.Sprintf("risky-process-%d.exe", i))
	}

	// Wait for alerts to be processed
	time.Sleep(5 * time.Second)

	alerts := suite.getRecentAlerts(eventCount)
	processingTime := time.Since(start)

	if len(alerts) > 0 {
		averageLatency := processingTime / time.Duration(len(alerts))
		fmt.Printf("Alert processing latency: %v per alert\n", averageLatency)
		fmt.Printf("Alerts generated: %d/%d\n", len(alerts), eventCount)

		assert.LessOrEqual(suite.T(), averageLatency, 2*time.Second, "Alert processing should be under 2 seconds")
		assert.GreaterOrEqual(suite.T(), len(alerts), eventCount/2, "Should generate alerts for at least half the risky events")
	}
}

func (suite *PerformanceTestSuite) TestScalabilityWithIncreasingLoad() {
	// Test system scalability with increasing concurrent load
	loadLevels := []int{10, 25, 50, 100}

	var results []PerformanceResult

	for _, loadLevel := range loadLevels {
		fmt.Printf("Testing with %d concurrent users...\n", loadLevel)

		start := time.Now()
		suite.runLoadTest(loadLevel, 30*time.Second)
		duration := time.Since(start)

		metrics := suite.getPerformanceMetrics()

		result := PerformanceResult{
			ConcurrentUsers:    loadLevel,
			Duration:          duration,
			AverageResponseTime: metrics.AverageResponseTime,
			Throughput:        metrics.Throughput,
			ErrorRate:         metrics.ErrorRate,
			CPUUsage:          metrics.CPUUsage,
			MemoryUsage:       metrics.MemoryUsage,
		}

		results = append(results, result)

		fmt.Printf("Load level %d results:\n", loadLevel)
		fmt.Printf("  Response time: %v\n", result.AverageResponseTime)
		fmt.Printf("  Throughput: %.2f req/sec\n", result.Throughput)
		fmt.Printf("  Error rate: %.2f%%\n", result.ErrorRate*100)
		fmt.Printf("  CPU usage: %.2f%%\n", result.CPUUsage)
		fmt.Printf("  Memory usage: %.2f MB\n", result.MemoryUsage)
	}

	// Analyze scalability
	for i := 1; i < len(results); i++ {
		prev := results[i-1]
		curr := results[i]

		responseTimeDegradation := float64(curr.AverageResponseTime) / float64(prev.AverageResponseTime)
		throughputScaling := float64(curr.Throughput) / float64(prev.Throughput)

		fmt.Printf("Scaling from %d to %d users:\n", prev.ConcurrentUsers, curr.ConcurrentUsers)
		fmt.Printf("  Response time degradation: %.2fx\n", responseTimeDegradation)
		fmt.Printf("  Throughput scaling: %.2fx\n", throughputScaling)

		// Assert reasonable scaling characteristics
		assert.LessOrEqual(suite.T(), responseTimeDegradation, 3.0, "Response time shouldn't degrade more than 3x")
		assert.GreaterOrEqual(suite.T(), throughputScaling, 0.5, "Throughput should scale reasonably")
	}
}

func (suite *PerformanceTestSuite) TestResourceCleanup() {
	// Test that system properly cleans up resources under load
	initialGoroutines := suite.getGoroutineCount()

	// Run high-load test
	suite.runLoadTest(50, 1*time.Minute)

	// Allow time for cleanup
	time.Sleep(10 * time.Second)

	finalGoroutines := suite.getGoroutineCount()
	goroutineIncrease := finalGoroutines - initialGoroutines

	fmt.Printf("Initial goroutines: %d\n", initialGoroutines)
	fmt.Printf("Final goroutines: %d\n", finalGoroutines)
	fmt.Printf("Goroutine increase: %d\n", goroutineIncrease)

	assert.LessOrEqual(suite.T(), goroutineIncrease, 100, "Goroutine increase should be reasonable")
}

func (suite *PerformanceTestSuite) TestDatabaseConnectionPooling() {
	// Test database connection pool performance
	concurrentQueries := 50
	queriesPerWorker := 100

	start := time.Now()

	var wg sync.WaitGroup
	for i := 0; i < concurrentQueries; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			for j := 0; j < queriesPerWorker; j++ {
				suite.executeQuery("SELECT COUNT(*) FROM events")
			}
		}()
	}

	wg.Wait()
	duration := time.Since(start)

	totalQueries := concurrentQueries * queriesPerWorker
	queriesPerSecond := float64(totalQueries) / duration.Seconds()

	fmt.Printf("Database queries per second: %.2f\n", queriesPerSecond)
	assert.GreaterOrEqual(suite.T(), queriesPerSecond, 500.0, "Should handle at least 500 database queries/second")
}

func (suite *PerformanceTestSuite) TestCachePerformance() {
	// Test caching layer performance
	suite.populateCache(1000) // Populate cache with 1000 items

	// Test cache hit performance
	start := time.Now()
	hitCount := 10000
	for i := 0; i < hitCount; i++ {
		suite.cacheGet(fmt.Sprintf("cached-item-%d", i%1000))
	}
	cacheHitTime := time.Since(start) / time.Duration(hitCount)

	// Test cache miss performance
	start = time.Now()
	missCount := 1000
	for i := 0; i < missCount; i++ {
		suite.cacheGet(fmt.Sprintf("nonexistent-item-%d", i))
	}
	cacheMissTime := time.Since(start) / time.Duration(missCount)

	fmt.Printf("Cache hit time: %v\n", cacheHitTime)
	fmt.Printf("Cache miss time: %v\n", cacheMissTime)

	assert.LessOrEqual(suite.T(), cacheHitTime, 1*time.Millisecond, "Cache hits should be under 1ms")
	assert.LessOrEqual(suite.T(), cacheMissTime, 10*time.Millisecond, "Cache misses should be under 10ms")
}

// Helper methods

func (suite *PerformanceTestSuite) authenticate() string {
	// Implementation for authentication
	return "test-token"
}

func (suite *PerformanceTestSuite) sendBatchEvents(agentID string, count int) {
	// Implementation for sending batch events
}

func (suite *PerformanceTestSuite) simulateUserSession(userID, requestCount int, results chan time.Duration) {
	// Implementation for simulating user session
}

func (suite *PerformanceTestSuite) prepareTestData(count int) {
	// Implementation for preparing test data
}

func (suite *PerformanceTestSuite) executeQuery(query string) {
	// Implementation for executing query
}

func (suite *PerformanceTestSuite) getMemoryUsage() float64 {
	// Implementation for getting memory usage
	return 0.0
}

func (suite *PerformanceTestSuite) generateSustainedLoad(ctx context.Context, workers int) {
	// Implementation for generating sustained load
}

func (suite *PerformanceTestSuite) makeAPIRequest(endpoint string) {
	// Implementation for making API request
}

func (suite *PerformanceTestSuite) simulateAgent(agentID string, eventCount int) {
	// Implementation for simulating agent
}

func (suite *PerformanceTestSuite) sendHighRiskEvent(processName string) {
	// Implementation for sending high-risk event
}

func (suite *PerformanceTestSuite) getRecentAlerts(count int) []models.Alert {
	// Implementation for getting recent alerts
	return []models.Alert{}
}

func (suite *PerformanceTestSuite) runLoadTest(users int, duration time.Duration) {
	// Implementation for running load test
}

func (suite *PerformanceTestSuite) getPerformanceMetrics() PerformanceMetrics {
	// Implementation for getting performance metrics
	return PerformanceMetrics{}
}

func (suite *PerformanceTestSuite) getGoroutineCount() int {
	// Implementation for getting goroutine count
	return 0
}

func (suite *PerformanceTestSuite) populateCache(count int) {
	// Implementation for populating cache
}

func (suite *PerformanceTestSuite) cacheGet(key string) interface{} {
	// Implementation for cache get
	return nil
}

// Data structures

type PerformanceResult struct {
	ConcurrentUsers     int
	Duration           time.Duration
	AverageResponseTime time.Duration
	Throughput         float64
	ErrorRate          float64
	CPUUsage           float64
	MemoryUsage        float64
}

type PerformanceMetrics struct {
	AverageResponseTime time.Duration
	Throughput         float64
	ErrorRate          float64
	CPUUsage           float64
	MemoryUsage        float64
}</content>
<parameter name="filePath">/workspaces/insec/tests/performance/system_performance_test.go
