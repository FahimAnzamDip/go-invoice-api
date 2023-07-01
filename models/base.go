package models

import (
	"log"
	"os"

	"github.com/joho/godotenv"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

var db *gorm.DB

var DEBUGMODE bool

func init() {
	DEBUGMODE = true

	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV for DB connection")
		return
	}
	connect()
	db = GetDB()

	if DEBUGMODE {
		// Auto Migration
		err := migrate()
		if err != nil {
			log.Println("Auto migration failed! Please check configuration.")
			return
		}
		// Seed Databasee
		err = seed()
		if err != nil {
			log.Println("Database seeding failed! Please check configuration.")
			return
		}
	}
}

func connect() {
	dbUser := os.Getenv("db_user")
	dbPass := os.Getenv("db_pass")
	dbHost := os.Getenv("db_host")
	dbPort := os.Getenv("db_port")
	dbName := os.Getenv("db_name")

	// dsn := "host=" + dbHost + " " + "user=" + dbUser + " " + "password=" + dbPass + " " + "dbname=" + dbName + " " + "port=" + dbPort + " sslmode=disable TimeZone=Asia/Dhaka"

	dsn := dbUser + ":" + dbPass + "@tcp(" + dbHost + ":" + dbPort + ")/" + dbName + "?charset=utf8mb4&parseTime=True&loc=Local"

	log.Println(dsn)

	d, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		panic(err)
	}
	db = d
}

func migrate() error {
	err := db.AutoMigrate(
		&User{},
		&Category{},
		&Tag{},
		&Product{},
		&Client{},
		&Invoice{},
		&Tax{},
		&InvoiceProduct{},
		&Payment{},
	)

	return err
}

func seed() error {
	// Create a Super User
	err := db.Create(&User{
		Name:     "Administrator",
		Email:    "super.admin@test.com",
		Mobile:   "12345678901",
		Password: "$2a$12$1ojxUBODleIVVvFo1Lvysu/WSpVXi8yUb2zq6SIWJe6M9OJv3yRf2", // 123456
		IsAdmin:  true,
	}).Error

	return err
}

func GetDB() *gorm.DB {
	return db
}
