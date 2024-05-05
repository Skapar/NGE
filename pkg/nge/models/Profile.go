package models

type Profile struct {
	Id   uint   `json:"id" gorm:"primaryKey;autoIncrement"`
	Text string `json:"text"`
	
}

