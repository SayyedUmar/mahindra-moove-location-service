package web

import (
	"fmt"
	"github.com/fvbock/endless"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"os"
)

func SetupServer() {
	port := os.Getenv("LOCATION_PORT")
	if port == "" {
		port = "4343"
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/drivers/{id}/heart_beat", TokenAuth(WriteHeartBeat)).Methods("POST")
	router.HandleFunc("/api/v2/drivers/{id}/update_current_location", TokenAuth(UpdateCurrentLocation)).Methods("POST")
	router.HandleFunc("/api/v3/drivers/{id}/location", TokenAuth(LocationSocket))
	log.Info("Starting ... ")
	log.Infof("Listening on port %s ... ", port)
	setupHeartBeatTimer()
	go hub.Run()

	handler := handlers.LoggingHandler(os.Stdout, router)

	err := endless.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		log.Error(err)
	}
}
