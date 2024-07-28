package main

import (
	"log"

	"github.com/ak-ko/ghop.git/cmd/api"
	"github.com/ak-ko/ghop.git/config"
	"github.com/ak-ko/ghop.git/db"
	"github.com/go-sql-driver/mysql"
)

func main(){
	database, err := db.InitDB(mysql.Config{
		User: config.ENV.DB_USER,
		Passwd: config.ENV.DB_PASSWORD,
		Addr: config.ENV.DB_HOST,
		DBName: config.ENV.DB_NAME,
		Net: "tcp",
		AllowNativePasswords: true,
		ParseTime: true,
	}); if err != nil {
		log.Fatal(err)
	}

	server := api.NewAPIServer(":8080", database)
	if err := server.Run(); err != nil {
		log.Fatal(err)
	}
}