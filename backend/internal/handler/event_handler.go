package handler

import (
	"encoding/json"
	"net/http"

	"goproject/internal/models"
	"goproject/internal/service"
	"goproject/internal/utils"
)

type EventHandler struct {
	Service *service.NotificationService
}

func (h *EventHandler) Receive(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var e models.TamperEvent

	err := json.NewDecoder(r.Body).Decode(&e)
	if err != nil {
		http.Error(w, "invalid input", http.StatusBadRequest)
		return
	}

	if e.MeterID == "" || e.TamperCode == 0 {
		http.Error(w, "missing fields", http.StatusBadRequest)
		return
	}

	if h.Service == nil {
		http.Error(w, "service unavailable", http.StatusInternalServerError)
		return
	}

	eventID, created, svcErr := h.Service.IngestEvent(e)
	if svcErr != nil {
		if svcErr == service.ErrUnknownMeterID {
			http.Error(w, "invalid meter_id", http.StatusBadRequest)
			return
		}
		http.Error(w, "notification error", http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(models.EventResponse{
		Status:     "ok",
		EventID:    eventID,
		MeterID:    e.MeterID,
		TamperCode: e.TamperCode,
		Timestamp:  utils.FormatIST(created.Timestamp),
	})
}