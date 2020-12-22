package main

import (
	"database/sql"
	"log"

	"github.com/spruce-cq/simplebank/api"
	db "github.com/spruce-cq/simplebank/db/sqlc"
	"github.com/spruce-cq/simplebank/util"

	_ "github.com/lib/pq"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config: ", err)
	}
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db: ", err)
	}

	store := db.NewStore(conn)
	server := api.NewServer(store)
	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start the server: ", err)
	}

}
