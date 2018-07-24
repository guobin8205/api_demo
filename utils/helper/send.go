package helper

import (
	"net/http"
	"io/ioutil"
	"strings"
)

func HttpGet(url string) string{
	resp, err := http.Get(url)
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}

func HttpPost(url,params string) string{
	resp, err := http.Post(url, "application/x-www-form-urlencoded",
	strings.NewReader(params))

	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return ""
	}
	return string(body)
}