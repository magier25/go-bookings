package main

import (
	"log"
	"time"

	"github.com/magier25/go-bookings/internal/models"
	mail "github.com/xhit/go-simple-mail"
)

func listenForMail() {
	go func() {
		for {
			msg := <-app.MailChan
			sendMessage(&msg)
		}
	}()
}

func sendMessage(m *models.MailData) {
	server := mail.NewSMTPClient()
	server.Host = "127.0.0.1"
	server.Port = 1025
	server.KeepAlive = false
	server.ConnectTimeout = 10 * time.Second
	server.SendTimeout = 10 * time.Second

	client, err := server.Connect()
	if err != nil {
		errorLog.Println(err)
	}

	email := mail.NewMSG()
	email.SetFrom(m.From).AddTo(m.To).SetSubject(m.Subject)
	email.SetBody(mail.TextHTML, m.Content)

	err = email.Send(client)
	if err != nil {
		log.Println(err)
	} else {
		log.Println("Email sent!")
	}
}
