package handler

import (
	"encoding/json"
	"net/http"
	"strconv"
	"time"

	"goproject/internal/repository"
	"goproject/internal/service"
)

type NotificationHandler struct {
	Service *service.NotificationService
}

func (h *NotificationHandler) GetAll(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodGet {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	if h.Service == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	q := r.URL.Query()

	var filter repository.NotificationFilter
	filter.MeterID = q.Get("meter_id")
	filter.Type = q.Get("type")

	if v := q.Get("tamper_code"); v != "" {
		if n, err := strconv.Atoi(v); err == nil {
			filter.TamperCode = &n
		}
	}

	if v := q.Get("from"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.From = &t
		}
	}
	if v := q.Get("to"); v != "" {
		if t, err := time.Parse(time.RFC3339, v); err == nil {
			filter.To = &t
		}
	}

	page := parseInt(q.Get("page"), 1)
	pageSize := parseInt(q.Get("page_size"), 25)

	data, err := h.Service.GetNotificationPage(filter, page, pageSize)
	if err != nil {
		http.Error(w, "db error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(data)
}

func parseInt(v string, fallback int) int {
	if v == "" {
		return fallback
	}
	n, err := strconv.Atoi(v)
	if err != nil {
		return fallback
	}
	return n
}