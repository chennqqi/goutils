package httputil

import (
	"encoding/json"
	"net/http"
)

func ReadJson(resp *http.Response, v interface{}) error {
	defer resp.Body.Close()
	d := json.NewDecoder(resp.Body)
	return d.Decode(v)
}
