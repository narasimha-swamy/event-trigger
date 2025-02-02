package models

import (
	"time"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

type EventLog struct {
	ID          uuid.UUID         `gorm:"type:uuid;primaryKey" json:"id"`
	TriggerID   uuid.UUID         `gorm:"index" json:"trigger_id"`
	TriggeredAt time.Time         `gorm:"index" json:"triggered_at"`
	Payload     map[string]string `gorm:"serializer:json" json:"payload,omitempty"`
	IsTest      bool              `json:"is_test"`
	Status      string            `gorm:"index;default:'active'" json:"status"`
	CreatedAt   time.Time         `json:"created_at"`
}

func (e *EventLog) BeforeCreate(tx *gorm.DB) error {
	e.ID = uuid.New()
	e.CreatedAt = time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
	return nil
}
