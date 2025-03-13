package model

import (
	"time"
)

//у одного бита может быть несколько лицензий
//и несколько лицензий у одного бита

type Beat struct{
	ID uint
	Author *string
	Title *string
	License *string
	Mood *string
	Date time.Time
	Genre *string
	FreeForNonProfit *uint
}