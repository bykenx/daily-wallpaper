package bing_source

import (
	"daily-wallpaper/sources"
	"fmt"
	"time"
)

const (
	UrlPrefix = "https://cn.bing.com/HPImageArchive.aspx?%s"
)

type BingSource struct{}

func (s BingSource) GetToday() (sources.TodayResponse, error) {
	payload := RequestPayload{
		Format:    "js",
		Index:     0,
		PageSize:  1,
		Timestamp: time.Now().Unix(),
	}
	var result Response
	err := sources.DispatchGetRequest(UrlPrefix, payload, &result)
	if err != nil {
		return sources.ImageItem{}, err
	}
	item := result.Images[0]
	return sources.ImageItem{
		Name:          item.Title,
		Url:           fmt.Sprintf("https://cn.bing.com%s", item.Url),
		Description:   "",
		Copyright:     item.Copyright,
		CopyrightLink: item.CopyrightLink,
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
	var result Response
	err := sources.DispatchGetRequest(UrlPrefix, payload, &result)
	if err == nil {
		for _, item := range result.Images {
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
	}
	return sources.ArchiveResponse{
		Items:    items,
		End:      true,
		Current:  0,
		PageSize: 8,
	}, nil
}
