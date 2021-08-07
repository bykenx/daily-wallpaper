package bing_source

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
)

const (
	UrlPrefix = "https://cn.bing.com/HPImageArchive.aspx"
)

func dispatchRequest(payload RequestPayload) []ImageItem {
	v, _ := query.Values(payload)
	println(v.Encode())
	res, err := http.Get(fmt.Sprintf("%s?%s", UrlPrefix, v.Encode()))
	if err != nil {
		return nil
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return nil
	}
	var resJson Response
	err = json.Unmarshal(data, &resJson)
	if err != nil {
		return nil
	}
	return resJson.Images
}