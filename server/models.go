package main

import (
	"time"

	"gorm.io/gorm"
)

type User struct {
	ID        string `json:"id" gorm:"primaryKey"`
	Email     string `json:"email" gorm:"uniqueIndex"`
	Dept      string `json:"dept"`
	RiskScore *int   `json:"risk_score,omitempty"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type OS struct {
	Family  string `json:"family"`
	Version string `json:"version"`
	Arch    string `json:"arch"`
}

type EventData struct {
	Type     string `json:"type"`
	ID       string `json:"id"`
	Category string `json:"category"`
}

type ProcessInfo struct {
	Name       string   `json:"name"`
	Cmd        []string `json:"cmd"`
	PPID       uint32   `json:"ppid"`
	PID        uint32   `json:"pid"`
	Hash       *string  `json:"hash,omitempty"`
	User       *string  `json:"user,omitempty"`
	CPUUsage   float32  `json:"cpu_usage"`
	MemoryUsage uint64  `json:"memory_usage"`
}

type NetworkInfo struct {
	SrcIP    *string `json:"src_ip,omitempty"`
	DstIP    string  `json:"dst_ip"`
	DstPort  uint16  `json:"dst_port"`
	Protocol string  `json:"protocol"`
	BytesIn  uint64  `json:"bytes_in"`
	BytesOut uint64  `json:"bytes_out"`
	Domain   *string `json:"domain,omitempty"`
}

type FileInfo struct {
	Path      string  `json:"path"`
	Operation string  `json:"operation"`
	Size      *uint64 `json:"size,omitempty"`
	Hash      *string `json:"hash,omitempty"`
	User      *string `json:"user,omitempty"`
}

type Agent struct {
	Version  string `json:"ver"`
	Mode     string `json:"mode"`
	Hostname string `json:"hostname"`
}

type Event struct {
	ID        uint      `json:"-" gorm:"primaryKey"`
	Timestamp time.Time `json:"ts" gorm:"index"`
	TenantID  string    `json:"tenant_id" gorm:"index"`
	HostID    string    `json:"host_id" gorm:"index"`
	User      User      `json:"user" gorm:"embedded;embeddedPrefix:user_"`
	OS        OS        `json:"os" gorm:"embedded;embeddedPrefix:os_"`
	Event     EventData `json:"event" gorm:"embedded;embeddedPrefix:event_"`
	Process   *ProcessInfo `json:"proc,omitempty" gorm:"embedded;embeddedPrefix:proc_"`
	Network   *NetworkInfo `json:"net,omitempty" gorm:"embedded;embeddedPrefix:net_"`
	File      *FileInfo    `json:"file,omitempty" gorm:"embedded;embeddedPrefix:file_"`
	Labels    []string `json:"labels" gorm:"type:text[]"`
	RiskHints []string `json:"risk_hints" gorm:"type:text[]"`
	Agent     Agent    `json:"agent" gorm:"embedded;embeddedPrefix:agent_"`
	SessionID string   `json:"session_id" gorm:"index"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type Alert struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	CreatedAt   time.Time `json:"created"`
	Severity    string    `json:"severity"`
	Title       string    `json:"title"`
	Description string    `json:"description,omitempty"`
	TenantID    string    `json:"tenant_id" gorm:"index"`
	RuleID      string    `json:"rule_id"`
	RuleVersion string    `json:"rule_version"`
	UEBAScore   int       `json:"ueba_score"`
	Entities    map[string]interface{} `json:"entities" gorm:"type:jsonb"`
	Evidence    []string  `json:"evidence" gorm:"type:text[]"`
	Status      string    `json:"status" gorm:"default:'open'"`
	Assignee    *string   `json:"assignee,omitempty"`
	ResolvedAt  *time.Time `json:"resolved_at,omitempty"`
}

type Rule struct {
	ID          string                 `json:"id" gorm:"primaryKey"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	Version     string                 `json:"version"`
	Enabled     bool                   `json:"enabled" gorm:"default:true"`
	Severity    string                 `json:"severity"`
	Conditions  map[string]interface{} `json:"conditions" gorm:"type:jsonb"`
	Actions     []string               `json:"actions" gorm:"type:text[]"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

type Policy struct {
	ID          string                 `json:"id" gorm:"primaryKey"`
	Name        string                 `json:"name"`
	Description string                 `json:"description"`
	TenantID    string                 `json:"tenant_id" gorm:"index"`
	Enabled     bool                   `json:"enabled" gorm:"default:true"`
	Config      map[string]interface{} `json:"config" gorm:"type:jsonb"`
	CreatedAt   time.Time
	UpdatedAt   time.Time
}

// Database migration
func AutoMigrate(db *gorm.DB) error {
	return db.AutoMigrate(
		&Event{},
		&Alert{},
		&Rule{},
		&Policy{},
	)
}
