package wallpaper

import (
	"daily-wallpaper/internal/platform"
	st "daily-wallpaper/internal/settings"
	"daily-wallpaper/internal/source/registry"
	"daily-wallpaper/internal/util"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"
)

const (
	autoUpdateRetryCount    = 5
	autoUpdateRetryInterval = 30 * time.Second
	autoUpdateCheckInterval = 5 * time.Minute
)

type AutoUpdater struct {
	mu        sync.Mutex
	inProcess bool
}

func NewAutoUpdater() *AutoUpdater {
	return &AutoUpdater{}
}

func (u *AutoUpdater) RunIfNeeded() {
	u.mu.Lock()
	defer u.mu.Unlock()

	settings := st.ReadSettings()
	if !isAutoUpdateDueToday(settings, time.Now()) {
		return
	}
	var imageUrl string
	err := retry("自动获取壁纸", autoUpdateRetryCount, autoUpdateRetryInterval, func() error {
		var err error
		imageUrl, err = fetchTodayImageUrl(settings)
		return err
	})
	if err != nil {
		log.Printf("任务执行失败: %s\n", err)
		return
	}
	u.inProcess = true
	defer func() {
		u.inProcess = false
	}()
	st.WriteSettings(st.Settings{CurrentImage: &imageUrl})
}

func (u *AutoUpdater) StartMissedChecker() {
	go func() {
		ticker := time.NewTicker(autoUpdateCheckInterval)
		defer ticker.Stop()
		for range ticker.C {
			u.RunIfNeeded()
		}
	}()
}

func (u *AutoUpdater) ApplyCurrentImage(settings st.Settings) {
	if settings.CurrentImage == nil || *settings.CurrentImage == "" {
		return
	}
	var savedPath string
	err := retry("下载并设置壁纸", autoUpdateRetryCount, autoUpdateRetryInterval, func() error {
		var err error
		savedPath, err = util.GetOrDownload(*settings.CurrentImage)
		if err != nil {
			return fmt.Errorf("文件下载失败: %w", err)
		}
		if err = platform.SetWallpaper(savedPath); err != nil {
			return fmt.Errorf("设置壁纸失败: %w", err)
		}
		return nil
	})
	if err != nil {
		log.Println(err)
		return
	}
	log.Println("切换壁纸成功")
	u.markSuccess()
}

func (u *AutoUpdater) markSuccess() {
	if !u.inProcess {
		return
	}
	today := time.Now().Format("2006-01-02")
	settings := st.ReadSettings()
	if settings.AutoUpdate == nil || !*settings.AutoUpdate {
		return
	}
	if settings.LastAutoUpdateDate != nil && *settings.LastAutoUpdateDate == today {
		return
	}
	st.WriteSettings(st.Settings{LastAutoUpdateDate: &today})
}

func retry(name string, count int, interval time.Duration, fn func() error) error {
	var err error
	for i := 1; i <= count; i++ {
		if err = fn(); err == nil {
			return nil
		}
		log.Printf("%s失败，第%d/%d次: %s\n", name, i, count, err)
		if i < count {
			time.Sleep(interval)
		}
	}
	return err
}

func normalizeUpdateTime(updateTime string) (int, int, error) {
	parts := strings.Split(updateTime, ":")
	if len(parts) < 2 {
		return 0, 0, fmt.Errorf("更新时间格式错误: %s", updateTime)
	}
	t, err := time.Parse("15:04", fmt.Sprintf("%s:%s", parts[0], parts[1]))
	if err != nil {
		return 0, 0, err
	}
	return t.Hour(), t.Minute(), nil
}

func isAutoUpdateDueToday(settings st.Settings, now time.Time) bool {
	if settings.AutoUpdate == nil || !*settings.AutoUpdate || settings.TimeToUpdate == nil || *settings.TimeToUpdate == "" {
		return false
	}
	today := now.Format("2006-01-02")
	if settings.LastAutoUpdateDate != nil && *settings.LastAutoUpdateDate == today {
		return false
	}
	hour, minute, err := normalizeUpdateTime(*settings.TimeToUpdate)
	if err != nil {
		log.Println(err)
		return false
	}
	return now.Hour()*60+now.Minute() >= hour*60+minute
}

func fetchTodayImageUrl(settings st.Settings) (string, error) {
	sourceName := "bing"
	if settings.CurrentSource != nil && *settings.CurrentSource != "" {
		sourceName = *settings.CurrentSource
	}
	source := registry.Get(sourceName)
	if source == nil {
		return "", fmt.Errorf("未知壁纸源: %s", sourceName)
	}
	res, err := source.GetToday()
	if err != nil {
		return "", err
	}
	imageUrl := res.Url
	if settings.QualityFirst != nil && *settings.QualityFirst && res.UrlHS != "" {
		imageUrl = res.UrlHS
	}
	if imageUrl == "" {
		return "", fmt.Errorf("壁纸源未返回图片地址")
	}
	return imageUrl, nil
}
