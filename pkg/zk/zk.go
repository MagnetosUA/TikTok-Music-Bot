package zk

import (
	"fmt"
	config2 "github.com/MagnetosUA/TikTok-Music-Bot/pkg/config"
	"github.com/PuerkitoBio/goquery"
	"golang.org/x/net/publicsuffix"
	"net/http"
	"net/http/cookiejar"
	"strings"
	"time"
)

type ZK struct {
	HostURL     string
	SearchURL   string
	DownloadURL string
}

func NewZK(config *config2.Config) *ZK {
	return &ZK{
		HostURL:     config.ZkHostURL,
		SearchURL:   config.ZkSearchURL,
		DownloadURL: config.ZkDownloadURL,
	}
}

func (zk *ZK) DownloadFullSongZK(message string) ([]string, error) {
	client, err := zk.setupClient()
	if err != nil {
		return nil, err
	}

	if response, err := http.Get(zk.HostURL); err != nil {
		return nil, err
	} else {
		cookie := response.Header["Set-Cookie"][2]
		keywordsList := strings.Split(message, " ")

		var keywords string

		for i, word := range keywordsList {
			if word == "оригинальный" || word == "звук" {
				continue
			}

			keywords += word
			if i != len(keywordsList)-1 {
				keywords += "+"
			}
		}

		request := zk.SearchURL + keywords

		reader := strings.NewReader("my_request")
		newRequest, err := http.NewRequest("GET", request, reader)
		if err != nil {
			return nil, err
		}

		newRequest.Header.Add("Cookie", cookie)

		// Send HTTP request and move the response to the variable
		res, err := client.Do(newRequest)
		if err != nil {
			return nil, err
		}

		defer res.Body.Close()
		if res.StatusCode != 200 {
			return nil, fmt.Errorf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			return nil, err
		}

		var band []string

		songsList := doc.Find(".whb_gr_r")
		songsList.Find(".song-xl").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			attr, _ := s.Attr("data-play")
			text := zk.DownloadURL + attr

			band = append(band, text)
		})

		return band, nil
	}

	return nil, nil
}

// SET COOKIEJAR
func (zk *ZK) setupClient() (client *http.Client, err error) {
	// Set cookiejar options
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	// Create new cookiejar for holding cookies
	jar, err := cookiejar.New(&options)
	if err != nil {
		return nil, err
	}

	// Create new http client with predefined options
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 60,
	}

	return client, nil
}
