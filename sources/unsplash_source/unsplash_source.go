package unsplash_source

import (
	"daily-wallpaper/sources"
	"github.com/anaskhan96/soup"
	"log"
	"strings"
)

const (
	UrlPrefix = "https://unsplash.com/napi/topics/wallpapers/photos?%s"
)

type UnsplashSource struct {
}

func (u UnsplashSource) GetToday() (sources.TodayResponse, error) {
	res, err := soup.Get("https://unsplash.com/")
	if err != nil {
		return sources.TodayResponse{}, err
	}
	current := soup.HTMLParse(res)
	current = current.Find("div", "data-test", "editorial-route")
	if current.Error == nil {
		current = current.Find("picture")
		if current.Error == nil {
			items := current.FindAll("source")
			if len(items) > 0 {
				var url, urlHS, urlP string
				for _, item := range items {
					attrs := item.Attrs()
					if strings.Contains(attrs["media"], "(min-width: 6400px)") {
						urlHS = sources.GetSafeUrl(attrs["srcset"])
					} else if strings.Contains(attrs["media"], "(min-width: 2000px)") {
						url = sources.GetSafeUrl(attrs["srcset"])
					} else if strings.Contains(attrs["media"], "(min-width: 200px)") {
						urlP = sources.GetSafeUrl(attrs["srcset"])
					}
				}
				return sources.ImageItem{
					Url:   url,
					UrlP:  urlP,
					UrlHS: urlHS,
				}, err
			}
		}
	}
	log.Fatal("查找节点失败")
	return sources.TodayResponse{}, nil
}

func (u UnsplashSource) GetArchive(param sources.ArchiveParam) (sources.ArchiveResponse, error) {
	payload := RequestPayload{
		Page:    param.Current,
		PerPage: param.PageSize,
	}
	var items []sources.ImageItem
	var result Response
	err := sources.DispatchGetRequest(UrlPrefix, payload, &result)
	if err == nil {
		for _, item := range result {
			items = append(items, sources.ImageItem{
				Name:          "",
				Url:           item.Urls.Regular,
				UrlHS:         item.Urls.Raw,
				UrlP:          item.Urls.Thumb,
				Description:   item.Description,
				Copyright:     "",
				CopyrightLink: "",
				SourceLink:    item.Links.Self,
				Author:        item.User.Username,
				Location:      "",
			})
		}
	}
	return sources.ArchiveResponse{
		Items:    items,
		End:      false, // 无限滚动
		Current:  param.Current,
		PageSize: param.PageSize,
	}, nil
}
