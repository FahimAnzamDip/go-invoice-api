package models

import "gorm.io/gorm"

type InvoiceProduct struct {
	gorm.Model
	InvoiceID   uint     `gorm:"not null;" json:"invoice_id"`
	Invoice     *Invoice `gorm:"foreignKey:InvoiceID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
	ProductID   uint     `gorm:"" json:"product_id"`
	ProductName string   `gorm:"" json:"product_name"`
	Product     *Product `gorm:"foreignKey:ProductID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"-"`
	Quantity    int      `gorm:"not null;" json:"quantity"`
	UnitPrice   float32  `gorm:"not null;type:integer;" json:"unit_price"`
	TaxID       *uint     `gorm:"" json:"tax_id"`
	Tax         *Tax     `gorm:"foreignKey:TaxID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"tax,omitempty"`
	SubTotal    float32  `gorm:"not null;type:integer;" json:"sub_total"`
}

// BeforeCreate is called implicitly just before creating an entry
func (invoiceProduct *InvoiceProduct) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("unit_price", invoiceProduct.UnitPrice*100)
	tx.Statement.SetColumn("sub_total", invoiceProduct.SubTotal*100)
	return nil
}

// BeforeUpdate is called implicitly just before updating an entry
func (invoiceProduct *InvoiceProduct) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("unit_price", invoiceProduct.UnitPrice*100)
	tx.Statement.SetColumn("sub_total", invoiceProduct.SubTotal*100)
	return nil
}

// AfterFind is called implicitly just after finding an entry
func (invoiceProduct *InvoiceProduct) AfterFind(tx *gorm.DB) error {
	tx.Statement.SetColumn("unit_price", invoiceProduct.UnitPrice/100)
	tx.Statement.SetColumn("sub_total", invoiceProduct.SubTotal/100)
	return nil
}
