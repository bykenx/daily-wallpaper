package bing_source

import (
	"daily-wallpaper/sources"
	"fmt"
	"time"
)

type BingSource struct {
}

func (s BingSource) GetToday() (sources.TodayResponse, error) {
	payload := RequestPayload{
		Format:    "js",
		Index:     0,
		PageSize:  1,
		Timestamp: time.Now().Unix(),
	}
	items := dispatchRequest(payload)
	if len(items) <= 0 {
		return sources.ImageItem{}, nil
	}
	return sources.ImageItem{
		Name:          items[0].Title,
		Url:           fmt.Sprintf("https://cn.bing.com%s", items[0].Url),
		Description:   "",
		Copyright:     items[0].Copyright,
		CopyrightLink: items[0].CopyrightLink,
		SourceLink:    "",
		Author:        "",
		Location:      "",
	}, nil
}

func (s BingSource) GetArchive(param sources.ArchiveParam) (sources.ArchiveResponse, error) {
	payload := RequestPayload{
		Format:    "js",
		Index:     1,
		PageSize:  8,
		Timestamp: time.Now().Unix(),
	}
	var items []sources.ImageItem
	res := dispatchRequest(payload)
	for _, item := range res {
		items = append(items, sources.ImageItem{
			Name:          item.Title,
			Url:           fmt.Sprintf("https://cn.bing.com%s", item.Url),
			Description:   "",
			Copyright:     item.Copyright,
			CopyrightLink: item.CopyrightLink,
			SourceLink:    "",
			Author:        "",
			Location:      "",
		})
	}
	return sources.ArchiveResponse{
		Items:    items,
		End:      true,
		Current:  0,
		PageSize: 8,
	}, nil
}
