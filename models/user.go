package models

import (
	"errors"
	"os"
	"strings"
	"time"

	u "github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

// Token : JWT claims struct
type Token struct {
	UserInfo
	jwt.RegisteredClaims
}

type User struct {
	gorm.Model
	Name            string `gorm:"type:varchar(255);" json:"name"`
	Email           string `gorm:"type:varchar(255);unique;not null;" json:"email"`
	Mobile          string `gorm:"type:varchar(30);unique" json:"mobile"`
	Password        string `gorm:"type:varchar(255)" json:"password,omitempty"`
	ConfirmPassword string `gorm:"-" json:"password_confirmation,omitempty"`
	CurrentPassword string `gorm:"-" json:"current_password,omitempty"`
	Token           string `gorm:"-" json:"token,omitempty"`
	IsAdmin         bool   `gorm:"not null;default:0;" json:"is_admin"`
}

type UserInfo struct {
	UserId  uint
	Name    string
	Email   string
	Mobile  string
	IsAdmin bool
}

type AuthUser struct {
	EmailOrMobile string `json:"email_or_mobile"`
	Password      string `json:"password"`
}

type UserResponse struct {
	ID      uint   `json:"id"`
	Name    string `json:"name"`
	Email   string `json:"email"`
	Mobile  string `json:"mobile"`
	IsAdmin bool   `json:"is_admin"`
}

// Validate incoming user details...
func (user *User) validate(action string) (map[string]interface{}, bool) {
	switch action {
	case "create":
		if !strings.Contains(user.Email, "@") {
			return u.Message(false, "Email address is required"), false
		}

		if len(user.Password) < 6 {
			return u.Message(false, "Password is required"), false
		}

		//Email must be unique
		temp := &User{}

		//check for errors and duplicate emails
		err := db.Where("email = ?", user.Email).First(&temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return u.Message(false, "Connection error. Please retry"), false
		}
		if temp.Email != "" {
			return u.Message(false, "Email already exists"), false
		}
		//check for errors and duplicate emails
		err = db.Where("mobile = ?", user.Mobile).First(&temp).Error
		if err != nil && err != gorm.ErrRecordNotFound {
			return u.Message(false, "Connection error. Please retry"), false
		}
		if temp.Mobile != "" {
			return u.Message(false, "Mobile number already exists"), false
		}
	case "update_password":
		if user.CurrentPassword == "" {
			return u.Message(false, "Current Password is required"), false
		}

		if user.Password == "" {
			return u.Message(false, "Password is required"), false
		}

		temp := &User{}

		err := db.Where("id = ?", user.ID).First(&temp).Error
		if err != nil {
			if err == gorm.ErrRecordNotFound {
				return u.Message(false, "User not found"), false
			}
			return u.Message(false, "Connection error. Please retry"), false
		}

		err = bcrypt.CompareHashAndPassword([]byte(temp.Password), []byte(user.CurrentPassword))
		if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
			return u.Message(false, "Current password does not match"), false
		}
	case "confirm_password":
		if user.ConfirmPassword != user.Password {
			return u.Message(false, "Password and confirm password do not match"), false
		}
	}

	return u.Message(true, "Requirement passed"), true
}

// Index to retrive list of all users
func (user *User) Index() map[string]interface{} {
	users := []UserResponse{}
	db.Table("users").Where("is_admin = ?", false).Find(&users)

	res := u.Message(true, "")
	res["data"] = users
	return res
}

