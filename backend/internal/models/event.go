package models

import "time"

type TamperEvent struct {
	ID          uint      `json:"id" gorm:"primaryKey"`
	MeterID     string    `json:"meter_id" gorm:"column:meter_id;index:idx_tamper_events_meter_time,priority:1"`
	TamperCode  int       `json:"tamper_code" gorm:"column:tamper_code"`
	Timestamp   time.Time `json:"timestamp" gorm:"column:timestamp;index:idx_tamper_events_meter_time,priority:2"`
	IsProcessed bool      `json:"is_processed" gorm:"column:processed;index"`
}

func (TamperEvent) TableName() string { return "tamper_events" }