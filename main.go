package main

import (
	"io"
	"os"
	"os/signal"

	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/models"
	"github.com/MOOVE-Network/location_service/redis"
	"github.com/MOOVE-Network/location_service/services"
	"github.com/MOOVE-Network/location_service/version"
	"github.com/MOOVE-Network/location_service/web"
	_ "github.com/joho/godotenv/autoload"
	log "github.com/sirupsen/logrus"
	"github.com/stvp/rollbar"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	version.PrintVersion()
	setupRollbar()
	defer rollbar.Wait()
	conn := db.InitSQLConnection()
	db.SetActiveDB(conn)
	defer closeConn(conn)

	pgConn := models.InitSQLConnection()
	models.SetActiveDB(pgConn)
	defer closeConn(pgConn)
	models.RunMigrations()

	redisClient := redis.SetupRedis()
	redis.SetClient(redisClient)
	defer closeConn(redisClient)

	os.Setenv("LOCATION_MAPS_API_KEY", "AIzaSyBqCo3g-P2P7MMSedPt-izL-oq1gL4804w")
	os.Setenv("FCM_API_KEY", "AAAAEXz4c6s:APA91bGE8CJdgnJQPsxgVAxM7asjVEdShM6wpq6jme5sShi4S095MdCuKOGSYuaE3AlrVYKeNrjcrfJQtIoafXabnOKI5kZOBarGkPKlo-pZah3f_iT6rvLA49CpJvx7UOf6JSzL0Zb4")
	os.Setenv("FCM_TOPIC_PREFIX", "local")

	services.InitDurationService(os.Getenv("LOCATION_MAPS_API_KEY"))
	services.InitGoogleRoadsService(os.Getenv("LOCATION_MAPS_API_KEY"))
	services.InitNotificationService(os.Getenv("FCM_API_KEY"), os.Getenv("FCM_TOPIC_PREFIX"))
	cancelable := make(chan bool)
	go services.StartETAServiceTimer(cancelable, db.CurrentDB())
	go cancelOnSignal(cancelable)

	web.SetupServer(db.CurrentDB())
}

func cancelOnSignal(cancelable chan bool) {
	sigIntChan := make(chan os.Signal, 1)
	signal.Notify(sigIntChan, os.Interrupt, os.Kill)
	//block on receiving signal
	_ = <-sigIntChan
	cancelable <- true
}

func closeConn(closable io.Closer) {
	err := closable.Close()
	if err != nil {
		log.Panic(err)
	}
}

func setupRollbar() {
	rollBarToken := os.Getenv("ROLLBAR_TOKEN")
	if rollBarToken != "" {
		rollbar.Token = rollBarToken
		rollbar.Environment = "production"
	}
}
