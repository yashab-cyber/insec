package main

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Event struct {
	Ts       string `json:"ts"`
	TenantID string `json:"tenant_id"`
	HostID   string `json:"host_id"`
	User     struct {
		ID    string `json:"id"`
		Email string `json:"email"`
		Dept  string `json:"dept"`
	} `json:"user"`
	Os struct {
		Family  string `json:"family"`
		Version string `json:"version"`
	} `json:"os"`
	Event struct {
		Type string `json:"type"`
		ID   string `json:"id"`
	} `json:"event"`
	Proc struct {
		Name string   `json:"name"`
		Cmd  []string `json:"cmd"`
		Ppid uint32   `json:"ppid"`
		Hash *string  `json:"hash,omitempty"`
	} `json:"proc"`
	Net     interface{} `json:"net,omitempty"`
	File    interface{} `json:"file,omitempty"`
	Labels  []string    `json:"labels"`
	RiskHints []string  `json:"risk_hints"`
	Agent   struct {
		Ver  string `json:"ver"`
		Mode string `json:"mode"`
	} `json:"agent"`
}

func main() {
	r := gin.Default()

	// CORS middleware
	r.Use(func(c *gin.Context) {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		c.Header("Access-Control-Allow-Headers", "Origin, Content-Type, Authorization")
		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}
		c.Next()
	})

	r.POST("/v1/events", func(c *gin.Context) {
		var events []Event
		if err := c.ShouldBindJSON(&events); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		for _, event := range events {
			log.Printf("Received event: %+v", event)
		}

		c.JSON(http.StatusOK, gin.H{"status": "ok", "count": len(events)})
	})

	r.GET("/v1/events", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"events": []interface{}{}})
	})

	r.GET("/v1/alerts", func(c *gin.Context) {
		c.JSON(http.StatusOK, gin.H{"alerts": []interface{}{}})
	})

	r.GET("/health", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"status": "healthy",
			"version": "1.0.0",
		})
	})

	fmt.Println("INSEC Server starting on :8080")
	r.Run(":8080")
}
