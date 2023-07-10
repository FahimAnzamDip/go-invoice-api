package handlers

import (
	"net/http"

	"github.com/fahimanzamdip/go-invoice-api/config"
	"github.com/fahimanzamdip/go-invoice-api/models"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
)

func PaymentReportHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	startDate := params.Get("start_date")
	endDate := params.Get("end_date")
	clientID := params.Get("client_id")
	paymentMethod := params.Get("payment_method")

	var payments []struct {
		ReceivedOn    string  `json:"received_on"`
		InvoiceRef    string  `json:"invoice_ref"`
		PaymentMethod string  `json:"payment_method"`
		ClientName    string  `json:"client_name"`
		Amount        float32 `json:"amount"`
	}

	query := config.GetDB().Model(&models.Payment{}).
		Select("payments.*, users.name AS client_name, invoices.reference AS invoice_ref").
		Joins("JOIN invoices ON invoices.id = payments.invoice_id").
		Joins("JOIN clients ON clients.id = invoices.client_id").
		Joins("JOIN users ON users.id = clients.user_id")
	if startDate != "" {
		query.Where("DATE(received_on) >= ?", startDate)
	}
	if endDate != "" {
		query.Where("DATE(received_on) <= ?", endDate)
	}
	if clientID != "" {
		query.Where("clients.id = ?", clientID)
	}
	if paymentMethod != "" {
		query.Where("payment_method = ?", paymentMethod)
	}
	err := query.Find(&payments).Error
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
	}

	res := u.Message(true, "")
	res["data"] = payments

	u.Respond(w, res)
}

func InvoiceReportHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	startDate := params.Get("start_date")
	endDate := params.Get("end_date")
	clientID := params.Get("client_id")
	status := params.Get("status")

	var invoices []struct {
		IssueDate   string  `json:"issue_date"`
		Reference   string  `json:"reference"`
		ClientID    uint    `json:"client_id"`
		ClientName  string  `json:"client_name"`
		TotalAmount float32 `json:"total_amount"`
		PaidAmount  float32 `json:"paid_amount"`
		DueAmount   float32 `json:"due_amount"`
		Status      string  `json:"status"`
	}

	query := config.GetDB().Model(&models.Invoice{}).
		Select("invoices.*, users.name AS client_name").
		Joins("JOIN clients ON clients.id = invoices.client_id").
		Joins("JOIN users ON users.id = clients.user_id")
	if startDate != "" {
		query.Where("DATE(issue_date) >= ?", startDate)
	}
	if endDate != "" {
		query.Where("DATE(issue_date) <= ?", endDate)
	}
	if clientID != "" {
		query.Where("client_id = ?", clientID)
	}
	if status != "" {
		query.Where("status = ?", status)
	}
	err := query.Find(&invoices).Error
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
	}

	res := u.Message(true, "")
	res["data"] = invoices

	u.Respond(w, res)
}

func ExpenseReportHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	startDate := params.Get("start_date")
	endDate := params.Get("end_date")
	purposeID := params.Get("purpose_id")

	var expenses []struct {
		Date        string  `json:"date"`
		Reference   string  `json:"reference"`
		PurposeID   uint    `json:"purpose_id"`
		PurposeName string  `json:"purpose_name"`
		Amount      float32 `json:"amount"`
	}

	query := config.GetDB().Model(&models.Expense{}).
		Select("expenses.*, purposes.name AS purpose_name").
		Joins("JOIN purposes ON purposes.id = expenses.purpose_id")
	if startDate != "" {
		query.Where("DATE(date) >= ?", startDate)
	}
	if endDate != "" {
		query.Where("DATE(date) <= ?", endDate)
	}
	if purposeID != "" {
		query.Where("purpose_id = ?", purposeID)
	}
	err := query.Find(&expenses).Error
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
	}

	res := u.Message(true, "")
	res["data"] = expenses

	u.Respond(w, res)
}

func ClientStatementReportHandler(w http.ResponseWriter, r *http.Request) {
	params := r.URL.Query()
	startDate := params.Get("start_date")
	endDate := params.Get("end_date")
	clientID := params.Get("clientID")

	var invoices []struct {
		IssueDate  string  `json:"issue_date"`
		Reference  string  `json:"reference"`
		ClientID   uint    `json:"client_id"`
		ClientName string  `json:"client_name"`
		PaidAmount float32 `json:"payments"`
		DueAmount  float32 `json:"balance"`
		Activity   string  `json:"activity"`
	}

	query := config.GetDB().Model(&models.Invoice{}).
		Select(`invoices.*, users.name AS client_name, 
			CASE WHEN invoices.due_amount > 0 THEN 'Payment Received' ELSE 'Invoice Generated' END AS activity`).
		Joins("JOIN clients ON clients.id = invoices.client_id").
		Joins("JOIN users ON users.id = clients.user_id")
	if startDate != "" {
		query.Where("DATE(issue_date) >= ?", startDate)
	}
	if endDate != "" {
		query.Where("DATE(issue_date) <= ?", endDate)
	}
	if clientID != "" {
		query.Where("ClientID = ?", clientID)
	}
	err := query.Find(&invoices).Error
	if err != nil {
		u.Respond(w, u.Message(false, err.Error()))
	}

	res := u.Message(true, "")
	res["data"] = invoices

	u.Respond(w, res)
}