// Create is to insert a new entry in the DB
func (user *User) Store() map[string]interface{} {
	if res, ok := user.validate("create"); !ok {
		return res
	}
	if res, ok := user.validate("confirm_password"); !ok {
		return res
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	user.Password = string(hashedPassword)

	err := db.Create(&user).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	if user.ID == 0 {
		return u.Message(false, "Failed to create user, connection error.")
	}

	//Create new JWT token for the newly registered user
	tkn := &Token{
		UserInfo: UserInfo{
			UserId:  user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Mobile:  user.Mobile,
			IsAdmin: user.IsAdmin,
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tkn)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString

	user.Password = "" //remove password from response

	res := u.Message(true, "User has been created")
	res["data"] = user
	return res
}

// Login is to autheticate an user either by email or mobile
func Login(emailOrMobile, password string) map[string]interface{} {
	user := &User{}
	err := db.Where("email = ? OR mobile = ?", emailOrMobile, emailOrMobile).First(user).Error
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return u.Message(false, "Email/Mobile not found")
		}
		return u.Message(false, "Connection error. Please retry")
	}

	err = bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
	if err != nil && err == bcrypt.ErrMismatchedHashAndPassword { //Password does not match!
		return u.Message(false, "Invalid login credentials. Please try again")
	}
	// Worked! Logged In
	user.Password = ""

	// Create JWT token
	expirationTime := time.Now().Add(180 * time.Minute) //token expires after 3 hours
	tkn := &Token{
		UserInfo: UserInfo{
			UserId:  user.ID,
			Name:    user.Name,
			Email:   user.Email,
			Mobile:  user.Mobile,
			IsAdmin: user.IsAdmin,
		},
		RegisteredClaims: jwt.RegisteredClaims{
			ExpiresAt: jwt.NewNumericDate(expirationTime),
		},
	}
	token := jwt.NewWithClaims(jwt.GetSigningMethod("HS256"), tkn)
	tokenString, _ := token.SignedString([]byte(os.Getenv("token_password")))
	user.Token = tokenString //Store the token in the response

	res := u.Message(true, "Logged In")
	res["data"] = user
	return res
}

// Show finds an user by ID
func (user *User) Show(id uint) map[string]interface{} {
	usr, err := user.exists(id)

	if err != nil {
		return u.Message(false, err.Error())
	}

	data := &UserResponse{
		ID:      usr.ID,
		Name:    usr.Name,
		Email:   usr.Email,
		Mobile:  usr.Mobile,
		IsAdmin: usr.IsAdmin,
	}

	res := u.Message(true, "")
	res["data"] = data
	return res
}

// ValidateUser validates the user from context
func (user *User) ValidateUser(id uint) map[string]interface{} {
	oldUser, err := user.exists(id)

	if err != nil {
		return u.Message(false, err.Error())
	}

	oldUser.Password = ""
	oldUser.Token = user.Token

	res := u.Message(true, "")
	res["data"] = oldUser
	return res
}

// Update function updates specific entry by ID
func (user *User) Update(id uint) map[string]interface{} {
	_, err := user.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Omit("password", "is_admin").Updates(&user).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	usr, _ := user.exists(id)
	usr.Password = ""

	res := u.Message(true, "Profile updated")
	res["data"] = usr

	return res
}

// UpdatePassword updates the user's password
func (user *User) UpdatePassword() map[string]interface{} {
	if res, ok := user.validate("update_password"); !ok {
		return res
	}
	if res, ok := user.validate("confirm_password"); !ok {
		return res
	}

	hashedPassword, _ := bcrypt.GenerateFromPassword([]byte(user.Password), bcrypt.DefaultCost)
	err := db.Model(&user).Select("password").Update("password", hashedPassword).Error

	if err != nil {
		return u.Message(false, err.Error())
	}

	user.Password = ""
	user.CurrentPassword = ""
	user.ConfirmPassword = ""

	res := u.Message(true, "Password updated")
	res["data"] = user
	return res
}

// Destroy permanently removes a entry
func (user *User) Destroy(id uint) map[string]interface{} {
	_, err := user.exists(id)
	if err != nil {
		return u.Message(false, err.Error())
	}

	err = db.Where("id = ?", id).Unscoped().Delete(&user).Error
	if err != nil {
		return u.Message(false, err.Error())
	}

	res := u.Message(true, "User Deleted Successfully")
	res["data"] = "1"

	return res
}

// exists function checks if the entry exists or not
func (user *User) exists(id uint) (*User, error) {
	usr := &User{}
	err := db.Where("id = ?", id).Take(usr).Error

	if errors.Is(err, gorm.ErrRecordNotFound) {
		return &User{}, errors.New("no record found")
	}

	if err != nil {
		return &User{}, err
	}

	return usr, nil
}
