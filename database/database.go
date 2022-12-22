package database

import (
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
)

const (
	DB_FILE = "./tinybeans_photos.db"
)

type Database interface {
	IsLiked(entryId, userId int64, frameId string) bool
	SaveLike(entryId, userId int64, frameId string, timestamp int64)
	SaveEntry(entryId, journalId, timestamp int64)
	GetMostRecentEntry(journal int64) int64
}

type database struct {
	*gorm.DB
}

func InitDatabase() Database {
	db, err := gorm.Open(sqlite.Open(DB_FILE), &gorm.Config{
		Logger: logger.Default.LogMode(logger.Silent),
	})
	if err != nil {
		log.Fatal("unable to open database")
	}

	db.AutoMigrate(&Entry{})
	db.AutoMigrate(&Like{})

	return &database{
		db,
	}
}

func (d database) IsLiked(entryId, userId int64, frameId string) bool {
	var l Like
	d.Where(d.Where(Like{
		EntryId: entryId,
		UserId:  userId,
	}).Or(Like{
		EntryId: entryId,
		FrameId: frameId,
	})).First(&l)
	return l.EntryId == entryId
}

func (d database) GetMostRecentEntry(journal int64) int64 {
	var e Entry
	d.Order("timestamp desc").Where(Entry{
		JournalId: journal,
	}).First(&e)
	return e.Timestamp
}

func (d database) SaveLike(entryId, userId int64, frameId string, timestamp int64) {
	d.Save(Like{
		EntryId:   entryId,
		UserId:    userId,
		FrameId:   frameId,
		Timestamp: timestamp,
	})
}

func (d database) SaveEntry(entryId, journalId, timestamp int64) {
	d.Save(Entry{
		EntryId:   entryId,
		JournalId: journalId,
		Timestamp: timestamp,
	})
}
