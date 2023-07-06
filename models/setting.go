package models

import (
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Setting struct {
	gorm.Model
	CompanyName       string `gorm:"not null;" json:"company_name"`
	CompanyAddress    string `gorm:"not null;" json:"company_address"`
	CompanyEmail      string `gorm:"not null;" json:"company_email"`
	CompanyMobile     string `gorm:"not null;" json:"company_mobile"`
	Logo              string `gorm:"" json:"logo"`
	Favicon           string `gorm:"" json:"favicon"`
	VatNumber         string `gorm:"" json:"vat_icon"`
	TimeZone          string `gorm:"" json:"time_zone"`
	DateFormat        string `gorm:"" json:"date_format"`
	CurrencySymbol    string `gorm:"type:varchar(10);" json:"currency_symbol"`
	DecimalSeparator  string `gorm:"type:varchar(3);" json:"decimal_separator"`
	ThousandSeparator string `gorm:"type:varchar(3);" json:"thousand_separator"`
	NumberOfDecimal   int8   `gorm:"" json:"number_of_decimal"`
	CurrencyPosition  string `gorm:"" json:"currency_position"`
}

// GetSettings returns all the settings data
func (setting *Setting) GetSettings() map[string]interface{} {
	err := db.First(&setting).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = setting

	return res
}

// Update updates all the settings in the db
func (setting *Setting) Update() map[string]interface{} {
	stng := Setting{}

	err := db.First(&stng).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", stng.ID).Updates(&setting).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = setting

	return res
}
