package main

var (
	StringEmpty = ""
	BooleanFalse = false
)

type Settings struct {
	CurrentSource       *string `json:"currentSource" yaml:"current_source"`
	CurrentImage        *string `json:"currentImage" yaml:"current_image"`
	AutoUpdate          *bool   `json:"autoUpdate" yaml:"auto_update"`
	TimeToUpdate        *string `json:"timeToUpdate" yaml:"time_to_update"`
	AutoRunAtSystemBoot *bool   `json:"autoRunAtSystemBoot" yaml:"auto_run_at_system_boot"`
	QualityFirst        *bool   `json:"qualityFirst" yaml:"quality_first"`
}

type DownloadHistory struct {
	Id         uint   `gorm:"primaryKey;autoIncrementIncrement"`
	Dir        string `gorm:""`
	Url        string `gorm:"index;not null"`
	CreateTime int64  `gorm:"autoCreateTime:milli"`
	UpdateTime int64  `gorm:"autoUpdateTime:milli"`
	Valid      bool   `gorm:"default:1"`
}
