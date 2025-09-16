package email

import (
	"bytes"
	"fmt"
	"html/template"

	"gopkg.in/gomail.v2"
)

type DiscountCode struct {
	Title  string
	Code   string
	Period string
}

type EmailData struct {
	To                string
	Subject           string
	Template          string
	StoreCredit       StoreCredit
	DiscountCodes     []DiscountCode
	UsedDiscountCodes []struct {
		Title    string
		UsedDate string
	}
}

type StoreCredit struct {
	Amount    float64
	Currency  string
	ExpiredAt string
}

func SendClubEmail(emailData EmailData) {

	tmpl, err := template.ParseFiles(emailData.Template)
	if err != nil {
		panic(err)
	}

	data := emailData

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		panic(err)
	}

	from := "neo@hiccpet.com"
	password := "Lijian@2025"

	to := emailData.To

	m := gomail.NewMessage()
	m.SetHeader("From", from)
	m.SetHeader("To", to)
	m.SetHeader("Subject", data.Subject)
	m.SetBody("text/html", body.String())

	d := gomail.NewDialer("smtp.office365.com", 587, from, password)
	// STARTTLS is enabled by default for port 587 in gomail

	if err := d.DialAndSend(m); err != nil {
		fmt.Println("Failed to send:", err)
		return
	}

	fmt.Println("Email sent successfully!")
}
