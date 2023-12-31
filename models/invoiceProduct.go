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
	UnitPrice   float32  `gorm:"not null;type:numeric(12,2);" json:"unit_price"`
	TaxID       *uint    `gorm:"" json:"tax_id"`
	Tax         *Tax     `gorm:"foreignKey:TaxID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"tax,omitempty"`
	SubTotal    float32  `gorm:"not null;type:numeric(12,2);" json:"sub_total"`
}
