package handlers

import (
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func IndexRecurringInvoiceHandler(w http.ResponseWriter, r *http.Request) {
	recurringInvoice := &models.InvoiceRecurring{}

	data := recurringInvoice.GenerateReccuringInvoice()
	u.Respond(w, data)
}
