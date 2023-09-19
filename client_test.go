package mailosaur

import (
	"encoding/base64"
	"fmt"
	"net/smtp"
	"os"
	"strings"
)

var client *MailosaurClient
var server string
var verifiedDomain string
var emails []*MessageSummary
var email *Message

func sendEmails(client *MailosaurClient, server string, quantity int) {
	for i := 0; i < quantity; i++ {
		sendEmail(client, server, "")
	}
}

func sendEmail(client *MailosaurClient, server string, sendToAddress string) error {
	host := os.Getenv("MAILOSAUR_SMTP_HOST")
	port := os.Getenv("MAILOSAUR_SMTP_PORT")

	if len(host) == 0 {
		host = "mailosaur.net"
	}

	if len(port) == 0 {
		port = "25"
	}

	randomString := getRandomString()

	sendFrom := fmt.Sprintf("%s %s <%s@%s>", randomString, randomString, randomString, verifiedDomain)
	toAddress := client.Servers.GenerateEmailAddress(server)

	if len(sendToAddress) != 0 {
		toAddress = sendToAddress
	}

	sendTo := fmt.Sprintf("%s %s <%s>", randomString, randomString, toAddress)
	subject := randomString + " subject"

	r := strings.NewReplacer("REPLACED_DURING_TEST", randomString)

	delimeter := "--==_mimepart_" + randomString

	catImage, _ := os.ReadFile("testing/cat.png")
	dogImage, _ := os.ReadFile("testing/dog.png")

	htmlFile, _ := os.ReadFile("testing/testEmail.html")
	htmlContent := r.Replace(string(htmlFile))

	textFile, _ := os.ReadFile("testing/testEmail.txt")
	textContent := r.Replace(string(textFile))

	c, err := smtp.Dial(host + ":" + port)
	if err != nil {
		return err
	}
	defer c.Close()
	if err = c.Mail(sendFrom); err != nil {
		return err
	}

	if err = c.Rcpt(toAddress); err != nil {
		return err
	}

	w, err := c.Data()
	if err != nil {
		return err
	}

	msg := "MIME-Version: 1.0\r\n"

	msg += fmt.Sprintf("From: %s\r\n", sendFrom)
	msg += fmt.Sprintf("To: %s\r\n", sendTo)
	msg += fmt.Sprintf("Subject: %s\r\n", subject)

	msg += fmt.Sprintf("Content-Type: multipart/alternative; boundary=\"%s\"\r\n", delimeter)

	msg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	msg += "Content-Type: text/plain; charset=\"utf-8\"\r\n"
	msg += "Content-Transfer-Encoding: base64\r\n"
	msg += "\r\n" + base64.StdEncoding.EncodeToString([]byte(textContent)) + "\r\n"

	msg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	msg += "Content-Type: text/html; charset=\"utf-8\"\r\n"
	msg += "Content-Transfer-Encoding: base64\r\n"
	msg += "\r\n" + base64.StdEncoding.EncodeToString([]byte(htmlContent)) + "\r\n"

	msg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	msg += "Content-Type: image/png; filename=cat.png\r\n"
	msg += "Content-Transfer-Encoding: base64\r\n"
	msg += "Content-Disposition: attachment; filename=cat.png\r\n"
	msg += "content-id: ii_1435fadb31d523f6\r\n"

	msg += "\r\n" + base64.StdEncoding.EncodeToString(catImage)

	msg += fmt.Sprintf("\r\n--%s\r\n", delimeter)
	msg += "Content-Type: image/png; filename=dog.png\r\n"
	msg += "Content-Transfer-Encoding: base64\r\n"
	msg += "Content-Disposition: attachment; filename=dog.png\r\n"
	msg += "content-id: ii_1435fadb31d523f7\r\n"

	msg += "\r\n" + base64.StdEncoding.EncodeToString(dogImage)

	fmt.Println("Sending email")

	_, err = w.Write([]byte(msg))
	if err != nil {
		return err
	}
	err = w.Close()
	if err != nil {
		return err
	}
	return c.Quit()
}
