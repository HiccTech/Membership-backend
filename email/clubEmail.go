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
	Subject       string
	DiscountCodes []DiscountCode
}

func SendClubEmail() {

	tmpl, err := template.ParseFiles("email/clubEmail.tmpl")
	if err != nil {
		panic(err)
	}

	data := EmailData{
		Subject: "通知：您已加入Hiccpet Club会员",
		DiscountCodes: []DiscountCode{
			{Title: "生日礼遇", Code: "BIRTHDAY2024", Period: "2025/09/03 - 2026/09/03"},
			{Title: "1V1宠物美容课程", Code: "GROOMING2024", Period: "2025/09/03 - 2026/09/03"},
			{Title: "新会员注册礼", Code: "WELCOME2024", Period: "2025/09/03 - 2026/09/03"},
			{Title: "宠物派对场地租赁8折优惠", Code: "PARTY20OFF", Period: "2025/09/03 - 2026/09/03"},
		},
	}

	var body bytes.Buffer
	if err := tmpl.Execute(&body, data); err != nil {
		panic(err)
	}

	from := "neo@hiccpet.com"
	password := "Lijian@2025"

	to := "812284688@qq.com"

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
