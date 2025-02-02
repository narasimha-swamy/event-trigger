package handlers

import (
	"errors"
	"fmt"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/narasimha-swamy/event-trigger/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

type TestScheduledRequest struct {
	DelayMinutes string `json:"delay_minutes"`
}

// create trigger based on the type (either API or scheduler)
func CreateTrigger(db *gorm.DB) gin.HandlerFunc {
	fmt.Println("came here successfully")
	return func(c *gin.Context) {
		var trigger models.Trigger
		if err := c.ShouldBindJSON(&trigger); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		if trigger.Type == models.ScheduledTrigger {
			if _, err := cron.ParseStandard(trigger.CronExpression); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid cron expression"})
				return
			}
			schedule, _ := cron.ParseStandard(trigger.CronExpression)
			trigger.NextRun = schedule.Next(time.Now().In(time.FixedZone("IST", 5*60*60+30*60)))
		} else if trigger.Type == models.APITrigger {
			trigger.Endpoint = "/api/triggers/fire/" + uuid.New().String()
		}

		if err := db.Create(&trigger).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to create trigger"})
			return
		}

		c.JSON(http.StatusCreated, trigger)
	}
}

// get all the triggers currently available
func GetTriggers(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		var triggers []models.Trigger
		result := db.Find(&triggers)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch triggers"})
			return
		}
		c.JSON(http.StatusOK, triggers)
	}
}

// get the specific trigger using trigger id
func GetTrigger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var trigger models.Trigger
		result := db.First(&trigger, "id = ?", id)
		if result.Error != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trigger not found"})
			return
		}
		c.JSON(http.StatusOK, trigger)
	}
}

// update the trigger data accordingly
func UpdateTrigger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var existingTrigger models.Trigger
		if err := db.First(&existingTrigger, "id = ?", id).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "trigger not found"})
			return
		}

		var updatedTrigger models.Trigger
		if err := c.ShouldBindJSON(&updatedTrigger); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
			return
		}

		// Update fields
		existingTrigger.Type = updatedTrigger.Type
		existingTrigger.CronExpression = updatedTrigger.CronExpression
		existingTrigger.IsRecurring = updatedTrigger.IsRecurring
		existingTrigger.APIPayload = updatedTrigger.APIPayload
		existingTrigger.IsActive = updatedTrigger.IsActive
		existingTrigger.UpdatedAt = time.Now().In(time.FixedZone("IST", 5*60*60+30*60))

		if err := db.Save(&existingTrigger).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update trigger"})
			return
		}

		c.JSON(http.StatusOK, existingTrigger)
	}
}

// Delete the specific trigger
func DeleteTrigger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		result := db.Delete(&models.Trigger{}, "id = ?", id)
		if result.Error != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete trigger"})
			return
		}
		if result.RowsAffected == 0 {
			c.JSON(http.StatusNotFound, gin.H{"error": "trigger not found"})
			return
		}
		c.Status(http.StatusNoContent)
	}
}


func TestTrigger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		triggerID := c.Param("id")
		fmt.Println("trigger id", triggerID)
		var trigger models.Trigger

		// Retrieve trigger with proper error handling
		if err := db.Where("id = ?", triggerID).First(&trigger).Error; err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "trigger not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": "database error"})
			return
		}

		var payload map[string]string
		var event models.EventLog

		switch trigger.Type {
		case models.APITrigger:
			fmt.Println("is a api trigger", triggerID)
			// Validate and parse payload for API triggers
			if err := c.ShouldBindJSON(&payload); err != nil {
				c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload format"})
				return
			}

			// Validate payload schema with placeholder support
			for key, expectedValue := range trigger.APIPayload {
				receivedValue, exists := payload[key]
				if !exists {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": fmt.Sprintf("missing required key: %s", key),
					})
					return
				}

				// Allow any value if schema contains "<value>" placeholder
				if expectedValue != "<value>" && receivedValue != expectedValue {
					c.JSON(http.StatusBadRequest, gin.H{
						"error": fmt.Sprintf("value mismatch for key '%s'", key),
					})
					return
				}
			}

			event = models.EventLog{
				TriggerID:   trigger.ID,
				TriggeredAt: time.Now().UTC(),
				Payload:     payload,
				IsTest:      true,
				Status:      "active",
			}
		
			// Create event log with error handling
			if err := db.Create(&event).Error; err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error":   "failed to log test event",
					"details": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, gin.H{
				"message":  "Trigger tested successfully",
				"event_id": event.ID,
				"is_test":  true,
			})

		case models.ScheduledTrigger:
			fmt.Println("is a scheduler id trigger", triggerID)
			var request TestScheduledRequest
			if err := c.ShouldBindJSON(&request); err != nil {
                c.JSON(http.StatusBadRequest, gin.H{"error": "invalid request"})
                return
            }

			// Schedule test execution 
			// Note: this fuction will trigger the event after certain minutes provided in the input 
			// format {"dealy_minutes": "10"} after 10 minutes
			go func(delay time.Duration) {
				time.Sleep(delay)
				event := models.EventLog{
					TriggerID:   trigger.ID,
					TriggeredAt: time.Now().UTC(),
					IsTest:      true,
					Status:      "active",
				}
				db.Create(&event)
			}(time.Duration(1) * time.Minute)

			c.JSON(http.StatusOK, gin.H{
				"message": fmt.Sprintf("Test scheduled to run in %s minutes", request.DelayMinutes),
			})
			return

		default:
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid trigger type"})
			return
		}

	}
}
