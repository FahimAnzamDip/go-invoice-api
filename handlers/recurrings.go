package handlers

import (
	"log"
	"net/http"
	"sync"

	"github.com/fahimanzamdip/go-invoice-api/config"
	"github.com/fahimanzamdip/go-invoice-api/models"
	"github.com/fahimanzamdip/go-invoice-api/services"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func GenerateRecurringHandler(w http.ResponseWriter, r *http.Request) {
	recurringInvoice := &models.InvoiceRecurring{}

	invoiceIDs := recurringInvoice.GenerateReccuringInvoice()

	invoices := []models.Invoice{}
	config.GetDB().Where("id IN ?", invoiceIDs).Find(&invoices)

	for _, invoice := range invoices {
		attchChan := make(chan string)
		var wg sync.WaitGroup
		wg.Add(2)
		invoiceCopy := invoice

		go func() {
			defer wg.Done()
			// generate pdf and get the attachment
			attachment, err := invoiceCopy.GeneratePDF()
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
			err := services.NewMailService().SendEmail([]string{invoiceCopy.Client.User.Email}, "Invoice From GoInvoicer",
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
	}

	res := u.Message(true, "Generated & sent invoices to cleint")

	u.Respond(w, res)
}
