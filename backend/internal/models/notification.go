package models

import "time"

type Notification struct {
	ID                uint      `json:"id" gorm:"primaryKey"`
	MeterID           string    `json:"meter_id" gorm:"column:meter_id;index:idx_escalation_notifications_meter_time,priority:1"`
	TamperCode        int       `json:"tamper_code" gorm:"column:tamper_code;index:idx_escalation_notifications_code_time,priority:1"`
	Message           string    `json:"message" gorm:"column:message"`
	Timestamp         time.Time `json:"timestamp" gorm:"column:timestamp"`

	// Columns present in your existing `escalation_notifications` table (per screenshots).
	Type    string `json:"type" gorm:"column:type"`
	Level   int    `json:"level" gorm:"column:level"`
	Status  string `json:"status" gorm:"column:status"`

	ScheduledFor *time.Time `json:"scheduled_for" gorm:"column:scheduled_for"`
	SentAt        *time.Time `json:"sent_at" gorm:"column:sent_at"`
	ResolvedAt    *time.Time `json:"resolved_at" gorm:"column:resolved_at"`

	// Not stored in your table (we compute from `tamper_code_desc` and return in API).
	TamperDescription string `json:"tamper_description" gorm:"-"`
}

func (Notification) TableName() string { return "escalation_notifications" }