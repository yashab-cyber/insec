package main

import (
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/sirupsen/logrus"
	"gorm.io/gorm"
)

type EventHandler struct {
	db     *gorm.DB
	logger *logrus.Logger
}

func NewEventHandler(db *gorm.DB, logger *logrus.Logger) *EventHandler {
	return &EventHandler{
		db:     db,
		logger: logger,
	}
}

// POST /v1/events - Ingest events from agents
func (h *EventHandler) IngestEvents(c *gin.Context) {
	var events []Event
	if err := c.ShouldBindJSON(&events); err != nil {
		h.logger.WithError(err).Error("Failed to bind events JSON")
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON format"})
		return
	}

	// Process events in batch
	for i := range events {
		events[i].ID = 0 // Let GORM auto-generate
		events[i].CreatedAt = time.Now()
		events[i].UpdatedAt = time.Now()
	}

	// Bulk insert
	if err := h.db.CreateInBatch(&events, 100).Error; err != nil {
		h.logger.WithError(err).Error("Failed to insert events")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to process events"})
		return
	}

	h.logger.WithField("count", len(events)).Info("Successfully ingested events")

	// Process events for rules and alerts
	go h.processEventsForAlerts(events)

	c.JSON(http.StatusOK, gin.H{
		"status": "ok",
		"count":  len(events),
	})
}

// GET /v1/events - Query events with filtering
func (h *EventHandler) GetEvents(c *gin.Context) {
	var events []Event
	query := h.db

	// Apply filters
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if hostID := c.Query("host_id"); hostID != "" {
		query = query.Where("host_id = ?", hostID)
	}
	if eventType := c.Query("event_type"); eventType != "" {
		query = query.Where("event_type = ?", eventType)
	}
	if startTime := c.Query("start_time"); startTime != "" {
		if t, err := time.Parse(time.RFC3339, startTime); err == nil {
			query = query.Where("timestamp >= ?", t)
		}
	}
	if endTime := c.Query("end_time"); endTime != "" {
		if t, err := time.Parse(time.RFC3339, endTime); err == nil {
			query = query.Where("timestamp <= ?", t)
		}
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Find(&events).Error; err != nil {
		h.logger.WithError(err).Error("Failed to query events")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query events"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"events": events,
		"page":   page,
		"limit":  limit,
	})
}

// GET /v1/alerts - Get alerts
func (h *EventHandler) GetAlerts(c *gin.Context) {
	var alerts []Alert
	query := h.db

	// Apply filters
	if tenantID := c.Query("tenant_id"); tenantID != "" {
		query = query.Where("tenant_id = ?", tenantID)
	}
	if status := c.Query("status"); status != "" {
		query = query.Where("status = ?", status)
	}
	if severity := c.Query("severity"); severity != "" {
		query = query.Where("severity = ?", severity)
	}

	// Pagination
	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "50"))
	offset := (page - 1) * limit

	if err := query.Offset(offset).Limit(limit).Order("created_at DESC").Find(&alerts).Error; err != nil {
		h.logger.WithError(err).Error("Failed to query alerts")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query alerts"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"alerts": alerts,
		"page":   page,
		"limit":  limit,
	})
}

// POST /v1/alerts/:id/actions - Execute alert actions
func (h *EventHandler) ExecuteAlertAction(c *gin.Context) {
	alertID := c.Param("id")
	var actionReq struct {
		Action string                 `json:"action" binding:"required"`
		Params map[string]interface{} `json:"params,omitempty"`
	}

	if err := c.ShouldBindJSON(&actionReq); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Find alert
	var alert Alert
	if err := h.db.First(&alert, alertID).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Alert not found"})
		return
	}

	// Execute action (simplified)
	h.logger.WithFields(logrus.Fields{
		"alert_id": alert.ID,
		"action":   actionReq.Action,
	}).Info("Executing alert action")

	// In a real implementation, this would trigger actual response actions
	// like isolating hosts, revoking tokens, creating tickets, etc.

	c.JSON(http.StatusOK, gin.H{
		"status": "action_executed",
		"alert_id": alert.ID,
		"action": actionReq.Action,
	})
}

// GET /v1/rules - Get detection rules
func (h *EventHandler) GetRules(c *gin.Context) {
	var rules []Rule
	if err := h.db.Find(&rules).Error; err != nil {
		h.logger.WithError(err).Error("Failed to query rules")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to query rules"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"rules": rules})
}

// POST /v1/rules - Create new rule
func (h *EventHandler) CreateRule(c *gin.Context) {
	var rule Rule
	if err := c.ShouldBindJSON(&rule); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid rule format"})
		return
	}

	rule.ID = uuid.New().String()
	rule.CreatedAt = time.Now()
	rule.UpdatedAt = time.Now()

	if err := h.db.Create(&rule).Error; err != nil {
		h.logger.WithError(err).Error("Failed to create rule")
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create rule"})
		return
	}

	h.logger.WithField("rule_id", rule.ID).Info("Created new rule")
	c.JSON(http.StatusCreated, rule)
}

// Process events for rule matching and alert generation
func (h *EventHandler) processEventsForAlerts(events []Event) {
	// Simplified rule processing - in production, this would be more sophisticated
	for _, event := range events {
		// Check for suspicious process patterns
		if event.Process != nil {
			suspiciousCmds := []string{"netcat", "ncat", "wget", "curl", "scp", "rclone"}
			for _, cmd := range event.Process.Cmd {
				for _, suspicious := range suspiciousCmds {
					if strings.Contains(cmd, suspicious) {
						h.createAlert(&event, "Suspicious data transfer tool detected", "high")
						break
					}
				}
			}
		}

		// Check for large file operations
		if event.File != nil && event.File.Size != nil && *event.File.Size > 100*1024*1024 { // 100MB
			h.createAlert(&event, "Large file operation detected", "medium")
		}
	}
}

func (h *EventHandler) createAlert(event *Event, title, severity string) {
	alert := Alert{
		CreatedAt:   time.Now(),
		Severity:    severity,
		Title:       title,
		TenantID:    event.TenantID,
		RuleID:      "RULE_AUTO",
		RuleVersion: "1.0",
		UEBAScore:   75,
		Entities: map[string]interface{}{
			"user": event.User.ID,
			"host": event.HostID,
		},
		Evidence: []string{event.Event.ID},
		Status:   "open",
	}

	if err := h.db.Create(&alert).Error; err != nil {
		h.logger.WithError(err).Error("Failed to create alert")
	} else {
		h.logger.WithField("alert_id", alert.ID).Info("Created alert")
	}
}
