package main

import (
	"github.com/MOOVE-Network/location_service/db"
	"github.com/MOOVE-Network/location_service/version"
	web "github.com/MOOVE-Network/location_service/web"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
	log.SetFormatter(&log.JSONFormatter{})
}

func main() {
	version.PrintVersion()
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
