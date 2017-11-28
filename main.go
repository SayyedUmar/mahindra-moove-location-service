package main

import (
	"github.com/MOOVE-Network/location_service/db"
	ident "github.com/MOOVE-Network/location_service/identity"
	log "github.com/sirupsen/logrus"
	"io"
	"os"
)

func init() {
	log.SetOutput(os.Stdout)
	log.SetLevel(log.DebugLevel)
}

func main() {
	conn := db.InitSQLConnection()
	db.SetActiveDB(conn)
	defer closeConn(conn)
	identity := ident.FetchIdentityByUID(conn, "moove.dinesh1651@gmail.com")
	if identity.IsValid("GWyFWywJ3mf5TJe77DBIHw", "v-pWLlS4pwf_L9qbx-FtYw") {
		log.Info("Yay! passwords match")
	}
}

func closeConn(closable io.Closer) {
	err := closable.Close()
	if err != nil {
		log.Panic(err)
	}
}
