package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/fahimanzamdip/go-invoice-api/services"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Payment struct {
	gorm.Model
	Reference     string   `gorm:"not null" json:"reference"`
	InvoiceID     uint     `gorm:"not null" json:"invoice_id"`
	Invoice       *Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"invoice,omitempty"`
	ReceivedOn    string   `gorm:"not null;type:date;" json:"received_on"`
	Amount        float32  `gorm:"type:numeric(12,2);not null;default:0;" json:"amount"`
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

	_, err := adjustInvoice(invoice, payment, "create")
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

// Update function updates specific entry by ID
func (payment *Payment) Update(id uint) map[string]interface{} {
	pay, err := payment.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	// Revert back to previous amounts
	currInv := &Invoice{}
	err = db.Where("id = ?", pay.InvoiceID).Take(currInv).Error
	if err != nil {
		return u.Message(false, err.Error())
	}
	currInv.PaidAmount -= pay.Amount
	currInv.DueAmount = currInv.TotalAmount - currInv.PaidAmount
	// Adjust amounts again
	_, err = adjustInvoice(currInv, payment, "update")
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&payment).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	newPay, _ := payment.exists(id)

	res := u.Message(true, "Payment Updated Successfully")
	res["data"] = newPay

	return res
}

// Destroy permanently removes a entry
func (payment *Payment) Destroy(id uint) map[string]interface{} {
	payment, err := payment.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	invoice := &Invoice{}

	_, err = adjustInvoice(invoice, payment, "destroy")
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&payment).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Payment Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (payment *Payment) exists(id uint) (*Payment, error) {
	pay := &Payment{}
	err := db.Where("id = ?", id).Take(pay).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Payment{}, errors.New("no record found")
	}

	if err != nil {
		return &Payment{}, err
	}

	return pay, nil
}

func adjustInvoice(inv *Invoice, pay *Payment, action string) (*Invoice, error) {
	switch action {
	case "create":
		err := db.Where("id = ?", pay.InvoiceID).Take(inv).Error
		if err != nil {
			return &Invoice{}, err
		}
		paid := inv.PaidAmount + pay.Amount
		due := inv.TotalAmount - paid
		inv.PaidAmount = paid
		inv.DueAmount = due
	case "update":
		paid := inv.PaidAmount + pay.Amount
		due := inv.TotalAmount - paid
		inv.PaidAmount = paid
		inv.DueAmount = due
	case "destroy":
		err := db.Where("id = ?", pay.InvoiceID).Take(inv).Error
		if err != nil {
			return &Invoice{}, err
		}
		paid := inv.PaidAmount - pay.Amount
		due := inv.TotalAmount - paid
		inv.PaidAmount = paid
		inv.DueAmount = due
	}

	if inv.DueAmount == inv.TotalAmount {
		inv.Status = "Unpaid"
	} else if inv.DueAmount == 0 {
		inv.Status = "Paid"
	} else {
		inv.Status = "Partially Paid"
	}

	err := db.Where("id = ?", inv.ID).Updates(inv).Error
	if err != nil {
		return &Invoice{}, err
	}

	return inv, nil
}

func (payment *Payment) GeneratePDF() (string, error) {
	err := db.Preload("Invoice.Client.User").Where("id = ?", payment.ID).First(&payment).Error
	if err != nil {
		return "", err
	}

	pdfData := struct {
		Payment     *Payment
		Setting     *Setting
		CurrentYear string
	}{
		Payment:     payment,
		Setting:     (&Setting{}).AppGetSettings(),
		CurrentYear: time.Now().Format("2006"),
	}
	path, err := services.NewPDFService().GeneratePayReceiptPDF(pdfData)
	if err != nil {
		return "", err
	}

	return path, nil
}
