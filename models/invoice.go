package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/fahimanzamdip/go-invoice-api/services"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Invoice struct {
	gorm.Model
	Reference       string            `gorm:"not null" json:"reference"`
	ClientID        uint              `gorm:"" json:"client_id"`
	Client          *Client           `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"client,omitempty"`
	Status          string            `gorm:"" json:"status"`
	IssueDate       string            `gorm:"type:date;" json:"issue_date"`
	DueDate         string            `gorm:"type:date;" json:"due_date"`
	Recurring       int8              `gorm:"not null;default:0;" json:"recurring"`
	RecurringCycle  string            `gorm:"" json:"recurring_cycle"`
	DiscountType    string            `gorm:"" json:"discount_type"`
	DiscountAmount  float32           `gorm:"type:numeric(12,2);not null;default:0;" json:"discount_amount"`
	TotalAmount     float32           `gorm:"type:numeric(12,2);not null;default:0;" json:"total_amount"`
	PaidAmount      float32           `gorm:"type:numeric(12,2);not null;default:0;" json:"paid_amount"`
	DueAmount       float32           `gorm:"type:numeric(12,2);not null;default:0;" json:"due_amount"`
	TaxAmount       float32           `gorm:"type:numeric(12,2);not null;default:0;" json:"tax_amount"`
	Terms           string            `gorm:"type:text;" json:"terms"`
	InvoiceProducts []*InvoiceProduct `gorm:"foreignKey:InvoiceID;" json:"invoice_products,omitempty"`
	PaymentMethod   string            `gorm:"-" json:"payment_method,omitempty"`
	Note            string            `gorm:"-" json:"note"`
	Payments        []*Payment        `gorm:"foreignKey:InvoiceID;" json:"payments,omitempty"`
}

// BeforeCreate is called implicitly just before creating an entry
func (invoice *Invoice) BeforeCreate(tx *gorm.DB) error {
	var maxID *int
	tx.Model(&Invoice{}).Select("MAX(id)").Scan(&maxID)

	var reference string
	if maxID == nil {
		reference = "#INV-00001"
	} else {
		reference = fmt.Sprintf("#INV-%05d", *maxID+1)
	}
	tx.Statement.SetColumn("reference", reference)
	return nil
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (invoice *Invoice) validate() (map[string]interface{}, bool) {
	if invoice.ClientID <= 0 {
		return u.Message(false, "Client is required"), false
	}
	if invoice.PaymentMethod == "" {
		return u.Message(false, "Payment Method is required"), false
	}
	if invoice.IssueDate == "" {
		return u.Message(false, "Issue date is required"), false
	}
	if invoice.PaidAmount < 0 {
		return u.Message(false, "Received amount is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (invoice *Invoice) Index(invType string) map[string]interface{} {
	invoices := []Invoice{}

	query := db.Model(&Invoice{})

	if invType == "" {
		query.Where("recurring != ?", 3)
	} else {
		query.Where("recurring = ?", 3)
	}

	err := query.Preload("Client.User").Find(&invoices).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = invoices

	return res
}

// Store function creates a new entry
func (invoice *Invoice) Store() map[string]interface{} {
	if res, ok := invoice.validate(); !ok {
		return res
	}

	if invoice.DueAmount == invoice.TotalAmount {
		invoice.Status = "Unpaid"
	} else if invoice.DueAmount == 0 {
		invoice.Status = "Paid"
	} else {
		invoice.Status = "Partially Paid"
	}

	err := db.Omit("InvoiceProducts", "Payments").Create(&invoice).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	for _, invoiceProduct := range invoice.InvoiceProducts {
		var productName string
		db.Model(&Product{}).Where("id = ?", invoiceProduct.ProductID).
			Select("name").Scan(&productName)
		invoiceProduct.ProductName = productName
		invoiceProduct.InvoiceID = invoice.ID

		err := db.Create(&invoiceProduct).Error
		if err != nil {
			return u.Message(false, err.Error())
		}
	}

	payment := &Payment{
		InvoiceID:     invoice.ID,
		ReceivedOn:    func() string { currentTime := time.Now().Format("2006-01-02"); return currentTime }(),
		Amount:        invoice.PaidAmount,
		PaymentMethod: invoice.PaymentMethod,
		Note:          invoice.Note,
	}

	err = db.Create(payment).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Invoice created successfully")
	res["data"] = invoice

	return res
}

// Show function returns specific entry by ID
func (invoice *Invoice) Show(id uint) map[string]interface{} {
	err := db.Preload("Client.User").Preload("InvoiceProducts.Tax").
		Preload("Payments").Where("id = ?", id).First(&invoice).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = invoice

	return res
}

// Update function updates specific entry by ID
func (invoice *Invoice) Update(id uint) map[string]interface{} {
	_, err := invoice.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	if invoice.DueAmount == invoice.TotalAmount {
		invoice.Status = "Unpaid"
	} else if invoice.DueAmount == 0 {
		invoice.Status = "Paid"
	} else {
		invoice.Status = "Partially Paid"
	}

	err = db.Where("id = ?", id).Omit("InvoiceProducts", "Payments").Updates(&invoice).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	inv, _ := invoice.exists(id)

	db.Unscoped().Model(&inv).Association("InvoiceProducts").Unscoped().Clear()

	for _, invoiceProduct := range invoice.InvoiceProducts {
		var productName string
		db.Model(&Product{}).Where("id = ?", invoiceProduct.ProductID).
			Select("name").Scan(&productName)
		invoiceProduct.ProductName = productName
		invoiceProduct.InvoiceID = inv.ID

		err := db.Create(&invoiceProduct).Error
		if err != nil {
			return u.Message(false, err.Error())
		}
	}

	res := u.Message(true, "Invoice Updated Successfully")
	res["data"] = inv

	return res
}

// Destroy permanently removes a entry
func (invoice *Invoice) Destroy(id uint) map[string]interface{} {
	_, err := invoice.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&invoice).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Invoice Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (invoice *Invoice) exists(id uint) (*Invoice, error) {
	inv := &Invoice{}
	err := db.Where("id = ?", id).Take(inv).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Invoice{}, errors.New("no record found")
	}

	if err != nil {
		return &Invoice{}, err
	}

	return inv, nil
}

func (invoice *Invoice) GeneratePDF() (string, error) {
	err := db.Preload("Client.User").Preload("InvoiceProducts.Tax").
		Preload("Payments").Where("id = ?", invoice.ID).First(&invoice).Error
	if err != nil {
		return "", err
	}

	pdfData := struct {
		Invoice     *Invoice
		Setting     *Setting
		CurrentYear string
	}{
		Invoice:     invoice,
		Setting:     (&Setting{}).AppGetSettings(),
		CurrentYear: time.Now().Format("2006"),
	}
	path, err := services.NewPDFService().GenerateInvoicePDF(pdfData)
	if err != nil {
		return "", err
	}

	return path, nil
}
