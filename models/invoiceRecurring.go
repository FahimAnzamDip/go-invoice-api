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

func (invoiceRecurring *InvoiceRecurring) GenerateReccuringInvoice() map[string]interface{} {
	invoices := []Invoice{}
	err := db.Where("recurring = ?", 1).Where("DATE(issue_date) <= ?", time.Now().Format("2006-01-02")).
		Find(&invoices).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = invoices

	return res
}
