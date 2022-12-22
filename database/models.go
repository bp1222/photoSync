package database

type Entry struct {
	EntryId   int64 `gorm:"primaryKey"`
	JournalId int64 `gorm:"primaryKey"`
	Timestamp int64
}

type Like struct {
	EntryId   int64  `gorm:"primaryKey"`
	UserId    int64  `gorm:"primaryKey"`
	FrameId   string `gorm:"primaryKey"`
	Timestamp int64
}
