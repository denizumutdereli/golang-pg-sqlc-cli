package main

import (
	"database/sql"
	"log"

	"github.com/denizumutdereli/golang-pg-sqlc-cli/api"
	db "github.com/denizumutdereli/golang-pg-sqlc-cli/db/sqlc"

	_ "github.com/lib/pq"
)

const (
	dbDriver      = "postgres"
	dbSource      = "postgresql://root:secret@192.168.99.100:5432/mservice?sslmode=disable"
	serverAddress = "0.0.0.0:8080"
)

var ()

func main() {

	var err error

	conn, err := sql.Open(dbDriver, dbSource)

	if err != nil {
		log.Fatal("DB Connection error:", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)

	err = server.Start(serverAddress)

	if err != nil {
		log.Fatal("Can not start the server:", err)
	}

}
