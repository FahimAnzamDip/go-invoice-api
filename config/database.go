package config

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

func init() {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV for DB connection")
		return
	}
}

func Connect() {
	dbDriver := os.Getenv("db_driver")
	dbUser := os.Getenv("db_user")
	dbPass := os.Getenv("db_pass")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	dbName := os.Getenv("db_name")

	var dsn string
	if dbDriver == "pgsql" {
		// dsn for postgres sql
		dsn = "host=" + dbHost + " " + "user=" + dbUser + " " + "password=" + dbPass + " " + "dbname=" + dbName + " " + "port=" + dbPort + " sslmode=disable TimeZone=Asia/Dhaka"
	} else if dbDriver == "mysql" {
		// dsn for mysql
		dsn = dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"
	} else {
		panic("Unsupported database driver! Please check you .env file.")
	}

	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}

func GetDB() *gorm.DB {
	return db
}
