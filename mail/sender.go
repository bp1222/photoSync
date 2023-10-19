package mail

import (
	"fmt"
	"math/rand"
	"net/http"
	"net/smtp"
	"os"
	"path"

	"github.com/jordan-wright/email"
)

func SendAuraEmail(frameId, image string) error {
	msg := email.NewEmail()
	msg.Subject = fmt.Sprintf("Image Upload: %s", randStringBytes(8))
	msg.From = os.Getenv("EMAIL_USERNAME")
	msg.To = []string{fmt.Sprintf("%s@send.auraframes.com", frameId)}
	if err := attachUrlImage(msg, image); err != nil {
		return err
	}

	emailHost := fmt.Sprintf("%s:%s", os.Getenv("EMAIL_HOST"), os.Getenv("EMAIL_PORT"))
	emailAuth := smtp.PlainAuth("", os.Getenv("EMAIL_USERNAME"), os.Getenv("EMAIL_PASSWORD"), os.Getenv("EMAIL_HOST"))
	if err := msg.Send(emailHost, emailAuth); err != nil {
		return err
	}

	return nil
}

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}

func attachUrlImage(msg *email.Email, url string) error {
	resp, err := http.Get(url)
	if err != nil {
		return err
	}
	defer resp.Body.Close()

	if _, err := msg.Attach(resp.Body, path.Base(resp.Request.URL.Path), "image/jpeg"); err != nil {
		return err
	}

	return nil
}
