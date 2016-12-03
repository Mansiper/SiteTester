package main

import (
	"os"
	"bufio"
	"log"
	"time"
	"strings"
	"strconv"
	"net/http"
	"io/ioutil"
	"encoding/json"
	"net/http/cookiejar"
	"golang.org/x/net/publicsuffix"
)

type PageInfo struct{
	Url string				`json:"url"`
	Content []string	`json:"content"`
}
type Config struct {
	Login string			`json:"login"`
	Password string		`json:"password"`
	Pages []PageInfo	`json:"pages"`
}

//==================================================================================================

func main() {
	const confFile = "configuration.json"
	var (
		err error
		resp *http.Response
		data string
		tmStart, tmEnd time.Time
		strResults, fileName, errorStr string
		notFound bool
		config Config
	)

	//Loading configuration
	if _, err = os.Stat(confFile); err == nil {
		file, err := ioutil.ReadFile(confFile)
		if err = json.Unmarshal(file, &config); err != nil {
			log.Fatal(err)
		}
	} else if err != nil {
		log.Fatal(err)
	}

	//Prepare client with cookies
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}
	jar, err := cookiejar.New(&options)
	if err != nil { log.Fatal(err) }
	client := http.Client{Jar: jar}

	//Log in
	Login(&client, config.Login, config.Password)

	//Pages walking
	for k, v := range config.Pages {
		log.Println(strconv.Itoa(k+1) + ": " + v.Url)
		tmStart = time.Now()
		resp, err = client.Get(v.Url)
		tmEnd = time.Now()
		data = GetBody(resp)

		strResults = ""
		notFound = false
		//Page content checking
		for _, sv := range v.Content {
			if strings.Contains(data, sv) {
				strResults += "    \"" + sv + "\": ok\r\n"
			} else {
				strResults += "    \"" + sv + "\": not found\r\n"
				notFound = true
			}
		}
		if notFound { errorStr = "error_" } else { errorStr = "" }
		data = "<!--\r\nTime start: " + tmStart.String() +
			"\r\ntime end: " + tmEnd.String() +
			"\r\nURL: " + v.Url +
			"\r\nHTTP code: " + strconv.Itoa(resp.StatusCode) +
			"\r\nStrings:\r\n" + strResults +
			"-->\r\n" + data

		//Replace wrong symbols for saving file
		fileName = errorStr +
			strings.Replace(
				strings.Replace(
					strings.Replace(
						strings.Replace(
							strings.Replace(
								strings.Replace(
									strings.Replace(
										strings.Replace(
											strings.Replace(v.Url,
											"\\", "_", -1),
										"/", "_", -1),
									"*", "_", -1),
								":", "_", -1),
							"?", "_", -1),
						"\"", "_", -1),
					"<", "_", -1),
				">", "_", -1),
			"|", "_", -1) + ".html"

		//Saving
		ioutil.WriteFile(fileName, []byte(data), 0644)
	}

	//Waiting for enter
	log.Println("Well done. Press Enter")
	in := bufio.NewReader(os.Stdin)
	in.ReadString('\n')
}
