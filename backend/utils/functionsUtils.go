package utils

import (
	"encoding/base64"
	"fmt"
	"log"
	"olhar-biblico/config"
	"regexp"

	"gopkg.in/gomail.v2"
)

func SendEmail(name, email, phone, address, subject string) error {
	decodedBytes, err := base64.StdEncoding.DecodeString(config.EmailPassword)
	if err != nil {
		log.Fatal("Erro ao decodificar:", err.Error())
	}

	host := "smtp.hostinger.com"
	port := 587
	user := ""
	password := string(decodedBytes)

	emailsToSend := []string{""}
	if config.Production {
		emailsToSend = append(emailsToSend, "")
	}

	msg := gomail.NewMessage()
	msg.SetHeader("From", user)
	msg.SetHeader("To", emailsToSend...)
	msg.SetHeader("Subject", subject)

	text := fmt.Sprintf(`
		<h1>Oi!</h1>
		<br/><h3>%s!</h3>
		<br/>

		<p>Dados:<br/>
			Nome: %s<br/>
			Email: %s<br/>
			Telefone: %s<br/>
			Endereço: %s</p>

		<h5>Atenciosamente.</h5>
	`,
		subject,
		name,
		email,
		phone,
		address,
	)

	msg.SetBody("text/html", text)

	dialer := gomail.NewDialer(host, port, user, password)
	if err := dialer.DialAndSend(msg); err != nil {
		return err
	}

	return nil
}

func VerifyPhoneNumber(phoneNumber string) (bool, error) {
	r, err := regexp.Compile("\\([0-9]{2,3}\\) [0-9]{4,5}-[0-9]{4}")
	if err != nil {
		return false, err
	}

	return r.MatchString(phoneNumber), nil
}
