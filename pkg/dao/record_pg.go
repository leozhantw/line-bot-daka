package dao

import (
	"time"

	"gorm.io/gorm"
)

type PGRecordDAO struct {
	db *gorm.DB
}

func NewPGRecordDAO(db *gorm.DB) RecordDAO {
	return PGRecordDAO{
		db: db,
	}
}

func (o PGRecordDAO) GetByUser(userID string, date time.Time) (*Record, error) {
	var r Record
	err := o.db.
		Where("user_id = ?", userID).
		Where("work_date = ?", date.Format("2006-01-02")).
		First(&r).
		Error

	return &r, err
}

func (o PGRecordDAO) ListByDate(date time.Time) (*[]Record, error) {
	var rs []Record
	err := o.db.Where("work_date = ?", date.Format("2006-01-02")).Find(&rs).Error

	return &rs, err
}

func (o PGRecordDAO) Create(record *Record) error {
	return o.db.Create(record).Error
}
