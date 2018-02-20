package web

import (
	"encoding/json"
	"fmt"
	"net/http"
	"os"

	"github.com/stvp/rollbar"

	"github.com/MOOVE-Network/location_service/version"
	"github.com/fvbock/endless"
	"github.com/gorilla/handlers"
	"github.com/gorilla/mux"
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

type RollbarLogger struct {
}

func (r *RollbarLogger) Println(args ...interface{}) {
	if rollbarEnabled() {
		rollbar.Message(rollbar.ERR, fmt.Sprintf("%v", args))
	}
	log.Error(args)
}

func SetupServer() {
	port := os.Getenv("LOCATION_PORT")
	if port == "" {
		port = "4343"
	}
	router := mux.NewRouter()
	router.Handle("/api/v3/metrics", promhttp.Handler())
	router.HandleFunc("/api/v1/drivers/{id}/heart_beat", LogRequestsMiddleware(TokenAuth(WriteHeartBeat))).Methods("POST")
	router.HandleFunc("/api/v2/drivers/{id}/update_current_location", LogRequestsMiddleware(TokenAuth(UpdateCurrentLocation))).Methods("POST")
	router.HandleFunc("/api/v3/trips/{id}/eta", LogRequestsMiddleware(GetTripETA)).Methods("GET")
	router.HandleFunc("/api/v3/trips/{id}/summary", LogRequestsMiddleware(TripSummary))
	router.HandleFunc("/api/v3/drivers/{id}/location", Auth(LocationSocket))

	router.HandleFunc("/api/v3/version", EchoVersion)
	log.Info("Starting ... ")
	log.Infof("Listening on port %s ... ", port)
	setupHeartBeatTimer()
	setupTripLocationTimer()
	go hub.Run()

	handler := handlers.LoggingHandler(os.Stdout, router)
	withRecovery := handlers.RecoveryHandler(handlers.RecoveryLogger(&RollbarLogger{}), handlers.PrintRecoveryStack(true))(handler)

	err := endless.ListenAndServe(fmt.Sprintf(":%s", port), withRecovery)
	if err != nil {
		log.Error(err)
	}
}

func rollbarEnabled() bool {
	return os.Getenv("ROLLBAR_TOKEN") != ""
}
