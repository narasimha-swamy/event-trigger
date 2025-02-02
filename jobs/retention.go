package jobs

import (
	"time"

	"github.com/narasimha-swamy/event-trigger/models"
	"gorm.io/gorm"
)

func StartRetentionJob(db *gorm.DB) {
	ticker := time.NewTicker(1 * time.Hour)
	defer ticker.Stop()

	// the status of the application will be denoted by status columm (archieved or active)

	for range ticker.C {
		// Archive events older than 2 hours
		db.Model(&models.EventLog{}).
			Where("status = 'active' AND triggered_at < ?",
				time.Now().In(time.FixedZone("IST", 5*60*60+30*60)).Add(-2*time.Hour)).
			Update("status", "archived")

		// Delete events older than 48 hours
		db.Where("status = 'archived' AND triggered_at < ?",
			time.Now().In(time.FixedZone("IST", 5*60*60+30*60)).Add(-48*time.Hour)).
			Delete(&models.EventLog{})
	}
}
