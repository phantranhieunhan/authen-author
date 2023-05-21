package main

import (
	"database/sql"
	"log"
	"time"

	"github.com/phantranhieunhan/authen-author/api"
	db "github.com/phantranhieunhan/authen-author/db/sqlc"
	"github.com/phantranhieunhan/authen-author/logger"
	"github.com/phantranhieunhan/authen-author/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}

	logger.Newlogger(logger.ConfigLogger{
		
	})
	forever := make(chan int)
	log := logger.GetLogger()
	for tick := range time.Tick(time.Millisecond) {
		log.Debug(tick)
	}
	<-forever
	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(conn)
	runGinServer(config, store)
}

func runGinServer(config util.Config, store db.Store) {
	server, err := api.NewServer(config, store)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
