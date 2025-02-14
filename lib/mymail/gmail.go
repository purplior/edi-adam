package mymail

import (
	"log"

	"github.com/purplior/sbec/domain/shared/logger"
	"gopkg.in/gomail.v2"
)

type (
	SendGmailRequest struct {
		To           string
		From         string
		FromPassword string
		Subject      string
		Body         string
	}
)

func SendGmail(request SendGmailRequest) error {
	message := gomail.NewMessage()
	message.SetHeader("From", request.From)
	message.SetHeader("To", request.To)
	message.SetHeader("Subject", request.Subject)
	message.SetBody("text/html", request.Body)

	logger.Debug("SendGmail()")
	logger.Debug("To: %s, From: %s", request.To, request.From)
	logger.Debug("Subject: %s", request.Subject)

	dialer := gomail.NewDialer(
		"smtp.gmail.com",
		587,
		request.From,
		request.FromPassword,
	)

	// 이메일 전송
	if err := dialer.DialAndSend(message); err != nil {
		return err
	}

	log.Printf("# Email Sent: %s", request.To)

	return nil
}
