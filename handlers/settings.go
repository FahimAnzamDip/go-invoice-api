package handlers

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func SettingsDataHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./data/settings.json"
	fileContent, err := ioutil.ReadFile(filePath)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
	}

	w.Header().Set("Content-Type", "application/json")
	w.Write(fileContent)
}

func GetSettingsHandler(w http.ResponseWriter, r *http.Request) {
	setting := &models.Setting{}

	data := setting.GetSettings()
	u.Respond(w, data)
}

func UpdateSettingsHandler(w http.ResponseWriter, r *http.Request) {
	setting := &models.Setting{}
	err := json.NewDecoder(r.Body).Decode(setting)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := setting.Update()
	u.Respond(w, res)
}
