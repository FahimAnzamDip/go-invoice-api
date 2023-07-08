package models

import (
	"log"

	"github.com/fahimanzamdip/go-invoice-api/config"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/patrickmn/go-cache"
	"gorm.io/gorm"
)

const SettingsCacheKey string = "settings"

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
	stngs := &Setting{}
	if settings, found := config.DBCache.Get(SettingsCacheKey); found {
		stngs = settings.(*Setting)
	} else {
		err := db.First(&stngs).Error
		if err != nil {
			return u.Message(false, err.Error())
		}
		config.DBCache.Set(SettingsCacheKey, stngs, cache.DefaultExpiration)
	}

	res := u.Message(true, "")
	res["data"] = stngs

	return res
}

// Update updates settings in the db
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

	config.DBCache.Delete(SettingsCacheKey)

	res := u.Message(true, "")
	res["data"] = setting

	return res
}

// GetSettingsIntr is for using settings only in backend
func (setting *Setting) AppGetSettings() *Setting {
	settings := Setting{}

	err := config.GetDB().First(&settings).Error
	if err != nil {
		log.Println(err.Error())
	}

	return &settings
}
