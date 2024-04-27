package main

import (
	"DiplomaV2/cmd/app"
	"DiplomaV2/config"
	"DiplomaV2/database"
)

func main() {
	conf := config.GetConfig()

	db := database.NewPostgresDatabase(conf)

	app.NewEchoServer(conf, db).Start()
}
