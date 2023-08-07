package models

import (
	"errors"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Purpose struct {
	gorm.Model
	Name        string `gorm:"not null" json:"name"`
	Description string `json:"description"`
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (purpose *Purpose) validate() (map[string]interface{}, bool) {
	if purpose.Name == "" {
		return u.Message(false, "Purpose name is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (purpose *Purpose) Index() map[string]interface{} {
	purposes := []Purpose{}

	err := db.Find(&purposes).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = purposes

	return res
}

// Store function creates a new entry
func (purpose *Purpose) Store() map[string]interface{} {
	if res, ok := purpose.validate(); !ok {
		return res
	}

	err := db.Create(&purpose).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Purpose created successfully")
	res["data"] = purpose

	return res
}

// Show function returns specific entry by ID
func (purpose *Purpose) Show(id uint) map[string]interface{} {
	err := db.Where("id = ?", id).First(&purpose).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = purpose

	return res
}

// Update function updates specific entry by ID
func (purpose *Purpose) Update(id uint) map[string]interface{} {
	_, err := purpose.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&purpose).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	prs, _ := purpose.exists(id)

	res := u.Message(true, "Purpose Updated Successfully")
	res["data"] = prs

	return res
}

// Destroy permanently removes a entry
func (purpose *Purpose) Destroy(id uint) map[string]interface{} {
	_, err := purpose.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&purpose).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Purpose Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (purpose *Purpose) exists(id uint) (*Purpose, error) {
	prs := &Purpose{}
	err := db.Where("id = ?", id).Take(prs).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Purpose{}, errors.New("no record found")
	}

	if err != nil {
		return &Purpose{}, err
	}

	return prs, nil
}
