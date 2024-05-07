package main

import (
	"DiplomaV2/config"
	"DiplomaV2/database"
	"DiplomaV2/server"
)

func main() {
	conf := config.GetConfig()

	db := database.NewPostgresDatabase(conf)

	server.NewEchoServer(conf, db).Start()
}
