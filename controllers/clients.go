package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexClientHandler(w http.ResponseWriter, r *http.Request) {
	client := &models.Client{}

	data := client.Index()
	u.Respond(w, data)
}

// Store function creates a new entry
func StoreClientHandler(w http.ResponseWriter, r *http.Request) {
	client := &models.Client{}

	err := json.NewDecoder(r.Body).Decode(client)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := client.Store()
	u.Respond(w, res)
}

func ShowClientHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	client := &models.Client{}

	data := client.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateClientHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	client := &models.Client{}
	err = json.NewDecoder(r.Body).Decode(client)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := client.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyClientHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	client := &models.Client{}

	res := client.Destroy(uint(ID))
	u.Respond(w, res)
}
