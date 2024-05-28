package main

import (
	"DiplomaV2/backend/internal/config"
	"DiplomaV2/backend/internal/database"
	"DiplomaV2/backend/server"
	_ "github.com/lib/pq"
)

func main() {
	conf := config.GetConfig()
	db := database.NewPostgresDatabase(conf)
	server.NewEchoServer(conf, db).Start()
}
