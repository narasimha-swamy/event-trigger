package handlers

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/narasimha-swamy/event-trigger/models"
	"gorm.io/gorm"
)


// get events using get request
func GetEvents(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var events []models.EventLog
        query := db.Model(&models.EventLog{})

        if status := c.Query("status"); status != "" {
            query = query.Where("status = ?", status)
        }
        if triggerID := c.Query("trigger_id"); triggerID != "" {
            query = query.Where("trigger_id = ?", triggerID)
        }

        // Sort by triggered_at in descending order
        query = query.Order("triggered_at DESC")

        if err := query.Find(&events).Error; err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to fetch events"})
            return
        }

        c.JSON(http.StatusOK, events)
    }
}

// fire the event with id 
// Note: not for testing purpose but for actually firing it
func FireAPITrigger(db *gorm.DB) gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")
		var trigger models.Trigger
		if err := db.First(&trigger, "id = ? AND type = ?", id, models.APITrigger).Error; err != nil {
			c.JSON(http.StatusNotFound, gin.H{"error": "API trigger not found"})
			return
		}

		var payload map[string]string
		if err := c.ShouldBindJSON(&payload); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "invalid payload"})
			return
		}

		// Create event log
		event := models.EventLog{
			TriggerID:   trigger.ID,
			TriggeredAt: time.Now().In(time.FixedZone("IST", 5*60*60+30*60)),
			Payload:     payload,
			IsTest:      true,
			Status:      "active",
		}

		if err := db.Create(&event).Error; err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to log event"})
			return
		}

		c.JSON(http.StatusOK, event)
	}
}
