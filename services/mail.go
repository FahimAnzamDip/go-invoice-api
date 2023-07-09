package services

import (
	"bytes"
	"html/template"
	"log"
	"os"
	"strconv"

	"github.com/fahimanzamdip/go-invoice-api/utils"
	"github.com/joho/godotenv"
	"gopkg.in/gomail.v2"
)

type MailService struct {
	smtpHost string
	smtpPort int
	smtpUser string
	smtpPass string
}

func NewMailService() *MailService {
	if err := godotenv.Load(); err != nil {
		log.Println("File .env not found, reading configuration from ENV")
		return &MailService{}
	}

	portStr := os.Getenv("mail_port")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		log.Println(err.Error())
		return &MailService{}
	}

	return &MailService{
		smtpHost: os.Getenv("mail_host"),
		smtpPort: port,
		smtpUser: os.Getenv("mail_username"),
		smtpPass: os.Getenv("mail_password"),
	}
}

func (mail *MailService) SendEmail(from string, to []string, subject string, templateName string, attachmentPath string, data interface{}) error {
	var err error

	templ := template.New(templateName)
	templ.Funcs(template.FuncMap{
		"Price": utils.FormatPrice,
		"Date":  utils.FormatToDate,
	})
	templ, err = templ.ParseFiles("./templates/email/" + templateName)
	if err != nil {
		return err
	}
	var body bytes.Buffer
	if err = templ.Execute(&body, data); err != nil {
		return err
	}

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body.String())
	if attachmentPath != "" {
		m.Attach(attachmentPath)
	}

	d := gomail.NewDialer(mail.smtpHost, mail.smtpPort, mail.smtpUser, mail.smtpPass)

	if err := d.DialAndSend(m); err != nil {
		return err
	}

	return nil
}
