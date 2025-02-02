package scheduler

import (
	"time"

	"github.com/narasimha-swamy/event-trigger/models"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
)

func StartScheduler(db *gorm.DB) {

	// have implemented the below logic to ensure that it runs exactly at every 1 min and not like 12:30:24
	// Calculate time until next minute
    now := time.Now().UTC()
    nextMinute := now.Truncate(time.Minute).Add(time.Minute)
    initialDelay := nextMinute.Sub(now)

    // Wait until the next minute starts
    time.Sleep(initialDelay)

	// Start the ticker
	ticker := time.NewTicker(1 * time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		var triggers []models.Trigger
		now := time.Now().In(time.FixedZone("IST", 5*60*60+30*60))
		db.Where("type = ? AND is_active = ? AND next_run <= ?",
			models.ScheduledTrigger, true, now).Find(&triggers)

		for _, trigger := range triggers {
			event := models.EventLog{
				TriggerID:   trigger.ID,
				TriggeredAt: now,
				Status:      "active",
			}
			db.Create(&event)

			// Update next run or deactivate
			if trigger.IsRecurring {
				schedule, _ := cron.ParseStandard(trigger.CronExpression)
				trigger.NextRun = schedule.Next(now)
				db.Save(&trigger)
			} else {
				db.Model(&trigger).Update("is_active", false)
			}
		}
	}
}
