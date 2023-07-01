package models

import (
	"fmt"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Reference     string   `gorm:"not null" json:"reference"`
	InvoiceID     uint     `gorm:"not null" json:"invoice_id"`
	Invoice       *Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ReceivedOn    string   `gorm:"not null;type:date;" json:"received_on"`
	Amount        float32  `gorm:"type:integer;not null;default:0;" json:"amount"`
	PaymentMethod string   `gorm:"not null;" json:"payment_method"`
	Note          string   `gorm:"type:text;" json:"note"`
}

// BeforeCreate is called implicitly just before creating an entry
func (payment *Payment) BeforeCreate(tx *gorm.DB) error {
	var maxID *int
	tx.Model(&Payment{}).Select("MAX(id)").Scan(&maxID)

	var reference string
	if maxID == nil {
		reference = "#PAY-00001"
	} else {
		reference = fmt.Sprintf("#PAY-%05d", *maxID+1)
	}
	tx.Statement.SetColumn("reference", reference)
	tx.Statement.SetColumn("amount", payment.Amount*100)
	return nil
}

// BeforeUpdate is called implicitly just before updating an entry
func (payment *Payment) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("amount", payment.Amount*100)
	return nil
}

// AfterFind is called implicitly just after finding an entry
func (payment *Payment) AfterFind(tx *gorm.DB) error {
	tx.Statement.SetColumn("amount", payment.Amount/100)
	return nil
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (payment *Payment) validate() (map[string]interface{}, bool) {
	if payment.InvoiceID <= 0 {
		return u.Message(false, "Invoice is required"), false
	}
	if payment.ReceivedOn == "" {
		return u.Message(false, "Received on date is required"), false
	}
	if payment.Amount <= 0 {
		return u.Message(false, "Amount is required"), false
	}
	if payment.PaymentMethod == "" {
		return u.Message(false, "Payment Method is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (payment *Payment) Index() map[string]interface{} {
	payments := []Payment{}
	err := db.Find(&payments).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = payments

	return res
}

// Store function creates a new entry
func (payment *Payment) Store() map[string]interface{} {
	if res, ok := payment.validate(); !ok {
		return res
	}

	invoice := &Invoice{}

	err := db.Where("id = ?", payment.InvoiceID).Take(invoice).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	paid := invoice.PaidAmount + payment.Amount
	due := invoice.TotalAmount - paid
	invoice.PaidAmount = paid
	invoice.DueAmount = due
	if invoice.DueAmount == invoice.TotalAmount {
		invoice.Status = "Unpaid"
	} else if invoice.DueAmount == 0 {
		invoice.Status = "Paid"
	} else {
		invoice.Status = "Partially Paid"
	}

	err = db.Where("id = ?", invoice.ID).Updates(invoice).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Create(&payment).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Payment created successfully")
	res["data"] = payment

	return res
}
