package main

import (
	"log"
	"os"

	"github.com/ak-ko/ghop.git/config"
	"github.com/ak-ko/ghop.git/db"
	mysqlDriver "github.com/go-sql-driver/mysql"
	"github.com/golang-migrate/migrate/v4"
	mysqlMigrate "github.com/golang-migrate/migrate/v4/database/mysql"
	_ "github.com/golang-migrate/migrate/v4/source/file"
)

func main(){
	database, err := db.InitDB(mysqlDriver.Config{
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

	driver, err := mysqlMigrate.WithInstance(database, &mysqlMigrate.Config{})
	if err != nil {
		log.Fatal(err)
	}


	m, err := migrate.NewWithDatabaseInstance(
		"file://cmd/migrate/migrations",
		"mysql",
		driver,
	); if err != nil {
		log.Fatal(err)
	}

	lastIndex := len(os.Args) - 1
	cmd := os.Args[lastIndex]


	if cmd == "up" {
		if err := m.Up(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}

	if cmd == "down" {
		if err := m.Down(); err != nil && err != migrate.ErrNoChange {
			log.Fatal(err)
		}
	}
}