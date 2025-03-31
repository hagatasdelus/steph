package model

import (
	"time"
)

type StepHistory struct {
	ID 			uint 		`gorm:"primaryKey;autoIncrement"`
	Date		time.Time	`gorm:"uniqueIndex;not null;type:date"`
	Steps		uint		`gorm:"not null"`
	CreatedAt	time.Time	`gorm:"autoCreateTime"`
	UpdatedAt	time.Time	`gorm:"autoUpdateTime"`
}
