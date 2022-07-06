package main

import (
	"memoapi/config"
	"memoapi/infra/db"
	"memoapi/infra/server"

	log "github.com/sirupsen/logrus"
)

func main() {
	conf := config.Load()
	db, err := db.NewPostgresDB(conf.RdbConfig.ConnectionString())
	if err != nil {
		log.Fatalf(err.Error())
	}
	s := server.NewServer(conf.AppPort(), conf.AppHost(), db, nil)
	s.Start()
}
