package mail

import (
	"fmt"
	"net/smtp"
)

type SMTPEmail struct {
	Host     string `yaml:"host"`
	Username string `yaml:"username"`
	Password string `yaml:"password"`
	Port     int    `yaml:"port"`
}

func (s SMTPEmail) Send(from, frameId, image string) error {
	msg, err := getNewEmail(from, frameId)
	if err != nil {
		return err
	}

	if err := attachUrlImage(msg, image); err != nil {
		return err
	}

	emailHost := fmt.Sprintf("%s:%d", s.Host, s.Port)
	emailAuth := smtp.PlainAuth("", s.Username, s.Password, s.Host)

	if err := msg.Send(emailHost, emailAuth); err != nil {
		return err
	}

	return nil
}
