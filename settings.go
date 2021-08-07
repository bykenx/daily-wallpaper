package main

import (
	"fmt"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"path/filepath"
)

func initSettings() {
	if !IsDir(appHome) {
		_ = os.Mkdir(appHome, defaultDirectoryCreatePermission)
	}
	p := filepath.Join(appHome, configFileName)
	if !IsFile(p) {
		settings := Settings{
			CurrentImage:        "",
			AutoUpdate:          false,
			AutoRunAtSystemBoot: false,
		}
		configBytes, _ := yaml.Marshal(settings)
		_ = ioutil.WriteFile(p, configBytes, defaultFileCreatePermission)
	}
}



func readSettings() Settings {
	var settings Settings
	configFilePath := filepath.Join(appHome, configFileName)
	configBytes, _ := ioutil.ReadFile(configFilePath)
	_ = yaml.Unmarshal(configBytes, &settings)
	return settings
}

func writeSettings(settings Settings) {
	home, _ := os.UserHomeDir()
	configFilePath := filepath.Join(home, fmt.Sprintf(".%s", packageName), configFileName)
	configBytes, _ := yaml.Marshal(settings)
	_ = ioutil.WriteFile(configFilePath, configBytes, defaultFileCreatePermission)
}
