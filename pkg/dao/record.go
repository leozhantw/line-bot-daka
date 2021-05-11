package dao

import (
	"time"
)

type Record struct {
	ID        uint `gorm:"primarykey"`
	UserID    string
	WorkDate  time.Time
	WorkedAt  time.Time `gorm:"autoCreateTime"`
	CreatedAt time.Time
	UpdatedAt time.Time
}

type RecordDAO interface {
	GetByUser(userID string, date time.Time) (*Record, error)
	ListByDate(date time.Time) (*[]Record, error)
	Create(record *Record) error
}
