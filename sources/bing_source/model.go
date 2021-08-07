package bing_source

type RequestPayload struct {
	Format string `url:"format"`
	Index  int          `url:"idx"`
	PageSize int          `url:"n"`
	Timestamp int64        `url:"nc"`
}

type ImageItem struct {
	StartDate     string        `json:"startdate"`
	FullStartDate string        `json:"fullstartdate"`
	EndDate       string        `json:"enddate"`
	Url           string        `json:"url"`
	UrlBase       string        `json:"urlbase"`
	Copyright     string        `json:"copyright"`
	CopyrightLink string        `json:"copyrightlink"`
	Title         string        `json:"title"`
	Quiz          string        `json:"quiz"`
	Wp            bool          `json:"wp"`
	Hsh           string        `json:"hsh"`
	Drk           int           `json:"drk"`
	Top           int           `json:"top"`
	Bot           int           `json:"bot"`
	Hs            []interface{} `json:"hs"`
}

type Response struct {
	Images []ImageItem `json:"images"`
}