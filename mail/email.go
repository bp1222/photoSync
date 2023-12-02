package mail

import (
	"fmt"
	"math/rand"
	"net/http"
	"path"

	"github.com/jordan-wright/email"
)

func getNewEmail(from, frameId string) (*email.Email, error) {
	msg := email.NewEmail()
	msg.Subject = fmt.Sprintf("Image Upload: %s", randStringBytes(8))
	msg.From = from
	msg.To = []string{fmt.Sprintf("%s@send.auraframes.com", frameId)}
	return msg, nil
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

func randStringBytes(n int) string {
	const letterBytes = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	b := make([]byte, n)
	for i := range b {
		b[i] = letterBytes[rand.Intn(len(letterBytes))]
	}
	return string(b)
}
