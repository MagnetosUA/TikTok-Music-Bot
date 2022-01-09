package shazam

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"strings"
)

// should be moved to env variables
const (
	ShazamResourceURL = "https://shazam.p.rapidapi.com/songs/detect"
	Host              = "shazam.p.rapidapi.com"
	Key               = "0dcc409e08msh55fe6be19bff0bcp192cf3jsn0d091d334fda"
)

type ShazamInterfce interface {
	ShazamDetect(songBytes string) (ShazamAPIResponse, error)
}

type Shazam struct {
	ResourceURL string
	Host        string
	Key         string
}

func NewShazam(url string) *Shazam {
	return &Shazam{
		ResourceURL: url,
		Host:        Host,
		Key:         Key,
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
