package models

import (
	"errors"
	"fmt"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type Client struct {
	gorm.Model
	Name       string `gorm:"-" json:"name,omitempty"`
	Email      string `gorm:"-" json:"email,omitempty"`
	Mobile     string `gorm:"-" json:"mobile,omitempty"`
	Password   string `gorm:"-" json:"password,omitempty"`
	UserID     uint   `gorm:"not null" json:"user_id"`
	User       *User  `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"user,omitempty"`
	Reference  string `gorm:"not null" json:"reference"`
	VatNumber  string `gorm:"" json:"vat_number"`
	Address    string `gorm:"type:text;" json:"address"`
	SendInvite bool   `gorm:"-" json:"send_invite,omitempty"`
}

// Validate validates the required parameters sent through the http request body
// returns message and true if the requirement is met
func (client *Client) validate() (map[string]interface{}, bool) {
	if client.Name == "" {
		return u.Message(false, "Client name is required"), false
	}
	if client.Email == "" {
		return u.Message(false, "Client email is required"), false
	}
	if client.Mobile == "" {
		return u.Message(false, "Client Mobile is required"), false
	}
	if len(client.Password) < 6 {
		return u.Message(false, "Account password is required"), false
	}

	//Email must be unique
	temp := &User{}

	//check for errors and duplicate emails
	err := db.Where("email = ?", client.Email).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Email != "" {
		return u.Message(false, "Email already exists"), false
	}
	//check for errors and duplicate emails
	err = db.Where("mobile = ?", client.Mobile).First(temp).Error
	if err != nil && err != gorm.ErrRecordNotFound {
		return u.Message(false, "Connection error. Please retry"), false
	}
	if temp.Mobile != "" {
		return u.Message(false, "Mobile number already exists"), false
	}
	// All the required parameters are present
	return u.Message(true, "success"), true
}

// BeforeCreate is called implicitly just before creating an entry
func (client *Client) BeforeCreate(tx *gorm.DB) error {
	var maxID *int
	tx.Model(&Client{}).Select("MAX(id)").Scan(&maxID)

	var reference string
	if maxID == nil {
		reference = "#CS-0001"
	} else {
		reference = fmt.Sprintf("#CS-%04d", *maxID+1)
	}
	tx.Statement.SetColumn("reference", reference)

	return nil
}

// Index function returns all entries
func (client *Client) Index() map[string]interface{} {
	clients := []Client{}
	err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email", "mobile")
	}).Find(&clients).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "")
	res["data"] = clients

	return res
}

// Store function creates a new entry
func (client *Client) Store() map[string]interface{} {
	if res, ok := client.validate(); !ok {
		return res
	}

	user := &User{
		Name:   client.Name,
		Email:  client.Email,
		Mobile: client.Mobile,
		Password: func() string {
			hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
			return string(hashedPassword)
		}(),
	}

	err := db.Create(user).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	client.UserID = user.ID

	err = db.Create(&client).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Client created successfully")
	res["data"] = client

	return res
}

// Show function returns specific entry by ID
func (client *Client) Show(id uint) map[string]interface{} {
	err := db.Preload("User", func(db *gorm.DB) *gorm.DB {
		return db.Select("id", "name", "email", "mobile")
	}).Where("id = ?", id).First(&client).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return u.Message(false, "Record not found")
	}

	res := u.Message(true, "")
	res["data"] = client

	return res
}

// Update function updates specific entry by ID
func (client *Client) Update(id uint) map[string]interface{} {
	clnt, err := client.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	user := &User{}
	db.Model(&clnt).Association("User").Find(user)

	user.Name = client.Name
	if user.Email != client.Email {
		user.Email = client.Email
	}
	if user.Mobile != client.Mobile {
		user.Mobile = client.Mobile
	}
	if client.Password != "" {
		hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(client.Password), bcrypt.DefaultCost)
		user.Password = string(hashedPassword)
	}
	err = db.Where("id = ?", user.ID).Updates(user).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Updates(&client).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	upClnt, err := client.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Client Updated Successfully")
	res["data"] = upClnt

	return res
}

// Destroy permanently removes a entry
func (client *Client) Destroy(id uint) map[string]interface{} {
	client, err := client.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	user := &User{}

	err = db.Where("id = ?", client.UserID).Unscoped().Delete(user).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "Client Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (client *Client) exists(id uint) (*Client, error) {
	clnt := &Client{}
	err := db.Where("id = ?", id).Take(clnt).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &Client{}, errors.New("no record found")
	}

	if err != nil {
		return &Client{}, err
	}

	return clnt, nil
}
