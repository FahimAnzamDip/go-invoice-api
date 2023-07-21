package models

import (
	"time"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type InvoiceRecurring struct {
	gorm.Model
	RecurringDate string   `gorm:"not null;type:date;" json:"recurring_date"`
	InvoiceID     uint     `gorm:"not null" json:"invoice_id"`
	Invoice       *Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"invoice,omitempty"`
	RefInvoiceID  uint     `gorm:"not null" json:"ref_invoice_id"`
	RefInvoice    *Invoice `gorm:"foreignKey:RefInvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"ref_invoice,omitempty"`
}

func (invoiceRecurring *InvoiceRecurring) GenerateReccuringInvoice() interface{} {
	invoices := []Invoice{}
	err := db.Where("recurring = ?", 1).Where("DATE(issue_date) <= ?", time.Now().Format("2006-01-02")).
		Preload("InvoiceProducts").Preload("Client.User").Find(&invoices).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	invoiceIDs := []uint{}

	for _, invoice := range invoices {
		recurring := &InvoiceRecurring{}
		err := db.Where("ref_invoice_id = ?", invoice.ID).Order("recurring_date desc").Limit(1).Find(recurring).Error
		if err != nil {
			return u.Message(false, err.Error())
		}

		var date time.Time
		if recurring.RecurringDate != "" {
			date, err = time.Parse("2006-01-02T15:04:05-07:00", recurring.RecurringDate)
			if err != nil {
				return u.Message(false, err.Error())
			}
		} else {
			date, err = time.Parse("2006-01-02T15:04:05-07:00", invoice.IssueDate)
			if err != nil {
				return u.Message(false, err.Error())
			}
		}

		var recurringDate time.Time
		if invoice.RecurringCycle == "Monthly" {
			recurringDate = date.AddDate(0, 1, 0)
		} else if invoice.RecurringCycle == "Quarterly" {
			recurringDate = date.AddDate(0, 3, 0)
		} else if invoice.RecurringCycle == "Semi Annually" {
			recurringDate = date.AddDate(0, 6, 0)
		} else if invoice.RecurringCycle == "Annually" {
			recurringDate = date.AddDate(1, 0, 0)
		} else {
			return u.Message(false, "Recurring cycle is not valid")
		}

		if time.Now().Format("2006-01-02") > recurringDate.Format("2006-01-02") {
			invoice.IssueDate = recurringDate.Format("2006-01-02")
			invoice.DueDate = recurringDate.AddDate(0, 0, 7).Format("2006-01-02")
			invoice.Reference = ""
			invoice.PaidAmount = 0
			invoice.DueAmount = invoice.TotalAmount
			invoice.Status = "Unpaid"
			invoice.Recurring = 3
			invoice.Terms = ""

			oldInvID := invoice.ID
			oldInvProds := invoice.InvoiceProducts
			invoice.ID = 0

			err = db.Omit("InvoiceProducts", "Payments").Create(&invoice).Error
			if err != nil {
				return u.Message(false, err.Error())
			}

			for _, invoiceProduct := range oldInvProds {
				invoiceProduct.InvoiceID = invoice.ID

				err := db.Omit("id").Create(&invoiceProduct).Error
				if err != nil {
					return u.Message(false, err.Error())
				}
			}

			invoiceRecurring = &InvoiceRecurring{
				RecurringDate: invoice.IssueDate,
				InvoiceID:     invoice.ID,
				RefInvoiceID:  oldInvID,
			}

			err = db.Create(invoiceRecurring).Error
			if err != nil {
				return u.Message(false, err.Error())
			}

			invoiceIDs = append(invoiceIDs, invoice.ID)
		}
	}

	return invoiceIDs
}
