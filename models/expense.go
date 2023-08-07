package models

import (
	"errors"
	"fmt"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Expense struct {
	gorm.Model
	Reference string   `gorm:"not null" json:"reference"`
	Date      string   `gorm:"not null" json:"date"`
	Name      string   `gorm:"not null" json:"name"`
	Amount    float32  `gorm:"not null;type:numeric(12,2);default:0;" json:"amount"`
	PurposeID *uint    `gorm:"" json:"purpose_id"`
	Purpose   *Purpose `gorm:"foreignKey:PurposeID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"purpose,omitempty"`
	Note      string   `gorm:"" json:"note"`
}

// BeforeCreate is called implicitly just before creating an entry
func (expense *Expense) BeforeCreate(tx *gorm.DB) error {
	var maxID *int
	tx.Model(&Expense{}).Select("MAX(id)").Scan(&maxID)

	var reference string
	if maxID == nil {
		reference = "#EXP-00001"
	} else {
		reference = fmt.Sprintf("#EXP-%05d", *maxID+1)
	}
	tx.Statement.SetColumn("reference", reference)
	return nil
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (expense *Expense) validate() (map[string]interface{}, bool) {
	if expense.Name == "" {
		return u.Message(false, "Expense name is required"), false
	}
	if expense.Date == "" {
		return u.Message(false, "Expense date is required"), false
	}
	if expense.Amount <= 0 {
		return u.Message(false, "Expense amount is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (expense *Expense) Index() map[string]interface{} {
	expenses := []Expense{}

	err := db.Preload("Purpose").Find(&expenses).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = expenses

	return res
}

// Store function creates a new entry
func (expense *Expense) Store() map[string]interface{} {
	if res, ok := expense.validate(); !ok {
		return res
	}

	err := db.Create(&expense).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Expense created successfully")
	res["data"] = expense

	return res
}

// Show function returns specific entry by ID
func (expense *Expense) Show(id uint) map[string]interface{} {
	err := db.Where("id = ?", id).First(&expense).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = expense

	return res
}

// Update function updates specific entry by ID
func (expense *Expense) Update(id uint) map[string]interface{} {
	_, err := expense.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&expense).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	exp, _ := expense.exists(id)

	res := u.Message(true, "Expense Updated Successfully")
	res["data"] = exp

	return res
}

// Destroy permanently removes a entry
func (expense *Expense) Destroy(id uint) map[string]interface{} {
	_, err := expense.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&expense).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Expense Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (expense *Expense) exists(id uint) (*Expense, error) {
	exp := &Expense{}
	err := db.Where("id = ?", id).Take(exp).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Expense{}, errors.New("no record found")
	}

	if err != nil {
		return &Expense{}, err
	}

	return exp, nil
}
