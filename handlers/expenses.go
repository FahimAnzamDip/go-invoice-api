package handlers

import (
	"encoding/json"
	"net/http"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func IndexExpenseHandler(w http.ResponseWriter, r *http.Request) {
	expense := &models.Expense{}

	data := expense.Index()
	u.Respond(w, data)
}

func StoreExpenseHandler(w http.ResponseWriter, r *http.Request) {
	expense := &models.Expense{}

	err := json.NewDecoder(r.Body).Decode(expense)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := expense.Store()
	u.Respond(w, res)
}

func ShowExpenseHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	expense := &models.Expense{}

	res := expense.Show(uint(ID))
	u.Respond(w, res)
}

func UpdateExpenseHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	expense := &models.Expense{}
	err = json.NewDecoder(r.Body).Decode(expense)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := expense.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyExpenseHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseUint(id, 10, 32)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	expense := &models.Expense{}

	res := expense.Destroy(uint(ID))
	u.Respond(w, res)
}
