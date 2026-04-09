package service

import (
	"errors"
	"fmt"
	"goproject/internal/models"
	"goproject/internal/repository"
	"goproject/internal/utils"
	"math"
	"time"
)

var ErrUnknownMeterID = errors.New("unknown meter_id")

type NotificationService struct {
	EventRepo *repository.EventRepository
	NotifRepo *repository.NotificationRepository
	WS        *WSManager
}

func (s *NotificationService) IngestEvent(e models.TamperEvent) (uint, models.Notification, error) {
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	eventID, err := s.EventRepo.SaveEvent(e)
	if err != nil {
		return 0, models.Notification{}, err
	}
	e.ID = eventID

	notification, err := s.ProcessEvent(e)
	if err != nil {
		return 0, models.Notification{}, err
	}

	if eventID != 0 {
		s.EventRepo.MarkProcessed(eventID)
	}
	return eventID, notification, nil
}

func (s *NotificationService) ProcessEvent(e models.TamperEvent) (models.Notification, error) {
	if s.EventRepo != nil {
		ok, err := s.EventRepo.MeterExists(e.MeterID)
		if err == nil && !ok {
			return models.Notification{}, ErrUnknownMeterID
		}
	}

	desc := s.EventRepo.GetTamperDescription(e.TamperCode)
	if desc == "" {
		desc = "Unknown Tamper Code"
	}

	// If the caller didn't set a timestamp, default to now.
	if e.Timestamp.IsZero() {
		e.Timestamp = time.Now()
	}

	if e.EventOccur == 1 {
		existing, found, err := s.NotifRepo.GetLatestPending(e.MeterID, e.TamperCode)
		if err != nil {
			return models.Notification{}, err
		}
		if !found {
			return models.Notification{
				MeterID:           e.MeterID,
				TamperCode:        e.TamperCode,
				TamperDescription: desc,
				Type:              "tamper",
				Level:             1,
				Status:            "resolved",
				Message:           fmt.Sprintf("Tamper Alert: %s ended", desc),
				Timestamp:         e.Timestamp,
			}, nil
		}

		existing.Status = "resolved"
		existing.Message = fmt.Sprintf("Tamper Alert: %s ended", desc)
		existing.Timestamp = e.Timestamp
		existing.ResolvedAt = &e.Timestamp
		existing.TamperDescription = desc
		if err := s.NotifRepo.UpdateNotification(&existing); err != nil {
			return models.Notification{}, err
		}
		if s.WS != nil {
			s.WS.Broadcast(existing)
		}
		return existing, nil
	}

	existing, found, err := s.NotifRepo.GetLatestPending(e.MeterID, e.TamperCode)
	if err != nil {
		return models.Notification{}, err
	}
	if found {
		existing.TamperDescription = desc
		return existing, nil
	}

	notification := models.Notification{
		MeterID:           e.MeterID,
		TamperCode:        e.TamperCode,
		TamperDescription: desc, // used for API response only (not stored in DB).
		Type:              "tamper",
		Level:             1,
		Status:            "pending",
		Message:           fmt.Sprintf("Tamper Alert: %s detected", desc),
		Timestamp:         e.Timestamp,
	}

	if err := s.NotifRepo.SaveNotificationWithResult(&notification); err != nil {
		return models.Notification{}, err
	}

	if s.WS != nil {
		s.WS.Broadcast(notification)
	}

	return notification, nil
}

func (s *NotificationService) GetNotificationPage(
	filter repository.NotificationFilter,
	page int,
	pageSize int,
) (models.NotificationListResponse, error) {
	if page <= 0 {
		page = 1
	}
	if pageSize <= 0 || pageSize > 500 {
		pageSize = 25
	}
	filter.Limit = pageSize
	filter.Offset = (page - 1) * pageSize

	items, err := s.NotifRepo.GetFiltered(filter)
	if err != nil {
		return models.NotificationListResponse{}, err
	}
	total, err := s.NotifRepo.CountFiltered(filter)
	if err != nil {
		return models.NotificationListResponse{}, err
	}

	for i := range items {
		items[i].TamperDescription = s.EventRepo.GetTamperDescription(items[i].TamperCode)
	}
	responseItems := make([]models.NotificationResponse, 0, len(items))
	for _, item := range items {
		responseItems = append(responseItems, toNotificationResponse(item))
	}

	totalPages := int64(math.Ceil(float64(total) / float64(pageSize)))
	if total == 0 {
		totalPages = 0
	}
	return models.NotificationListResponse{
		Items: responseItems,
		Pagination: models.PaginationMeta{
			Page:       page,
			PageSize:   pageSize,
			TotalItems: total,
			TotalPages: totalPages,
		},
	}, nil
}

func toNotificationResponse(item models.Notification) models.NotificationResponse {
	return models.NotificationResponse{
		ID:                item.ID,
		MeterID:           item.MeterID,
		TamperCode:        item.TamperCode,
		TamperDescription: item.TamperDescription,
		Message:           item.Message,
		Type:              item.Type,
		Level:             item.Level,
		Status:            item.Status,
		ScheduledFor:      formatOptionalIST(item.ScheduledFor),
		SentAt:            formatOptionalIST(item.SentAt),
		ResolvedAt:        formatOptionalIST(item.ResolvedAt),
		Timestamp:         utils.FormatIST(item.Timestamp),
	}
}

func formatOptionalIST(t *time.Time) *string {
	if t == nil {
		return nil
	}
	value := utils.FormatIST(*t)
	return &value
}

func (s *NotificationService) ProcessEvents() {

	events, _ := s.EventRepo.GetUnprocessedEvents()

	for _, e := range events {
		_, _ = s.ProcessEvent(e)

		s.EventRepo.MarkProcessed(e.ID)
	}
}