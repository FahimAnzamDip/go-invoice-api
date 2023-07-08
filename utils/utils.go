package utils

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/config"
	"github.com/leekchan/accounting"
)

type CurrencySetting struct {
	CurrencySymbol    string
	DecimalSeparator  string
	ThousandSeparator string
	NumberOfDecimal   int8
	CurrencyPosition  string
}

type DateSetting struct {
	TimeZone   string
	DateFormat string
}

func Message(status bool, message string) map[string]interface{} {
	return map[string]interface{}{"status": status, "message": message}
}

func Respond(w http.ResponseWriter, data map[string]interface{}) {
	w.Header().Add("Content-Type", "application/json")
	json.NewEncoder(w).Encode(data)
}

func FormatPrice(value float32) string {
	settings := CurrencySetting{}
	err := config.GetDB().Table("settings").First(&settings).Error
	if err != nil {
		log.Println(err.Error())
	}

	ac := accounting.Accounting{
		Symbol:    "",
		Precision: int(settings.NumberOfDecimal),
		Thousand:  settings.ThousandSeparator,
		Decimal:   settings.DecimalSeparator,
	}

	formattedValue := ac.FormatMoney(value)

	switch settings.CurrencyPosition {
	case "prefix":
		formattedValue = settings.CurrencySymbol + formattedValue
	case "prefix_with_space":
		formattedValue = settings.CurrencySymbol + " " + formattedValue
	case "suffix":
		formattedValue += settings.CurrencySymbol
	case "suffix_with_space":
		formattedValue += " " + settings.CurrencySymbol
	}

	return formattedValue
}

func FormatDate() {

}
