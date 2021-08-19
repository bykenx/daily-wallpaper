package main

import (
	"github.com/imdario/mergo"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

var lastModifyTime int64
var settings Settings

func initSettings() {
	if !IsDir(appHome) {
		_ = os.Mkdir(appHome, defaultDirectoryCreatePermission)
	}
	p := filepath.Join(appHome, configFileName)
	if !IsFile(p) {
		settings := Settings{
			CurrentSource:       &StringEmpty,
			CurrentImage:        &StringEmpty,
			AutoUpdate:          &BooleanFalse,
			TimeToUpdate:        &StringEmpty,
			AutoRunAtSystemBoot: &BooleanFalse,
			QualityFirst:        &BooleanFalse,
		}
		configBytes, _ := yaml.Marshal(settings)
		_ = ioutil.WriteFile(p, configBytes, defaultFileCreatePermission)
	}
	settings = readSettings()
	lastModifyTime = GetLastModifyTime(p)
}

func readSettings() Settings {
	p := filepath.Join(appHome, configFileName)
	if lastModifyTime < GetLastModifyTime(p) {
		configFilePath := filepath.Join(appHome, configFileName)
		configBytes, _ := ioutil.ReadFile(configFilePath)
		_ = yaml.Unmarshal(configBytes, &settings)
	}
	return settings
}

func writeSettings(settings Settings) {
	dst := readSettings()
	_ = mergo.Merge(&dst, settings, mergo.WithOverride)
	configFilePath := filepath.Join(appHome, configFileName)
	configBytes, _ := yaml.Marshal(dst)
	_ = ioutil.WriteFile(configFilePath, configBytes, defaultFileCreatePermission)
}
