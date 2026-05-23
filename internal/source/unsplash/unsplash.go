package unsplash

import (
	"daily-wallpaper/internal/source"
	"github.com/anaskhan96/soup"
	"log"
	"strings"
)

const (
	UrlPrefix = "https://unsplash.com/napi/topics/wallpapers/photos?%s"
)

type Source struct{}

func (u Source) GetToday() (source.TodayResponse, error) {
	res, err := soup.Get("https://unsplash.com/")
	if err != nil {
		return source.TodayResponse{}, err
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
						urlHS = source.GetSafeUrl(attrs["srcset"])
					} else if strings.Contains(attrs["media"], "(min-width: 2000px)") {
						url = source.GetSafeUrl(attrs["srcset"])
					} else if strings.Contains(attrs["media"], "(min-width: 200px)") {
						urlP = source.GetSafeUrl(attrs["srcset"])
					}
				}
				return source.ImageItem{
					Url:   url,
					UrlP:  urlP,
					UrlHS: urlHS,
				}, err
			}
		}
	}
	log.Println("查找节点失败")
	return source.TodayResponse{}, nil
}

func (u Source) GetArchive(param source.ArchiveParam) (source.ArchiveResponse, error) {
	payload := RequestPayload{
		Page:    param.Current,
		PerPage: param.PageSize,
	}
	items := make([]source.ImageItem, 0)
	var result Response
	err := source.DispatchGetRequest(UrlPrefix, payload, &result)
	if err == nil {
		for _, item := range result {
			items = append(items, source.ImageItem{
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
	return source.ArchiveResponse{
		Items:    items,
		End:      false, // 无限滚动
		Current:  param.Current,
		PageSize: param.PageSize,
	}, nil
}
