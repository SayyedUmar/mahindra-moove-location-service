package web

import (
	"fmt"
	"github.com/DataDog/dd-trace-go/tracer"
	"github.com/DataDog/dd-trace-go/tracer/contrib/gorilla/muxtrace"
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
	tracer := muxtrace.NewMuxTracer("my-web-app", tracer.DefaultTracer)
	tracer.HandleFunc(router, "/api/v1/drivers/{id}/heart_beat", TokenAuth(WriteHeartBeat)).Methods("POST")
	//router.HandleFunc("/api/v1/drivers/{id}/heart_beat", TokenAuth(WriteHeartBeat)).Methods("POST")
	tracer.HandleFunc(router, "/api/v2/drivers/{id}/update_current_location", TokenAuth(UpdateCurrentLocation)).Methods("POST")
	//router.HandleFunc("/api/v2/drivers/{id}/update_current_location", TokenAuth(UpdateCurrentLocation)).Methods("POST")
	tracer.HandleFunc(router, "/api/v3/drivers/{id}/location", TokenAuth(LocationSocket))
	//router.HandleFunc("/api/v3/drivers/{id}/location", TokenAuth(LocationSocket))
	log.Info("Starting ... ")
	log.Infof("Listening on port %s ... ", port)
	setupHeartBeatTimer()
	err := http.ListenAndServe(fmt.Sprintf(":%s", port), router)
	if err != nil {
		log.Panic(err)
	}
}
