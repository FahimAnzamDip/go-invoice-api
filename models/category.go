package models

import (
	"errors"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Category struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (category *Category) validate() (map[string]interface{}, bool) {
	if category.Name == "" {
		return u.Message(false, "Category name is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (category *Category) Index() map[string]interface{} {
	categories := []Category{}

	err := db.Find(&categories).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = categories

	return res
}

// Store function creates a new entry
func (category *Category) Store() map[string]interface{} {
	if res, ok := category.validate(); !ok {
		return res
	}

	err := db.Create(&category).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Category created successfully")
	res["data"] = category

	return res
}

// Update function updates specific entry by ID
func (category *Category) Update(id uint) map[string]interface{} {
	_, err := category.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&category).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	cat, _ := category.exists(id)

	res := u.Message(true, "Category Updated Successfully")
	res["data"] = cat

	return res
}

// Destroy permanently removes a entry
func (category *Category) Destroy(id uint) map[string]interface{} {
	_, err := category.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&category).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Category Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (category *Category) exists(id uint) (*Category, error) {
	cat := &Category{}
	err := db.Where("id = ?", id).Take(cat).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Category{}, errors.New("no record found")
	}

	if err != nil {
		return &Category{}, err
	}

	return cat, nil
}
