package models

import "time"

type Meter struct {
	MeterID   string    `json:"meter_id" gorm:"column:meter_id;primaryKey"`
	MeterName string    `json:"meter_name" gorm:"column:meter_name"`
	Location  string    `json:"location" gorm:"column:location"`
	CreatedAt time.Time `json:"created_at" gorm:"column:created_at;autoCreateTime"`
}

func (Meter) TableName() string { return "meters" }

