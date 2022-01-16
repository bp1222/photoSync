package main

type Entry struct {
	EntryId   int64 `gorm:"primaryKey"`
	Timestamp int64
}

type Like struct {
	EntryId   int64 `gorm:"primaryKey"`
	UserId    int64 `gorm:"primaryKey"`
	Timestamp int64
}
