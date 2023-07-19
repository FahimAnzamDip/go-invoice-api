package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"sync"

	"github.com/fahimanzamdip/go-invoice-api/config"
	"github.com/fahimanzamdip/go-invoice-api/models"
	"github.com/fahimanzamdip/go-invoice-api/services"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func SendEstimateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	estimate := &models.Estimate{}
	err = config.GetDB().Where("id = ?", ID).Preload("Client.User").Preload("EstimateProducts.Tax").
		First(estimate).Error
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
		return
	}

	attchChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// generate pdf and get the attachment
		attachment, err := estimate.GeneratePDF()
		if err != nil {
			log.Println(err.Error())
			u.Respond(w, u.Message(false, err.Error()))
			return
		}
		attchChan <- attachment
	}()

	go func() {
		defer wg.Done()
		attachment := <-attchChan
		// send email to client with attachment
		err = services.NewMailService().SendEmail([]string{estimate.Client.User.Email}, "Estimate From GoInvoicer",
			"estimate-mail.html",
			attachment, "")
		if err != nil {
			log.Println(err.Error())
			u.Respond(w, u.Message(false, "Estimate created. But can not send email!"))
			return
		} else {
			u.RemoveFile(attachment)
			err = config.GetDB().Model(&estimate).Where("id = ?", ID).Update("status", "Sent").Error
			if err != nil {
				u.Respond(w, u.Message(false, err.Error()))
				return
			}
		}
	}()

	wg.Wait()

	u.Respond(w, u.Message(true, "Estimate sent to client successfully"))
}

func IndexEstimateHandler(w http.ResponseWriter, r *http.Request) {
	estimate := &models.Estimate{}

	data := estimate.Index()
	u.Respond(w, data)
}

func StoreEstimateHandler(w http.ResponseWriter, r *http.Request) {
	estimate := &models.Estimate{}
	estimate.Status = "Pending"

	err := json.NewDecoder(r.Body).Decode(estimate)
	if err != nil {
		log.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := estimate.Store()

	u.Respond(w, res)
}

func ShowEstimateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	estimate := &models.Estimate{}

	data := estimate.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateEstimateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	estimate := &models.Estimate{}
	err = json.NewDecoder(r.Body).Decode(estimate)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := estimate.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyEstimateHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	estimate := &models.Estimate{}

	data := estimate.Destroy(uint(ID))
	u.Respond(w, data)
}
