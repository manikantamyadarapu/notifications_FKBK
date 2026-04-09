package main

import (
	"log"
	"net/http"
	"os"

	"github.com/robfig/cron/v3"

	"goproject/config"
	"goproject/docs"
	"goproject/internal/handler"
	"goproject/internal/models"
	"goproject/internal/repository"
	"goproject/internal/service"
)

func main() {
	db := config.ConnectGormDB()

	// Auto-migrate required tables for the assessment.
	_ = db.AutoMigrate(
		&models.Meter{},
		&models.TamperCodeDesc{},
		&models.TamperEvent{},
		&models.Notification{},
	)

	eventRepo := &repository.EventRepository{DB: db}
	notifRepo := &repository.NotificationRepository{DB: db}

	ws := service.NewWSManager()

	notifService := &service.NotificationService{
		EventRepo: eventRepo,
		NotifRepo: notifRepo,
		WS:        ws,
	}

	c := cron.New()
	_, err := c.AddFunc("@every 5s", func() {
		notifService.ProcessEvents()
	})
	if err != nil {
		log.Fatalf("failed to register cron job: %v", err)
	}
	c.Start()
	defer c.Stop()

	eventHandler := &handler.EventHandler{Service: notifService}
	notifHandler := &handler.NotificationHandler{Service: notifService}

	mux := http.NewServeMux()
	mux.HandleFunc("/event", eventHandler.Receive)
	mux.HandleFunc("/notifications", notifHandler.GetAll)
	mux.HandleFunc("/ws", ws.Handle)
	docs.RegisterRoutes(mux)

	port := os.Getenv("SERVER_PORT")
	if port == "" {
		port = "8080"
	}

	log.Printf("Server running on :%s", port)
	log.Fatal(http.ListenAndServe(":"+port, withCORS(mux)))
}

func withCORS(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// For local frontend integration, keep CORS fully open.
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET,POST,OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == http.MethodOptions {
			w.WriteHeader(http.StatusNoContent)
			return
		}

		next.ServeHTTP(w, r)
	})
}
