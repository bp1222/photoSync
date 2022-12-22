package main

import (
	"math/rand"
	"os"
	"time"

	"github.com/bp1222/photoSync/mail"
	tbApi "github.com/bp1222/tinybeans-api/go-client"
	"github.com/joho/godotenv"
	log "github.com/sirupsen/logrus"

	"github.com/bp1222/photoSync/database"
	"github.com/bp1222/photoSync/tinybeans"
)

var (
	config    = Config{}
	Tinybeans tinybeans.Tinybeans
	Database  database.Database
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load environment")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	loadConfig()

	Database = database.InitDatabase()
	Tinybeans = tinybeans.InitTinybeans(Database)

	for _, journal := range config.Journals {
		log.Infof("Iterating journal (%d)", journal.Id)
		doProcessTinybeansJournal(journal)
	}
}

func doProcessTinybeansJournal(journal Journal) {
	since := Database.GetMostRecentEntry(journal.Id)
	log.Infof("Most recent entry for journal (%d) on (%d)", journal.Id, since)

	for {
		entries, _ := Tinybeans.GetJournalEntriesSince(journal.Id, 200, since)
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

	Database.SaveEntry(entry.GetId(), journalId, entry.GetLastUpdatedTimestamp())
}

func isUserTrackedForJournal(userId, journalId int64) *User {
	var trackedUsers []User

	for _, j := range config.Journals {
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

func doProcessTinybeansEntryEmotion(journalId int64, entry tbApi.Entry, emotion tbApi.Emotion) {
	if user := isUserTrackedForJournal(emotion.GetUserId(), journalId); user != nil {
		log.Infof("Found emotion on journal (%d), on entry (%d) for user (%d)", journalId, entry.GetId(), user.Id)
		for _, frameId := range user.FrameIds {
			if !Database.IsLiked(entry.GetId(), user.Id, frameId) {
				if doMail, found := os.LookupEnv("LIVE_SEND_MAIL"); found && doMail == "true" {
					log.Infof("Image being sent to Aura Frame (%s): %s", frameId, *entry.Blobs.O)
					mail.SendAuraEmail(frameId, *entry.Blobs.O)
				} else {
					log.Infof("TEST: Image would be sent to Aura Frame: %s", *entry.Blobs.O)
				}
			}
			Database.SaveLike(entry.GetId(), user.Id, frameId, emotion.GetLastUpdatedTimestamp())
		}
	}
}
