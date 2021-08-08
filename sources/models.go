package sources

type size struct {
	Height int    `json:"height"`
	Width  int    `json:"width"`
	Unit   string `json:"unit"`
}

type ImageItem struct {
	Name          string `json:"name"`
	Url           string `json:"url"`
	UrlP          string `json:"urlP"`  // 预览图片
	UrlHS         string `json:"urlHS"` // 高清图片
	Description   string `json:"description"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightLink"`
	SourceLink    string `json:"sourceLink"`
	Author        string `json:"author"`
	Location      string `json:"location"`
	RawSize       size   `json:"rawSize"` // 原始大小
}

type PaginationParam struct {
	Current  int `json:"current" form:"current"`
	PageSize int `json:"pageSize" form:"pageSize"`
}

type TodayResponse = ImageItem

type ArchiveResponse struct {
	Items    []ImageItem `json:"items"`
	End      bool        `json:"end"`
	Current  int         `json:"current"`
	PageSize int         `json:"pageSize"`
}

type ArchiveParam = PaginationParam

type Description struct {
	Name        string `json:"name"`
	Description string `json:"description"`
	Enabled     bool   `json:"enabled"`
}
