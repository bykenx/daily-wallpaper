package settings

import (
	"daily-wallpaper/constant"
	"daily-wallpaper/utils"
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
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
}

type ModifyCallback func(prev, current Settings)

var lastModifyTime int64
var settings Settings
var modifyCallback ModifyCallback

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
		}
		configBytes, _ := yaml.Marshal(settings)
		_ = ioutil.WriteFile(p, configBytes, constant.DefaultFileCreatePermission)
	}
	settings = ReadSettings()
	lastModifyTime = utils.GetLastModifyTime(p)
	return settings
}

func ReadSettings() Settings {
	p := filepath.Join(constant.AppHome, constant.ConfigFileName)
	if lastModifyTime < utils.GetLastModifyTime(p) {
		configFilePath := filepath.Join(constant.AppHome, constant.ConfigFileName)
		configBytes, _ := ioutil.ReadFile(configFilePath)
		_ = yaml.Unmarshal(configBytes, &settings)
	}
	return settings
}

func WriteSettings(settings Settings) {
	dst := ReadSettings()
	if modifyCallback != nil {
		modifyCallback(dst, settings)
	}
	_ = mergo.Merge(&dst, settings, mergo.WithOverride)
	configFilePath := filepath.Join(constant.AppHome, constant.ConfigFileName)
	configBytes, _ := yaml.Marshal(dst)
	_ = ioutil.WriteFile(configFilePath, configBytes, constant.DefaultFileCreatePermission)
}

func RegisterModifyCallback(callback ModifyCallback) {
	modifyCallback = callback
}
