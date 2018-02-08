package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/MOOVE-Network/location_service/version"
	"github.com/fvbock/endless"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	log "github.com/sirupsen/logrus"
)

func LogRequestsMiddleware(next http.HandlerFunc) http.HandlerFunc {
	return func(rw http.ResponseWriter, r *http.Request) {
		log.Infof("Serving %s %s ", r.Method, r.URL.Path)
		next(rw, r)
	}
}

func EchoVersion(rw http.ResponseWriter, req *http.Request) {
	enc := json.NewEncoder(rw)
	err := enc.Encode(version.GetVersion())
	if err != nil {
		ErrorWithMessage(err.Error()).Respond(rw, 500)
	}
	rw.Header().Add("Content-Type", "application/json")
}

func SetupServer() {
	port := os.Getenv("LOCATION_PORT")
	if port == "" {
		port = "4343"
	}
	router := mux.NewRouter()
	router.Handle("/metrics", promhttp.Handler())
	router.HandleFunc("/api/v1/drivers/{id}/heart_beat", LogRequestsMiddleware(TokenAuth(WriteHeartBeat))).Methods("POST")
	router.HandleFunc("/api/v2/drivers/{id}/update_current_location", LogRequestsMiddleware(TokenAuth(UpdateCurrentLocation))).Methods("POST")
	router.HandleFunc("/api/v3/trips/{id}/eta", LogRequestsMiddleware(TokenAuth(GetTripETA))).Methods("GET")
	prometheus.InstrumentHandlerFunc("socket", Auth(LocationSocket))
	router.HandleFunc("/api/v3/version", EchoVersion)
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
