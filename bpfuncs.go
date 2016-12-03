package main

import (
	"log"
	"net/url"
	"strings"
	"net/http"
	"io/ioutil"
)

//--------------------------------------------------------------------------------------------------

func GetBody(resp *http.Response) string {
	data, err := ioutil.ReadAll(resp.Body)
	resp.Body.Close()
	if err != nil { return "" }
	return string(data)
}

func Login(client *http.Client, login string, password string) bool {
	const (
		siteLogin = "http://example.com/login"
		checkString = "Please login"
	)

	resp, err := (*client).PostForm(siteLogin, url.Values{
		"login": {login},
		"password": {password},
		"anotherLoginParam": {"param"},
	})
	if err != nil {
		log.Println(err)
		return false
	}
	data := GetBody(resp)
	if strings.Contains(data, checkString) { return false }
	return true
}
