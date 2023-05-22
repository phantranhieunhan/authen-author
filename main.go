package main

import (
	"database/sql"
	"log"

	redisAdapter "github.com/phantranhieunhan/authen-author/adapter/redis"
	"github.com/phantranhieunhan/authen-author/api"
	"github.com/phantranhieunhan/authen-author/db/redis"
	db "github.com/phantranhieunhan/authen-author/db/sqlc"
	"github.com/phantranhieunhan/authen-author/logger"
	"github.com/phantranhieunhan/authen-author/util"
)

func main() {
	config, err := util.LoadConfig(".")
	if err != nil {
		log.Fatal("cannot load config")
	}
	// init logger
	logger.Newlogger(logger.ConfigLogger{})

	// demo logger into file
	// forever := make(chan int)
	// log := logger.GetLogger()
	// for tick := range time.Tick(time.Millisecond) {
	// 	log.Debug(tick)
	// }
	// <-forever

	conn, err := sql.Open(config.DBDriver, config.DBSource)
	if err != nil {
		log.Fatal("cannot connect to db")
	}

	store := db.NewStore(conn)
	redisAdapter, err := redisAdapter.New(
		redisAdapter.WithAddress("redis-16062.c267.us-east-1-4.ec2.cloud.redislabs.com:16062"),
		redisAdapter.WithDatabase(0),
		redisAdapter.WithPassword("cUPZU2TtmHTTojHGiFVhkDE7VHKu2t6E"),
	)
	if err != nil {
		log.Fatal("cannot connect redis")
	}
	sessionStore := redis.NewStore(redisAdapter)
	runGinServer(config, store, sessionStore)
}

func runGinServer(config util.Config, store db.Store, session *redis.RedisStore) {
	server, err := api.NewServer(config, store, session)
	if err != nil {
		log.Fatal("cannot create server")
	}

	err = server.Start(config.ServerAddress)
	if err != nil {
		log.Fatal("cannot start server")
	}
}
