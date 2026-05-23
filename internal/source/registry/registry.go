package registry

import (
	"daily-wallpaper/internal/source"
	"daily-wallpaper/internal/source/bing"
	"daily-wallpaper/internal/source/unsplash"
)

func Get(name string) source.Source {
	switch name {
	case "bing":
		return bing.Source{}
	case "unsplash":
		return unsplash.Source{}
	default:
		return nil
	}
}

func Descriptions() []source.Description {
	return []source.Description{
		{Name: "bing", Description: "bing壁纸", Enabled: true},
		{Name: "unsplash", Description: "Unsplash", Enabled: true},
	}
}
