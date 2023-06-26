package tasks

import (
	"SMS-panel/models"
	"log"
	"time"

	"gorm.io/gorm"
	"gorm.io/gorm/clause"
)

func RentNumberTask(db *gorm.DB) TaskFunc {
	return func() {
		now := time.Now()
		log.Println("Start ------------------------", now)
		var userNumberObjects []models.UserNumbers

		tx := db.Begin()

		err := tx.Clauses(clause.Returning{Columns: []clause.Column{{Name: "number_id"}}}).
			Where("user_numbers.end_date < ?", now).
			Delete(&userNumberObjects).Error
		log.Println(userNumberObjects, "=====")
		if err != nil || len(userNumberObjects) == 0 {
			log.Println("There is no userNumber to delete")
			return
		}

		var userNumbersIDs []int
		for _, userNumber := range userNumberObjects {
			userNumbersIDs = append(userNumbersIDs, int(userNumber.NumberID))
		}

		err = tx.Table("sender_numbers").
			Where("id IN ?", userNumbersIDs).
			Updates(map[string]interface{}{"is_exclusive": false}).Error
		if err != nil {
			log.Println("Can't update senderNumbers")
			return
		}

		tx.Commit()
	}
}
