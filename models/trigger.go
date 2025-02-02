package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type TriggerType string

const (
	ScheduledTrigger TriggerType = "scheduled"
	APITrigger       TriggerType = "api"
)

type Trigger struct {
	ID             uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	Type           TriggerType       `gorm:"not null" json:"type"`
	CronExpression string            `json:"cron_expression,omitempty"`
	NextRun        time.Time         `json:"next_run,omitempty"`
	IsRecurring    bool              `json:"is_recurring"`
	APIPayload     map[string]string `gorm:"serializer:json" json:"api_payload,omitempty"`
	Endpoint       string            `json:"endpoint,omitempty"`
	IsActive       bool              `gorm:"default:true" json:"is_active"`
	CreatedAt      time.Time         `json:"created_at"`
	UpdatedAt      time.Time         `json:"updated_at"`
	DeletedAt      gorm.DeletedAt    `gorm:"index" json:"-"`
	HitCount       int               `gorm:"default:0" json:"hit_count"`
}

func (t *Trigger) BeforeCreate(tx *gorm.DB) error {
	t.ID = uuid.New()
	t.CreatedAt = time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	t.UpdatedAt = time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	return nil
}
