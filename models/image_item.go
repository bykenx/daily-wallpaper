package models

type ImageItem struct {
	Id         uint   `gorm:"primaryKey;autoIncrementIncrement"`
	Dir        string `gorm:""`
	Url        string `gorm:"index;not null"`
	CreateTime int64  `gorm:"autoCreateTime:milli"`
	UpdateTime int64  `gorm:"autoUpdateTime:milli"`
	Valid      bool   `gorm:"default:1"`
	Sha1       string `gorm:""`
}
