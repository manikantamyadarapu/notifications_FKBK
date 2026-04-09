package repository

import (
	"goproject/internal/models"
	"time"

	"gorm.io/gorm"
)

type NotificationRepository struct {
	DB *gorm.DB
}

func (r *NotificationRepository) SaveNotification(n models.Notification) error {
	return r.DB.Create(&n).Error
}

func (r *NotificationRepository) SaveNotificationWithResult(n *models.Notification) error {
	return r.DB.Create(n).Error
}

func (r *NotificationRepository) GetLatestPending(meterID string, tamperCode int) (models.Notification, bool, error) {
	var n models.Notification
	err := r.DB.
		Where("meter_id = ? AND tamper_code = ? AND status = ?", meterID, tamperCode, "pending").
		Order("timestamp desc").
		First(&n).Error
	if err == gorm.ErrRecordNotFound {
		return models.Notification{}, false, nil
	}
	if err != nil {
		return models.Notification{}, false, err
	}
	return n, true, nil
}

func (r *NotificationRepository) UpdateNotification(n *models.Notification) error {
	return r.DB.Save(n).Error
}

func (r *NotificationRepository) GetAll() ([]models.Notification, error) {
	var list []models.Notification
	err := r.DB.Order("timestamp desc").Find(&list).Error
	return list, err
}

type NotificationFilter struct {
	MeterID    string
	TamperCode *int
	Type       string
	From       *time.Time
	To         *time.Time
	Limit      int
	Offset     int
}

func (r *NotificationRepository) buildFilterQuery(f NotificationFilter) *gorm.DB {
	tx := r.DB.Model(&models.Notification{})
	if f.MeterID != "" {
		tx = tx.Where("meter_id = ?", f.MeterID)
	}
	if f.TamperCode != nil {
		tx = tx.Where("tamper_code = ?", *f.TamperCode)
	}
	if f.Type != "" {
		tx = tx.Where("type = ?", f.Type)
	}
	if f.From != nil {
		tx = tx.Where("timestamp >= ?", *f.From)
	}
	if f.To != nil {
		tx = tx.Where("timestamp <= ?", *f.To)
	}
	return tx
}

func (r *NotificationRepository) CountFiltered(f NotificationFilter) (int64, error) {
	var total int64
	err := r.buildFilterQuery(f).Count(&total).Error
	return total, err
}

func (r *NotificationRepository) GetFiltered(f NotificationFilter) ([]models.Notification, error) {
	limit := f.Limit
	if limit <= 0 || limit > 500 {
		limit = 100
	}
	offset := f.Offset
	if offset < 0 {
		offset = 0
	}

	var list []models.Notification
	tx := r.buildFilterQuery(f)

	err := tx.Order("timestamp desc").Limit(limit).Offset(offset).Find(&list).Error
	if err == gorm.ErrRecordNotFound {
		return []models.Notification{}, nil
	}
	return list, err
}