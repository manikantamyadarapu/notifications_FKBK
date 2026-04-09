package repository

import (
	"goproject/internal/models"

	"gorm.io/gorm"
)

type EventRepository struct {
	DB *gorm.DB
}

func (r *EventRepository) MeterExists(meterID string) (bool, error) {
	var m models.Meter
	err := r.DB.Select("meter_id").First(&m, "meter_id = ?", meterID).Error
	if err == gorm.ErrRecordNotFound {
		return false, nil
	}
	return err == nil, err
}

func (r *EventRepository) SaveEvent(e models.TamperEvent) (uint, error) {
	e.IsProcessed = false
	if err := r.DB.Create(&e).Error; err != nil {
		return 0, err
	}
	return e.ID, nil
}

func (r *EventRepository) GetUnprocessedEvents() ([]models.TamperEvent, error) {
	var events []models.TamperEvent
	err := r.DB.
		Where("processed = ?", false).
		Order("timestamp asc").
		Find(&events).Error
	return events, err
}

func (r *EventRepository) MarkProcessed(id uint) {
	r.DB.Model(&models.TamperEvent{}).Where("id = ?", id).Update("processed", true)
}

func (r *EventRepository) GetTamperDescription(code int) string {
	// Mandatory lookup. We intentionally do NOT hardcode descriptions.
	//
	// Your DB schema might use either snake_case (tamper_code/tamper_desc)
	// or camelCase (tamperCode/tamperDesc). We try both to stay compatible.
	var desc string

	// Attempt 1: snake_case
	if err := r.DB.
		Raw(`SELECT tamper_desc FROM tamper_code_desc WHERE tamper_code = ? LIMIT 1`, code).
		Scan(&desc).Error; err == nil && desc != "" {
		return desc
	}

	// Attempt 2: camelCase quoted identifiers (Postgres preserves case when quoted)
	desc = ""
	if err := r.DB.
		Raw(`SELECT "tamperDesc" FROM tamper_code_desc WHERE "tamperCode" = ? LIMIT 1`, code).
		Scan(&desc).Error; err == nil && desc != "" {
		return desc
	}

	return "Unknown Tamper Code"
}