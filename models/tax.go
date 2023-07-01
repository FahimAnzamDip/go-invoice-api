package models

import (
	"errors"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Tax struct {
	gorm.Model
	Name       string `gorm:"" json:"name"`
	Percentage int    `gorm:"" json:"percentage"`
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (tax *Tax) validate() (map[string]interface{}, bool) {
	if tax.Name == "" {
		return u.Message(false, "Tax name is required"), false
	}
	if tax.Percentage <= 0 {
		return u.Message(false, "Tax percentage is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (tax *Tax) Index() map[string]interface{} {
	taxes := []Tax{}

	err := db.Find(&taxes).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = taxes

	return res
}

// Store function creates a new entry
func (tax *Tax) Store() map[string]interface{} {
	if res, ok := tax.validate(); !ok {
		return res
	}

	err := db.Create(&tax).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Tax created successfully")
	res["data"] = tax

	return res
}

// Show function returns specific entry by ID
func (tax *Tax) Show(id uint) map[string]interface{} {
	err := db.Where("id = ?", id).First(&tax).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = tax

	return res
}

// Update function updates specific entry by ID
func (tax *Tax) Update(id uint) map[string]interface{} {
	_, err := tax.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&tax).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	tx, _ := tax.exists(id)

	res := u.Message(true, "Tax Updated Successfully")
	res["data"] = tx

	return res
}

// Destroy permanently removes a entry
func (tax *Tax) Destroy(id uint) map[string]interface{} {
	_, err := tax.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&tax).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Tax Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (tax *Tax) exists(id uint) (*Tax, error) {
	tx := &Tax{}
	err := db.Where("id = ?", id).Take(tx).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Tax{}, errors.New("no record found")
	}

	if err != nil {
		return &Tax{}, err
	}

	return tx, nil
}
