package main

import (
	"daily-wallpaper/icon"
	"fmt"

	"github.com/getlantern/systray"
)

func main() {
	systray.Run(onReady, onExit)
}

func onExit() {
	fmt.Println("on Exit.")
}

func onReady() {
	fmt.Println("on Ready.")
	systray.SetTitle("每日一图")
	systray.SetIcon(icon.Data)
	everydayItem := systray.AddMenuItemCheckbox("每日一图", "每日自动更新壁纸", true)
	moreSettingItem := systray.AddMenuItem("更多设置", "更多设置")
	quitItem := systray.AddMenuItem("退出", "退出应用程序")

	initSettings()
	startServer()

	go func() {
		for {
			select {
			case <-quitItem.ClickedCh:
				fmt.Println("Exit App.")
				systray.Quit()
			case <-everydayItem.ClickedCh:
				if everydayItem.Checked() {
					everydayItem.Uncheck()
					fmt.Println("关闭每日一图")
				} else {
					everydayItem.Check()
					fmt.Println("开启每日更新")
				}
			case <-moreSettingItem.ClickedCh:
				OpenUrl("http://127.0.0.1:3000")
			}
		}
	}()
}
