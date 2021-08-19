package main

import (
	"daily-wallpaper/api"
	"daily-wallpaper/icon"
	"fmt"
	"github.com/getlantern/systray"
	"log"
)

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {
	fmt.Println("on Exit.")
}

func onReady() {

	InitSettings()
	StartServer()

	task := NewTask(func() {
		settings := ReadSettings()
		name := settings.CurrentSource
		if name == nil || *name == "" {
			*name = "bing"
		}
		source := GetSource(*name)
		res, err := source.GetToday()
		if err != nil {
			log.Printf("任务执行失败: %s\n", err)
			return
		}
		settings.CurrentImage = &res.Url
		WriteSettings(settings)
	})

	if settings.TimeToUpdate != nil && *settings.TimeToUpdate != "" {
		task.RunAt(*settings.TimeToUpdate)
	}
	if settings.AutoUpdate != nil && *settings.AutoUpdate {
		task.Start()
	}

	RegisterSettingsModifyCallback(func(settings Settings) {
		if settings.CurrentImage != nil && *settings.CurrentImage != "" {
			savedPath, err := downloadFileAndSave(*settings.CurrentImage)
			if err != nil {
				log.Printf("文件下载失败: %s\n", err)
				return
			}
			err = api.SetWallpaper(savedPath)
			if err != nil {
				log.Printf("设置壁纸失败: %s\n", err)
				return
			}
			log.Println("切换壁纸成功")
		}
		if settings.TimeToUpdate != nil && *settings.TimeToUpdate != "" {
			log.Printf("自动切换壁纸时间设置为: %s", *settings.TimeToUpdate)
			task.RunAt(*settings.TimeToUpdate)
		}
		if settings.AutoUpdate != nil {
			if *settings.AutoUpdate {
				log.Println("开启自动更新")
				task.Restart()
			} else {
				log.Println("关闭自动更新")
				task.Stop()
			}
		}
	})

	fmt.Println("on Ready.")
	systray.SetTitle("每日一图")
	systray.SetIcon(icon.Data)
	everydayItem := systray.AddMenuItemCheckbox("每日一图", "每日自动更新壁纸", true)
	moreSettingItem := systray.AddMenuItem("更多设置", "更多设置")
	quitItem := systray.AddMenuItem("退出", "退出应用程序")

	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				fmt.Println("Exit App.")
				systray.Quit()
			case <-everydayItem.ClickedCh:
				if everydayItem.Checked() {
					everydayItem.Uncheck()
					task.Stop()
					fmt.Println("关闭每日一图")
				} else {
					everydayItem.Check()
					task.Start()
					fmt.Println("开启每日更新")
				}
			case <-moreSettingItem.ClickedCh:
				OpenUrl("http://127.0.0.1:3000")
			}
		}
	}()
}
