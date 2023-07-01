package controllers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := &models.Tag{}

	data := tag.Index()
	u.Respond(w, data)
}

func StoreTagHandler(w http.ResponseWriter, r *http.Request) {
	tag := &models.Tag{}

	err := json.NewDecoder(r.Body).Decode(tag)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := tag.Store()
	u.Respond(w, res)
}

func ShowTagHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	tag := &models.Tag{}

	data := tag.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateTagHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	tag := &models.Tag{}
	err = json.NewDecoder(r.Body).Decode(tag)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := tag.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyTagHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	tag := &models.Tag{}

	res := tag.Destroy(uint(ID))
	u.Respond(w, res)
}
