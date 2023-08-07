package handlers

import (
	"encoding/json"
	"log"
	"net/http"
	"os"
	"strconv"
	"sync"

	"github.com/fahimanzamdip/go-invoice-api/models"
	"github.com/fahimanzamdip/go-invoice-api/services"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/go-chi/chi/v5"
)

func PaymentMethodsHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./data/payment_methods.json"
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

func IndexPaymentHandler(w http.ResponseWriter, r *http.Request) {
	payment := &models.Payment{}

	data := payment.Index()
	u.Respond(w, data)
}

func StorePaymentHandler(w http.ResponseWriter, r *http.Request) {
	payment := &models.Payment{}

	err := json.NewDecoder(r.Body).Decode(payment)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := payment.Store()

	attchChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// generate pdf and get the attachment
		attachment, err := payment.GeneratePDF()
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
		// send email to client with the attachment
		err = services.NewMailService().
			SendEmail([]string{payment.Invoice.Client.User.Email}, "Payment Received. GoInvoicer",
				"payment-mail.html",
				attachment,
				struct {
					Reference string
					Amount    float32
				}{
					Reference: payment.Reference,
					Amount:    payment.Amount,
				})
		if err != nil {
			log.Println(err.Error())
			u.Respond(w, u.Message(false, "Payment created. But can not send email!"))
			return
		} else {
			u.RemoveFile(attachment)
		}
	}()

	wg.Wait()

	u.Respond(w, res)
}

func ShowPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	payment := &models.Payment{}

	res := payment.Show(uint(ID))
	u.Respond(w, res)
}

func UpdatePaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	payment := &models.Payment{}
	err = json.NewDecoder(r.Body).Decode(payment)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := payment.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyPaymentHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	payment := &models.Payment{}

	res := payment.Destroy(uint(ID))
	u.Respond(w, res)
}
