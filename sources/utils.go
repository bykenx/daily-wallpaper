package sources

import (
	"encoding/json"
	"fmt"
	"github.com/google/go-querystring/query"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

func DispatchGetRequest(url string, payload, result interface{}) error {
	v, _ := query.Values(payload)
	res, err := http.Get(fmt.Sprintf(url, v.Encode()))
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return err
	}
	err = json.Unmarshal(data, result)
	if err != nil {
		return err
	}
	return nil
}

func GetSafeUrl(s string) string {
	parts := strings.SplitN(s, "?", 2)
	if len(parts) > 1 {
		parts[1] = url.PathEscape(parts[1])
	}
	return strings.Join(parts, "?")
}