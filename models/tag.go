package models

import (
	"errors"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Tag struct {
	gorm.Model
	Name     string     `gorm:"not null" json:"name"`
	Products []*Product `gorm:"many2many:product_tags;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"-"`
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (tag *Tag) validate() (map[string]interface{}, bool) {
	if tag.Name == "" {
		return u.Message(false, "Tag name is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (tag *Tag) Index() map[string]interface{} {
	tags := []Tag{}

	err := db.Find(&tags).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = tags

	return res
}

// Store function creates a new entry
func (tag *Tag) Store() map[string]interface{} {
	if res, ok := tag.validate(); !ok {
		return res
	}

	err := db.Create(&tag).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Tag created successfully")
	res["data"] = tag

	return res
}

// Show function returns specific entry by ID
func (tag *Tag) Show(id uint) map[string]interface{} {
	err := db.Where("id = ?", id).First(&tag).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = tag

	return res
}

// Update function updates specific entry by ID
func (tag *Tag) Update(id uint) map[string]interface{} {
	_, err := tag.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&tag).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	tg, _ := tag.exists(id)

	res := u.Message(true, "Tag Updated Successfully")
	res["data"] = tg

	return res
}

// Destroy permanently removes a entry
func (tag *Tag) Destroy(id uint) map[string]interface{} {
	_, err := tag.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&tag).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Tag Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (tag *Tag) exists(id uint) (*Tag, error) {
	tg := &Tag{}
	err := db.Where("id = ?", id).Take(tg).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Tag{}, errors.New("no record found")
	}

	if err != nil {
		return &Tag{}, err
	}

	return tg, nil
}
