package sources

type ImageItem struct {
	Name          string `json:"name"`
	Url           string `json:"url"`
	Description   string `json:"description"`
	Copyright     string `json:"copyright"`
	CopyrightLink string `json:"copyrightLink"`
	SourceLink    string `json:"sourceLink"`
	Author        string `json:"author"`
	Location      string `json:"location"`
}

type PaginationParam struct {
	Current  int `json:"current"`
	PageSize int `json:"pageSize"`
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
