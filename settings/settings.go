package settings

import (
	"daily-wallpaper/constant"
	"daily-wallpaper/utils"
	"os"
	"path/filepath"

	"gopkg.in/yaml.v2"
)

type ModifyCallback func(settings Settings, changed FieldChanged)

var (
	lastModifyTime int64
	settings       Settings
	modifyCallback ModifyCallback
)

var (
	StringEmpty  = ""
	BooleanFalse = false
)

type Settings struct {
	CurrentSource       *string `json:"currentSource" yaml:"current_source"`
	CurrentImage        *string `json:"currentImage" yaml:"current_image"`
	AutoUpdate          *bool   `json:"autoUpdate" yaml:"auto_update"`
	TimeToUpdate        *string `json:"timeToUpdate" yaml:"time_to_update"`
	AutoRunAtSystemBoot *bool   `json:"autoRunAtSystemBoot" yaml:"auto_run_at_system_boot"`
	QualityFirst        *bool   `json:"qualityFirst" yaml:"quality_first"`
	LastAutoUpdateDate  *string `json:"lastAutoUpdateDate" yaml:"last_auto_update_date"`
}

type FieldChanged int

const (
	CurrentSourceChanged = 1 << iota
	CurrentImageChanged
	AutoUpdateChanged
	TimeToUpdateChanged
	AutoRunAtSystemBootChanged
	QualityFirstChanged
	LastAutoUpdateDateChanged
)

func mergeSettings(dst *Settings, src Settings) (res FieldChanged) {
	res = 0
	if src.CurrentSource != nil && dst.CurrentSource != src.CurrentSource {
		dst.CurrentSource = src.CurrentSource
		res += CurrentSourceChanged
	}
	if src.CurrentImage != nil && dst.CurrentImage != src.CurrentImage {
		dst.CurrentImage = src.CurrentImage
		res += CurrentImageChanged
	}
	if src.AutoUpdate != nil && dst.AutoUpdate != src.AutoUpdate {
		dst.AutoUpdate = src.AutoUpdate
		res += AutoUpdateChanged
	}
	if src.TimeToUpdate != nil && dst.TimeToUpdate != src.TimeToUpdate {
		dst.TimeToUpdate = src.TimeToUpdate
		res += TimeToUpdateChanged
	}
	if src.AutoRunAtSystemBoot != nil && dst.AutoRunAtSystemBoot != src.AutoRunAtSystemBoot {
		dst.AutoRunAtSystemBoot = src.AutoRunAtSystemBoot
		res += AutoRunAtSystemBootChanged
	}
	if src.QualityFirst != nil && dst.QualityFirst != src.QualityFirst {
		dst.QualityFirst = src.QualityFirst
		res += QualityFirstChanged
	}
	if src.LastAutoUpdateDate != nil && dst.LastAutoUpdateDate != src.LastAutoUpdateDate {
		dst.LastAutoUpdateDate = src.LastAutoUpdateDate
		res += LastAutoUpdateDateChanged
	}
	return
}

func InitSettings() Settings {
	if !utils.IsDir(constant.AppHome) {
		_ = os.Mkdir(constant.AppHome, constant.DefaultDirectoryCreatePermission)
	}
	p := filepath.Join(constant.AppHome, constant.ConfigFileName)
	if !utils.IsFile(p) {
		settings := Settings{
			CurrentSource:       &StringEmpty,
			CurrentImage:        &StringEmpty,
			AutoUpdate:          &BooleanFalse,
			TimeToUpdate:        &StringEmpty,
			AutoRunAtSystemBoot: &BooleanFalse,
			QualityFirst:        &BooleanFalse,
			LastAutoUpdateDate:  &StringEmpty,
		}
		configBytes, _ := yaml.Marshal(settings)
		_ = os.WriteFile(p, configBytes, constant.DefaultFileCreatePermission)
	}
	settings = ReadSettings()
	lastModifyTime = utils.GetLastModifyTime(p)
	return settings
}

func ReadSettings() Settings {
	p := filepath.Join(constant.AppHome, constant.ConfigFileName)
	if lastModifyTime < utils.GetLastModifyTime(p) {
		configFilePath := filepath.Join(constant.AppHome, constant.ConfigFileName)
		configBytes, _ := os.ReadFile(configFilePath)
		_ = yaml.Unmarshal(configBytes, &settings)
	}
	return settings
}

func WriteSettings(src Settings) {
	dst := ReadSettings()
	res := mergeSettings(&dst, src)
	configFilePath := filepath.Join(constant.AppHome, constant.ConfigFileName)
	configBytes, _ := yaml.Marshal(dst)
	_ = os.WriteFile(configFilePath, configBytes, constant.DefaultFileCreatePermission)
	settings = dst
	lastModifyTime = utils.GetLastModifyTime(configFilePath)
	if res != 0 && modifyCallback != nil {
		modifyCallback(dst, res)
	}
}

func RegisterModifyCallback(callback ModifyCallback) {
	modifyCallback = callback
}
