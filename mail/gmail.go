package mail

import (
	"context"
	"encoding/base64"
	"os"

	"golang.org/x/oauth2/google"
	"google.golang.org/api/gmail/v1"
	"google.golang.org/api/option"

	log "github.com/sirupsen/logrus"
)

type GmailEmail struct {
}

func (_ GmailEmail) getService(from string) *gmail.Service {
	ctx := context.Background()
	b, err := os.ReadFile("credentials.json")
	if err != nil {
		log.Fatalf("Unable to read client secret file: %v", err)
	}

	config, err := google.JWTConfigFromJSON(b, gmail.GmailSendScope)
	if err != nil {
		log.Fatalf("Unable to parse client secret file to config: %v", err)
	}

	// Act on authority of
	config.Subject = from

	svc, err := gmail.NewService(ctx, option.WithHTTPClient(config.Client(ctx)))
	if err != nil {
		log.Fatalf("Unable to retrieve Gmail client: %v", err)
	}

	return svc
}

func (g GmailEmail) Send(from, frameId, image string) error {
	svc := g.getService(from)

	msg, err := getNewEmail(from, frameId)
	if err != nil {
		return err
	}

	if err := attachUrlImage(msg, image); err != nil {
		return err
	}

	bytes, err := msg.Bytes()
	if err != nil {
		return err
	}

	message := &gmail.Message{
		Raw: base64.URLEncoding.EncodeToString(bytes),
	}

	_, err = svc.Users.Messages.Send("me", message).Do()
	if err != nil {
		return err
	}
	return nil
}
