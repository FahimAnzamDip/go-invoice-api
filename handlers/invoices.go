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

func RecurringCyclesHandler(w http.ResponseWriter, r *http.Request) {
	filePath := "./data/recurring_cycles.json"
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

func IndexInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	invType := params.Get("invoice_type")

	invoice := &models.Invoice{}

	data := invoice.Index(invType)
	u.Respond(w, data)
}

func StoreInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	invoice := &models.Invoice{}

	err := json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		log.Println(err)
		u.Respond(w, u.Message(false, "Error while decoding request body"))
		return
	}

	res := invoice.Store()

	attchChan := make(chan string)
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		// generate pdf and get the attachment
		attachment, err := invoice.GeneratePDF()
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
		err = services.NewMailService().SendEmail([]string{invoice.Client.User.Email}, "Invoice From GoInvoicer",
			"invoice-mail.html",
			attachment, "")
		if err != nil {
			log.Println(err.Error())
			u.Respond(w, u.Message(false, "Invoice created. But can not send email!"))
			return
		} else {
			u.RemoveFile(attachment)
		}
	}()

	wg.Wait()

	u.Respond(w, res)
}

func ShowInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	invoice := &models.Invoice{}

	data := invoice.Show(uint(ID))
	u.Respond(w, data)
}

func UpdateInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	invoice := &models.Invoice{}
	err = json.NewDecoder(r.Body).Decode(invoice)
	if err != nil {
		u.Respond(w, u.Message(false, "Invalid Request"))
		return
	}

	res := invoice.Update(uint(ID))
	u.Respond(w, res)
}

func DestroyInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	id := chi.URLParam(r, "id")
	ID, err := strconv.ParseInt(id, 10, 64)
	if err != nil {
		u.Respond(w, u.Message(false, "Error while parsing ID"))
		return
	}

	invoice := &models.Invoice{}

	data := invoice.Destroy(uint(ID))
	u.Respond(w, data)
}
