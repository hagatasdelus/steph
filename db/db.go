package db

import (
	"log"
	"os"
	"path/filepath"
	"time"

	"github.com/hagatasdelus/steph/model"
	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

var (
	db *gorm.DB
)

func InitDB() error {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		return err
	}

	appDir := filepath.Join(homeDir, ".steph")
	if err := os.MkdirAll(appDir, 0755); err != nil {
		return err
	}

	dbPath := filepath.Join(appDir, "steph.db")

	newLogger := logger.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags),
		logger.Config{
			SlowThreshold:             200 * time.Millisecond,
			LogLevel:                  logger.Warn,
			IgnoreRecordNotFoundError: true,
			Colorful:                  true,
		},
	)

	db, err = gorm.Open(sqlite.Open(dbPath), &gorm.Config{
		Logger: newLogger,
	})
	if err != nil {
		return err
	}

	if err := db.AutoMigrate(&model.StepHistory{}); err != nil {
		return err
	}

	return nil
}

func GetDB() *gorm.DB {
	return db
}

func SetSteps(date time.Time, steps uint) error {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var history model.StepHistory
	result := db.Where("date = ?", date).First(&history)

	if result.Error == nil {
		history.Steps = steps
		return db.Save(&history).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		newHistory := model.StepHistory{
			Date:  date,
			Steps: steps,
		}
		return db.Create(&newHistory).Error
	} else {
		return result.Error
	}
}

func AddSteps(date time.Time, steps uint) error {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var history model.StepHistory
	result := db.Where("date = ?", date).First(&history)

	if result.Error == nil {
		history.Steps += steps
		return db.Save(&history).Error
	} else if result.Error == gorm.ErrRecordNotFound {
		newHistory := model.StepHistory{
			Date:  date,
			Steps: steps,
		}
		return db.Create(&newHistory).Error
	} else {
		return result.Error
	}
}

func GetStepsForDate(date time.Time) (uint, error) {
	date = time.Date(date.Year(), date.Month(), date.Day(), 0, 0, 0, 0, date.Location())

	var history model.StepHistory
	result := db.Where("date = ?", date).First(&history)

	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return 0, nil
		}
		return 0, result.Error
	}

	return history.Steps, nil
}

func GetAllStepHistories() ([]model.StepHistory, error) {
	var histories []model.StepHistory
	result := db.Order("date").Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	return histories, nil
}

func GetStepHistoriesByDateRange(startDate, endDate time.Time) ([]model.StepHistory, error) {
	var histories []model.StepHistory
	result := db.Where("date BETWEEN ? AND ?", startDate, endDate).Order("date").Find(&histories)
	if result.Error != nil {
		return nil, result.Error
	}
	return histories, nil
}

