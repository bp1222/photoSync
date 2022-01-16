package main

import (
	"fmt"
	"log"
	"math/rand"
	"os"
	"time"

	"github.com/joho/godotenv"
)

const (
	TINYBEANS_JOURNAL = int64(1328947)
	TINYBEANS_DAVE_ID = int64(3248901)
)

func main() {
	if err := godotenv.Load(); err != nil {
		log.Fatal("unable to load environment")
	}

	rand.Seed(time.Now().UTC().UnixNano())

	db, err := GetDatabase()
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	tiny := InitTinybeans(db)

	if err := tiny.Authenticate(); err != nil {
		log.Fatal("unable to authenticate to tinybeans")
	}

	since := tiny.GetMostRecentEntry(TINYBEANS_JOURNAL)
	firstRun := since == 0

	for {
		entries, _ := tiny.GetJournalEntriesSince(TINYBEANS_JOURNAL, 200, since)

		if entries.Entries == nil {
			break
		}
		for _, entry := range *entries.Entries {
			fmt.Fprintf(os.Stdout, "Entry %d, %s, %d\n", entry.GetId(), entry.GetClientRef(), entry.GetLastUpdatedTimestamp())

			db.Save(&Entry{
				EntryId:   entry.GetId(),
				Timestamp: entry.GetLastUpdatedTimestamp(),
			})

			if entry.Emotions != nil {
				for _, emotion := range *entry.Emotions {
					if emotion.GetUserId() == TINYBEANS_DAVE_ID {
						if !firstRun && !tiny.IsLikedBy(entry.GetId(), TINYBEANS_DAVE_ID) {
							SendAuraEmail(*entry.Blobs.O)
						}
						db.Save(Like{
							EntryId:   entry.GetId(),
							UserId:    emotion.GetId(),
							Timestamp: emotion.GetLastUpdatedTimestamp(),
						})
					}
				}
			}
			since = entry.GetLastUpdatedTimestamp()
		}

		if entries.GetNumEntriesRemaining() == 0 {
			break
		}
		time.Sleep(time.Second * 2)
	}
	fmt.Println("Done with Loading Previous Images")
}
