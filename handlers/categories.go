package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := &models.Category{}

	data := category.Index()
	u.Respond(w, data)
}

func StoreCategoryHandler(w http.ResponseWriter, r *http.Request) {
	category := &models.Category{}

	err := json.NewDecoder(r.Body).Decode(category)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := category.Store()
	u.Respond(w, res)
}

func ShowCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	category := &models.Category{}

	data := category.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	category := &models.Category{}
	err = json.NewDecoder(r.Body).Decode(category)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := category.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyCategoryHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	category := &models.Category{}

	res := category.Destroy(uint(ID))
	u.Respond(w, res)
}
