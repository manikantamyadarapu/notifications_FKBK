package models

type EventResponse struct {
	Status     string `json:"status"`
	EventID    uint   `json:"event_id"`
	MeterID    string `json:"meter_id"`
	TamperCode int    `json:"tamper_code"`
	Timestamp  string `json:"timestamp"`
}

type NotificationResponse struct {
	ID                uint    `json:"id"`
	MeterID           string  `json:"meter_id"`
	TamperCode        int     `json:"tamper_code"`
	TamperDescription string  `json:"tamper_description"`
	Message           string  `json:"message"`
	Type              string  `json:"type"`
	Level             int     `json:"level"`
	Status            string  `json:"status"`
	ScheduledFor      *string `json:"scheduled_for"`
	SentAt            *string `json:"sent_at"`
	ResolvedAt        *string `json:"resolved_at"`
	Timestamp         string  `json:"timestamp"`
}

type PaginationMeta struct {
	Page       int   `json:"page"`
	PageSize   int   `json:"page_size"`
	TotalItems int64 `json:"total_items"`
	TotalPages int64 `json:"total_pages"`
}

type NotificationListResponse struct {
	Items      []NotificationResponse `json:"items"`
	Pagination PaginationMeta `json:"pagination"`
}

