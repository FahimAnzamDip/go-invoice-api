package models

import (
	"errors"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Product struct {
	gorm.Model
	Name        string    `gorm:"not null;" json:"name"`
	Code        string    `gorm:"" json:"code"`
	Description string    `json:"description"`
	Price       float32   `gorm:"type:integer;not null;" json:"price"`
	Stock       int       `gorm:"not null;default:0" json:"stock"`
	CategoryID  *uint     `json:"category_id"`
	Category    *Category `gorm:"foreignKey:CategoryID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"category"`
	ImagePath   string    `json:"image_url"`
	Tags        []*Tag    `gorm:"many2many:product_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"tags"`
	TagIDs      []uint    `gorm:"-" json:"tag_ids,omitempty"`
}

// BeforeCreate is called implicitly just before creating an entry
func (product *Product) BeforeCreate(tx *gorm.DB) error {
	tx.Statement.SetColumn("price", product.Price*100)
	return nil
}

// BeforeUpdate is called implicitly just before updating an entry
func (product *Product) BeforeUpdate(tx *gorm.DB) error {
	tx.Statement.SetColumn("price", product.Price*100)
	return nil
}

// AfterFind is called implicitly just after finding an entry
func (product *Product) AfterFind(tx *gorm.DB) error {
	tx.Statement.SetColumn("price", product.Price/100)
	return nil
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (product *Product) validate() (map[string]interface{}, bool) {
	if product.Name == "" {
		return u.Message(false, "Product name is required"), false
	}
	if product.Price <= 0 {
		return u.Message(false, "Product price is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (product *Product) Index() map[string]interface{} {
	products := []Product{}
	err := db.Preload("Category").Preload("Tags").Find(&products).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = products

	return res
}

// Store function creates a new entry
func (product *Product) Store() map[string]interface{} {
	if res, ok := product.validate(); !ok {
		return res
	}

	err := db.Create(&product).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	for _, tagID := range product.TagIDs {
		tag := &Tag{}
		tag.ID = tagID
		db.Model(&product).Association("Tags").Append(tag)
	}

	res := u.Message(true, "Product created successfully")
	res["data"] = product

	return res
}

// Show function returns specific entry by ID
func (product *Product) Show(id uint) map[string]interface{} {
	err := db.Preload("Category").Preload("Tags").Where("id = ?", id).First(&product).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = product

	return res
}

// Update function updates specific entry by ID
func (product *Product) Update(id uint) map[string]interface{} {
	_, err := product.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&product).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	prod, _ := product.exists(id)

	db.Unscoped().Model(&prod).Association("Tags").Unscoped().Clear()

	for _, tagID := range product.TagIDs {
		tag := &Tag{}
		tag.ID = tagID
		db.Model(&prod).Association("Tags").Append(tag)
	}

	res := u.Message(true, "Product Updated Successfully")
	res["data"] = prod

	return res
}

// Destroy permanently removes a entry
func (product *Product) Destroy(id uint) map[string]interface{} {
	_, err := product.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&product).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Product Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (product *Product) exists(id uint) (*Product, error) {
	prod := &Product{}
	err := db.Where("id = ?", id).Take(prod).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Product{}, errors.New("no record found")
	}

	if err != nil {
		return &Product{}, err
	}

	return prod, nil
}
