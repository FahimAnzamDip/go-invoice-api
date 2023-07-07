package models

import (
	"log"

	"github.com/fahimanzamdip/go-invoice-api/config"
	"gorm.io/gorm"
)

var db *gorm.DB
var debug bool

func init() {
	debug = true

	config.Connect()
	db = config.GetDB()

	if debug {
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
		&Purpose{},
		&Expense{},
		&Setting{},
	)

	return err
}

func seed() error {
	var err error

	// Create a Super User
	// err = db.Create(&User{
	// 	Name:     "Administrator",
	// 	Email:    "super.admin@test.com",
	// 	Mobile:   "12345678901",
	// 	Password: "$2a$12$1ojxUBODleIVVvFo1Lvysu/WSpVXi8yUb2zq6SIWJe6M9OJv3yRf2", // 123456
	// 	IsAdmin:  true,
	// }).Error
	// if err != nil {
	// 	return err
	// }

	// Create settings
	// err = db.Create(&Setting{
	// 	CompanyName:       "Company",
	// 	CompanyAddress:    "Dhaka, Bangladesh",
	// 	CompanyEmail:      "company@mail.com",
	// 	CompanyMobile:     "12345678901",
	// 	Logo:              "public/uploads/logo.png",
	// 	Favicon:           "public/uploads/favicon.png",
	// 	VatNumber:         "",
	// 	TimeZone:          "Asia/Dhaka",
	// 	DateFormat:        "d-m-Y",
	// 	CurrencySymbol:    "$",
	// 	DecimalSeparator:  ".",
	// 	ThousandSeparator: ",",
	// 	NumberOfDecimal:   2,
	// 	CurrencyPosition:  "prefix_with_space",
	// }).Error
	// if err != nil {
	// 	return err
	// }

	return err
}
