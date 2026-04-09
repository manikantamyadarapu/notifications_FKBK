package api

import (
	"net/http"

	"goproject/internal/handler"
	"goproject/internal/service"
)

// SetupRoutes registers all routes
func SetupRoutes(
	eventHandler *handler.EventHandler,
	notifHandler *handler.NotificationHandler,
	ws *service.WSManager,
) {

	http.HandleFunc("/event", eventHandler.Receive)

	http.HandleFunc("/notifications", notifHandler.GetAll)

	http.HandleFunc("/ws", ws.Handle)
}