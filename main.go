package main

import (
	"fmt"
	"net/http"
	"net/url"
	"time"

	"github.com/bp1222/photoSync/database"
	"github.com/bp1222/photoSync/mail"
	"github.com/bp1222/photoSync/tinybeans"
	tbApi "github.com/bp1222/tinybeans-api/go-client"
	log "github.com/sirupsen/logrus"
)

var (
	config = Config{}
	tb     tinybeans.Tinybeans
	db     database.Database
	sender mail.Sender
)

func main() {
	loadConfig()
	setupDatabase()
	setupTinybeans()
	setupSender()

	for _, journal := range config.Tinybeans.Journals {
		log.Infof("Iterating journal (%d)", journal.Id)
		doProcessTinybeansJournal(journal)
	}
}

func setupDatabase() {
	db = database.InitDatabase()
}

func setupTinybeans() {
	tinybeansOpts := []tinybeans.OptionFunc{
		tinybeans.WithDatabase(db),
	}

	if config.Mitm != nil {
		proxyUrl, _ := url.Parse(fmt.Sprintf("%s:%d", config.Mitm.Host, config.Mitm.Port))
		tinybeansOpts = append(tinybeansOpts, tinybeans.WithClient(&http.Client{
			Transport: &http.Transport{
				Proxy: http.ProxyURL(proxyUrl),
			},
		}))
	}

	tb = tinybeans.InitTinybeans(config.Tinybeans, tinybeansOpts...)
}

func setupSender() {
	if (config.Sender.Smtp == nil && config.Sender.Gmail == nil) ||
		(config.Sender.Smtp != nil && config.Sender.Gmail != nil) {
		log.Fatal("Must have one, and only one sender config")
	}

	if config.Sender.Smtp != nil {
		sender = mail.SMTPEmail{}
	}

	if config.Sender.Gmail != nil {
		sender = mail.GmailEmail{}
	}
}

func doProcessTinybeansJournal(journal tinybeans.Journal) {
	since := db.GetMostRecentEntry(journal.Id)
	log.Infof("Most recent entry for journal (%d) on (%d)", journal.Id, since)

	for {
		entries, _ := tb.GetJournalEntriesSince(journal.Id, 200, since)
		log.Infof("Found %d new entries", len(entries.GetEntries()))

		if entries.Entries == nil {
			break
		}

		for _, entry := range entries.Entries {
			log.Infof("Processing journal (%d) entry (%d)", journal.Id, entry.GetId())
			doProcessTinybeansEntry(journal.Id, entry)
			since = entry.GetLastUpdatedTimestamp()
		}

		if entries.GetNumEntriesRemaining() == 0 {
			break
		}

		// Be Nice
		time.Sleep(time.Second * 2)
	}
}

func doProcessTinybeansEntry(journalId int64, entry tbApi.Entry) {
	if entry.Emotions != nil {
		for _, emotion := range entry.Emotions {
			doProcessTinybeansEntryEmotion(journalId, entry, emotion)
		}
	}

	db.SaveEntry(entry.GetId(), journalId, entry.GetLastUpdatedTimestamp())
}

func doProcessTinybeansEntryEmotion(journalId int64, entry tbApi.Entry, emotion tbApi.Emotion) {
	if user := isUserTrackedForJournal(emotion.GetUserId(), journalId); user != nil {
		log.Infof("Found emotion on journal (%d), on entry (%d) for user (%d)", journalId, entry.GetId(), user.Id)
		for _, frameId := range user.FrameIds {
			if !db.IsLiked(entry.GetId(), user.Id, frameId) {
				if config.Live {
					log.Infof("Image being sent to Aura Frame (%s): %s", frameId, *entry.Blobs.O)
					if err := sender.Send(config.Sender.From, frameId, *entry.Blobs.O); err != nil {
						log.Fatal("Email failed to send", err)
					}
				} else {
					log.Infof("TEST: Image would be sent to Aura Frame: %s", frameId)
				}
			}
			db.SaveLike(entry.GetId(), user.Id, frameId, emotion.GetLastUpdatedTimestamp())
		}
	}
}

func isUserTrackedForJournal(userId, journalId int64) *tinybeans.User {
	var trackedUsers []tinybeans.User

	for _, j := range config.Tinybeans.Journals {
		if journalId == j.Id {
			trackedUsers = j.Users
			break
		}
	}

	if trackedUsers == nil {
		return nil
	}

	for _, user := range trackedUsers {
		if user.Id == userId {
			return &user
		}
	}

	return nil
}
