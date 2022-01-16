package shazam

import (
	"encoding/json"
	"github.com/MagnetosUA/TikTok-Music-Bot/pkg/config"
	"io/ioutil"
	"net/http"
	"strings"
)

type ShazamInterfce interface {
	ShazamDetect(songBytes string) (ShazamAPIResponse, error)
}

type Shazam struct {
	ResourceURL string
	Host        string
	Key         string
}

func NewShazam(config config.Config) *Shazam {
	return &Shazam{
		ResourceURL: config.ShazamResourceURL,
		Host:        config.ShazamHost,
		Key:         config.ShazamKey,
	}
}

type ShazamAPIResponse struct {
	Track ShazamTrack `json:"track"`
}

type ShazamTrack struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

func (s *Shazam) ShazamDetect(songBytes string) (ShazamAPIResponse, error) {
	payload := strings.NewReader(songBytes)

	req, err := http.NewRequest("POST", s.ResourceURL, payload)
	if err != nil {
		return ShazamAPIResponse{}, err
	}

	req.Header.Add("x-rapidapi-host", s.Host)
	req.Header.Add("x-rapidapi-key", s.Key)

	res, err := http.DefaultClient.Do(req)
	if err != nil {
		return ShazamAPIResponse{}, err
	}

	defer res.Body.Close()
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		return ShazamAPIResponse{}, err
	}

	var response ShazamAPIResponse

	err = json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}
