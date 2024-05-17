package main

import (
	"DiplomaV2/config"
	"DiplomaV2/database"
	"DiplomaV2/server"
	"fmt"
	_ "github.com/lib/pq"
)

func main() {
	// This stupis shit opens up some type of config or smth idk
	conf := config.GetConfig()
	// I am so gay. I love Dimash a lot and want to kiss him. I send nudes to madi
	fmt.Println("!!BEKARIZZLER 1.0 COOL PROGRAM ACTIVATED!!")
	db := database.NewPostgresDatabase(conf)
	server.NewEchoServer(conf, db).Start()
}
