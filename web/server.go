package web

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fvbock/endless"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	log "github.com/sirupsen/logrus"
)

func LogRequestsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Infof("Serving %s %s ", r.Method, r.URL.Path)
		next(rw, r)
	}
}

func SetupServer() {
	port := os.Getenv("LOCATION_PORT")
	if port == "" {
		port = "4343"
	}
	router := mux.NewRouter()
	router.HandleFunc("/api/v1/drivers/{id}/heart_beat", LogRequestsMiddleware(TokenAuth(WriteHeartBeat))).Methods("POST")
	router.HandleFunc("/api/v2/drivers/{id}/update_current_location", LogRequestsMiddleware(TokenAuth(UpdateCurrentLocation))).Methods("POST")
	router.HandleFunc("/api/v3/drivers/{id}/location", Auth(LocationSocket))
	log.Info("Starting ... ")
	log.Infof("Listening on port %s ... ", port)
	setupHeartBeatTimer()
	setupTripLocationTimer()
	go hub.Run()

	handler := handlers.LoggingHandler(os.Stdout, router)

	err := endless.ListenAndServe(fmt.Sprintf(":%s", port), handler)
	if err != nil {
		log.Error(err)
	}
}
