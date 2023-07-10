package handlers

import (
	"encoding/json"
	"net/http"
	"os"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func SettingsDataHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./data/settings.json"
	file, err := os.Open(filePath)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}
	defer file.Close()

	fileInfo, err := file.Stat()
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	fileSize := fileInfo.Size()
	fileContent := make([]byte, fileSize)

	_, err = file.Read(fileContent)
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
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
