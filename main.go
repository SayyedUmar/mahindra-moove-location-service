package main

import (
	"io"
	"os"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/services"
	"github.com/MOOVE-Network/location_service/version"
	web "github.com/MOOVE-Network/location_service/web"
	log "github.com/sirupsen/logrus"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	version.PrintVersion()
	services.InitDurationService(os.Getenv("LOCATION_MAPS_API_KEY"))
	services.InitNotificationService(os.Getenv("FCM_API_KEY"), os.Getenv("FCM_TOPIC_PREFIX"))
	conn := db.InitSQLConnection()
	db.SetActiveDB(conn)
	defer closeConn(conn)
	web.SetupServer()
}

func closeConn(closable io.Closer) {
	err := closable.Close()
	if err != nil {
		log.Panic(err)
	}
}
