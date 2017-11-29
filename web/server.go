package web

import (
	"fmt"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
	"net/http"
	"os"
)

func SetupServer() {
	port := os.Getenv("LOCATION_PORT")
	if port == "" {
		port = "4343"
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/drivers/{id}/heart_beat", TokenAuth(WriteHeartBeat))
	router.HandleFunc("/api/v2/drivers/{id}/update_current_location", TokenAuth(UpdateCurrentLocation))
	log.Info("Starting ... ")
	log.Infof("Listening on port %s ... ", port)
	setupHeartBeatTimer()
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Panic(err)
	}
}
