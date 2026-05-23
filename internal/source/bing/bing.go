package bing

import (
	"daily-wallpaper/internal/source"
	"fmt"
	"time"
)

const (
	UrlPrefix         = "https://cn.bing.com/HPImageArchive.aspx?%s"
	ImageUrlPrefix    = "https://cn.bing.com%s"
	NormalImageWidth  = 1920
	NormalImageHeight = 1080
	UHDImageWidth     = 3840
	UHDImageHeight    = 2160
)

type Source struct{}

func (s Source) GetToday() (source.TodayResponse, error) {
	payload := RequestPayload{
		Format:    "js",
		Index:     0,
		PageSize:  1,
		Timestamp: time.Now().Unix(),
		Pid:       "hp",
		UHD:       1,
		UHDWidth:  UHDImageWidth,
		UHDHeight: UHDImageHeight,
	}
	var result Response
	err := source.DispatchGetRequest(UrlPrefix, payload, &result)
	if err != nil {
		return source.ImageItem{}, err
	}
	item := result.Images[0]
	return buildImageItem(item), nil
}

func (s Source) GetArchive(param source.ArchiveParam) (source.ArchiveResponse, error) {
	payload := RequestPayload{
		Format:    "js",
		Index:     1,
		PageSize:  8,
		Timestamp: time.Now().Unix(),
		Pid:       "hp",
		UHD:       1,
		UHDWidth:  UHDImageWidth,
		UHDHeight: UHDImageHeight,
	}
	items := make([]source.ImageItem, 0)
	var result Response
	err := source.DispatchGetRequest(UrlPrefix, payload, &result)
	if err == nil {
		for _, item := range result.Images {
			items = append(items, buildImageItem(item))
		}
	}
	return source.ArchiveResponse{
		Items:    items,
		End:      true,
		Current:  0,
		PageSize: 8,
	}, nil
}

func buildImageItem(item ImageItem) source.ImageItem {
	normalUrl := fmt.Sprintf("%s_%dx%d.jpg", item.UrlBase, NormalImageWidth, NormalImageHeight)
	return source.ImageItem{
		Name:          item.Title,
		Url:           fmt.Sprintf(ImageUrlPrefix, normalUrl),
		UrlHS:         fmt.Sprintf(ImageUrlPrefix, item.Url),
		Description:   "",
		Copyright:     item.Copyright,
		CopyrightLink: item.CopyrightLink,
		SourceLink:    "",
		Author:        "",
		Location:      "",
	}
}
