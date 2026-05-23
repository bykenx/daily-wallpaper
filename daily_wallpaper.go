package main

import (
	"daily-wallpaper/api"
	"daily-wallpaper/db"
	"daily-wallpaper/icon"
	"daily-wallpaper/server"
	st "daily-wallpaper/settings"
	tk "daily-wallpaper/task"
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

	settings := st.InitSettings()
	db.OpenDB()
	server.StartServer()

	autoUpdater := tk.NewAutoUpdater()
	task := tk.NewTask(autoUpdater.RunIfNeeded)

	if *settings.AutoUpdate && *settings.TimeToUpdate != "" {
		log.Printf("开启自动更新，更新时间设置为: %s\n", *settings.TimeToUpdate)
		task.RunAt(*settings.TimeToUpdate).Start()
	}
	autoUpdater.StartMissedChecker()

	st.RegisterModifyCallback(func(s st.Settings, changed st.FieldChanged) {
		if changed&st.CurrentImageChanged != 0 {
			autoUpdater.ApplyCurrentImage(s)
		}
		if changed&st.TimeToUpdateChanged != 0 {
			log.Printf("更新时间设置为: %s", *s.TimeToUpdate)
			task.RunAt(*s.TimeToUpdate).Restart()
		}
		if changed&st.AutoUpdateChanged != 0 {
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
		if changed&st.AutoRunAtSystemBootChanged != 0 {
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
	api.SetTrayIcon(icon.Data)
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
				st.WriteSettings(st.Settings{AutoUpdate: &checked})
			case <-moreSettingItem.ClickedCh:
				utils.OpenUrl("http://127.0.0.1:9001")
			case <-startAtLoginItem.ClickedCh:
				checked := !startAtLoginItem.Checked()
				st.WriteSettings(st.Settings{AutoRunAtSystemBoot: &checked})
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
