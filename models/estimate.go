package models

import (
	"errors"
	"fmt"
	"time"

	"github.com/fahimanzamdip/go-invoice-api/services"
	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"gorm.io/gorm"
)

type Estimate struct {
	gorm.Model
	Reference        string             `gorm:"not null" json:"reference"`
	ClientID         uint               `gorm:"" json:"client_id"`
	Client           *Client            `gorm:"foreignKey:ClientID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"client,omitempty"`
	Status           string             `gorm:"" json:"status"`
	Date             string             `gorm:"type:date;" json:"date"`
	DiscountType     string             `gorm:"" json:"discount_type"`
	DiscountAmount   float32            `gorm:"type:numeric(12,2);not null;default:0;" json:"discount_amount"`
	TotalAmount      float32            `gorm:"type:numeric(12,2);not null;default:0;" json:"total_amount"`
	TaxAmount        float32            `gorm:"type:numeric(12,2);not null;default:0;" json:"tax_amount"`
	Terms            string             `gorm:"type:text;" json:"terms"`
	EstimateProducts []*EstimateProduct `gorm:"foreignKey:EstimateID;" json:"estimate_products,omitempty"`
	Note             string             `gorm:"" json:"note"`
}

// BeforeCreate is called implicitly just before creating an entry
func (estimate *Estimate) BeforeCreate(tx *gorm.DB) error {
	var maxID *int
	tx.Model(&Estimate{}).Select("MAX(id)").Scan(&maxID)

	var reference string
	if maxID == nil {
		reference = "#EST-00001"
	} else {
		reference = fmt.Sprintf("#EST-%05d", *maxID+1)
	}
	tx.Statement.SetColumn("reference", reference)
	return nil
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (estimate *Estimate) validate() (map[string]interface{}, bool) {
	if estimate.ClientID <= 0 {
		return u.Message(false, "Client is required"), false
	}
	if estimate.Date == "" {
		return u.Message(false, "Date is required"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// Index function returns all entries
func (estimate *Estimate) Index() map[string]interface{} {
	estimates := []Estimate{}
	err := db.Preload("Client.User").Find(&estimates).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = estimates

	return res
}

// Store function creates a new entry
func (estimate *Estimate) Store() map[string]interface{} {
	if res, ok := estimate.validate(); !ok {
		return res
	}

	err := db.Omit("EstimateProducts").Create(&estimate).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	for _, estimateProduct := range estimate.EstimateProducts {
		var productName string
		db.Model(&Product{}).Where("id = ?", estimateProduct.ProductID).
			Select("name").Scan(&productName)
		estimateProduct.ProductName = productName
		estimateProduct.EstimateID = estimate.ID

		err := db.Create(&estimateProduct).Error
		if err != nil {
			return u.Message(false, err.Error())
		}
	}

	res := u.Message(true, "Estimate created successfully")
	res["data"] = estimate

	return res
}

// Show function returns specific entry by ID
func (estimate *Estimate) Show(id uint) map[string]interface{} {
	err := db.Preload("Client.User").Preload("EstimateProducts.Tax").
		Where("id = ?", id).First(&estimate).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = estimate

	return res
}

// Update function updates specific entry by ID
func (estimate *Estimate) Update(id uint) map[string]interface{} {
	_, err := estimate.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Omit("EstimateProducts").Updates(&estimate).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	est, _ := estimate.exists(id)

	db.Unscoped().Model(&est).Association("EstimateProducts").Unscoped().Clear()

	for _, estimateProduct := range estimate.EstimateProducts {
		var productName string
		db.Model(&Product{}).Where("id = ?", estimateProduct.ProductID).
			Select("name").Scan(&productName)
		estimateProduct.ProductName = productName
		estimateProduct.EstimateID = est.ID

		err := db.Create(&estimateProduct).Error
		if err != nil {
			return u.Message(false, err.Error())
		}
	}

	res := u.Message(true, "Estimate updated successfully")
	res["data"] = est

	return res
}

// Destroy permanently removes a entry
func (estimate *Estimate) Destroy(id uint) map[string]interface{} {
	_, err := estimate.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&estimate).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Estimate Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (estimate *Estimate) exists(id uint) (*Estimate, error) {
	est := &Estimate{}
	err := db.Where("id = ?", id).Take(est).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Estimate{}, errors.New("no record found")
	}

	if err != nil {
		return &Estimate{}, err
	}

	return est, nil
}

func (estimate *Estimate) GeneratePDF() (string, error) {
	err := db.Preload("Client.User").Preload("EstimateProducts.Tax").
		Where("id = ?", estimate.ID).First(&estimate).Error
	if err != nil {
		return "", err
	}

	pdfData := struct {
		Estimate    *Estimate
		Setting     *Setting
		CurrentYear string
	}{
		Estimate:    estimate,
		Setting:     (&Setting{}).AppGetSettings(),
		CurrentYear: time.Now().Format("2006"),
	}
	path, err := services.NewPDFService().GenerateEstimatePDF(pdfData)
	if err != nil {
		return "", err
	}

	return path, nil
}
