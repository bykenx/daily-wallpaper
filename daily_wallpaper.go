package main

import (
	"daily-wallpaper/api"
	"daily-wallpaper/db"
	"daily-wallpaper/icon"
	"daily-wallpaper/server"
	settings2 "daily-wallpaper/settings"
	task2 "daily-wallpaper/task"
	"daily-wallpaper/utils"
	"log"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {
	db.CloseDB()
	log.Println("on Exit.")
}

func onReady() {
	checkedChan := make(chan bool, 1)
	checkedChan2 := make(chan bool, 1)

	settings := settings2.InitSettings()
	db.OpenDB()
	server.StartServer()

	task := task2.NewTask(func() {
		settings := settings2.ReadSettings()
		source := utils.GetSource(*settings.CurrentSource)
		res, err := source.GetToday()
		if err != nil {
			log.Printf("任务执行失败: %s\n", err)
			return
		}
		settings2.WriteSettings(settings2.Settings{CurrentImage: &res.Url})
	})

	if *settings.AutoUpdate && *settings.TimeToUpdate != "" {
		log.Printf("开启自动更新，更新时间设置为: %s\n", *settings.TimeToUpdate)
		task.RunAt(*settings.TimeToUpdate).Start()
	}

	settings2.RegisterModifyCallback(func(s settings2.Settings, changed settings2.FieldChanged) {
		if changed&settings2.CurrentImageChanged != 0 {
			savedPath, err := api.DownloadFileAndSave(*s.CurrentImage)
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
		if changed&settings2.TimeToUpdateChanged != 0 {
			log.Printf("更新时间设置为: %s", *s.TimeToUpdate)
			task.RunAt(*s.TimeToUpdate)
		}
		if changed&settings2.AutoUpdateChanged != 0 {
			if *s.AutoUpdate {
				log.Println("开启自动更新")
				checkedChan <- true
				task.Restart()
			} else {
				log.Println("关闭自动更新")
				checkedChan <- false
				task.Stop()
			}
		}
		if changed&settings2.AutoRunAtSystemBootChanged != 0 {
			if *s.AutoRunAtSystemBoot {
				log.Println("开启开机自启")
				checkedChan2 <- true
			} else {
				log.Println("关闭开机自启")
				checkedChan2 <- false
			}
			api.SetStartAtLogin(*s.AutoRunAtSystemBoot)
		}
	})

	//systray.SetTitle("每日一图")
	systray.SetIcon(icon.Data)
	startAtLoginItem := systray.AddMenuItemCheckbox("开机自启", "开机自启", *settings.AutoRunAtSystemBoot)
	everydayItem := systray.AddMenuItemCheckbox("每日一图", "每日自动更新壁纸", *settings.AutoUpdate)
	moreSettingItem := systray.AddMenuItem("更多设置", "更多设置")
	quitItem := systray.AddMenuItem("退出", "退出应用程序")

	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				log.Println("Exit App.")
				systray.Quit()
			case <-everydayItem.ClickedCh:
				checked := !everydayItem.Checked()
				settings2.WriteSettings(settings2.Settings{AutoUpdate: &checked})
			case <-moreSettingItem.ClickedCh:
				utils.OpenUrl("http://127.0.0.1:9001")
			case <-startAtLoginItem.ClickedCh:
				checked := !startAtLoginItem.Checked()
				settings2.WriteSettings(settings2.Settings{AutoRunAtSystemBoot: &checked})
			case v := <-checkedChan:
				if v {
					everydayItem.Check()
				} else {
					everydayItem.Uncheck()
				}
			case v := <-checkedChan2:
				if v {
					startAtLoginItem.Check()
				} else {
					startAtLoginItem.Uncheck()
				}
			}
		}
	}()
}
