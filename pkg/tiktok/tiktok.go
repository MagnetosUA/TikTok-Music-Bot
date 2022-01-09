package tiktok

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/MagnetosUA/TikTok-Music-Bot/pkg/shazam"
	"io"
	"io/ioutil"
	"net/http"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
)

const hostURL = "https://www.tiktok.com/"

type TikTok struct {
	HostURL string
}

func NewTikTok() *TikTok {
	return &TikTok{
		HostURL: hostURL,
	}
}

type ResponseTikTok struct {
	Html string `json:"html"`
}

//var defaultCookies = make(map[string]string)

// referer header is not necessary while accessing API, but is must when
// downloading videos. The same headers and cookies must be used both when
// access API and downloading videos.
var defaultHeaders = map[string]string{
	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
}

func (tk *TikTok) GetFirstTimeCookies() (map[string]string, error) {
	req, err := http.NewRequest(http.MethodGet, tk.HostURL, nil)
	if err != nil {
		return nil, err
	}

	for name, value := range defaultHeaders {
		req.Header.Set(name, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(tk.HostURL +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return nil, err
	}

	c := make(map[string]string)
	// Ref: https://stackoverflow.com/a/53023010
	// Ref: https://gist.github.com/rowland/984989
	for _, cookie := range resp.Cookies() {
		c[cookie.Name] = cookie.Value
	}

	return c, nil
}

func (tk *TikTok) GetSongName(shazam *shazam.Shazam, originalSoundURL string) (string, error) {
	playURL, err := tk.getPlayURL(originalSoundURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	filePath := "music.mp3"

	fmt.Println("Play url:", playURL)

	_, err = tk.DownloadFile(filePath, playURL)
	if err != nil {
		return "", err
	}

	fmt.Println("Downloaded: " + playURL)

	sourceSong := "music"
	sourceSongName := fmt.Sprintf("%v.mp3", sourceSong)
	outputShortSongName := fmt.Sprintf("%v_short.mp3", sourceSong)
	outputSongName := fmt.Sprintf("%v.wav", sourceSong)

	exec.Command("ffmpeg", "-y", "-i", sourceSongName, "-t", "4", outputShortSongName).Output()

	exec.Command("ffmpeg", "-y", "-i", outputShortSongName, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "48000", outputSongName).Output()

	//Open file on disk.
	f, err := os.Open(outputSongName)
	if err != nil {
		return "", err
	}

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, err := ioutil.ReadAll(reader)
	if err != nil {
		return "", err
	}

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	song, err := shazam.ShazamDetect(encoded)
	if err != nil {
		fmt.Println(err.Error())
		return "", err
	}

	fmt.Println(song.Track.Subtitle, " ", song.Track.Title)

	return song.Track.Subtitle + " " + song.Track.Title, nil
}

func (tk *TikTok) getPlayURL(originalSoundURL string) (string, error) {
	var playURL string

	fmt.Println("\n\n", originalSoundURL, "\n\n")
	defaultCookies, err := tk.GetFirstTimeCookies()

	url := tk.HostURL + originalSoundURL
	b, err := tk.SendHttpRequest(url, http.MethodGet, defaultCookies, GetHeaders())
	if err != nil {
		return playURL, err
	}

	pattern := regexp.MustCompile(`"playUrl":"(.*?)","coverLarge"`)
	matches := pattern.FindSubmatch(b)
	//println(len(matches))
	if len(matches) != 2 {
		err = errors.New("trouble getting __NEXT_DATA__")
		return playURL, err
	}

	playURL = string(matches[1])
	playURL = strings.Replace(playURL, "\\u002F", "/", 5)

	return playURL, nil
}

// GetHeaders returns HTTP request headers.
func GetHeaders() map[string]string {
	return defaultHeaders
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func (tk *TikTok) DownloadFile(filepath string, url string) (*os.File, error) {
	var file *os.File
	// Get the data
	resp, err := http.Get(url)
	fmt.Println(resp)
	if err != nil {
		return file, err
	}
	defer resp.Body.Close()

	// Create the file
	file, err = os.Create(filepath)
	if err != nil {
		return file, err
	}
	defer file.Close()

	// Write the body to file
	_, err = io.Copy(file, resp.Body)
	return file, err
}

// SendHttpRequest sends HTTP request.
func (tk *TikTok) SendHttpRequest(url, method string, cookies, headers map[string]string) ([]byte, error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return nil, err
	}

	for name, value := range cookies {
		req.AddCookie(&http.Cookie{Name: name, Value: value})
	}

	for name, value := range headers {
		req.Header.Set(name, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return nil, err
	}

	return ioutil.ReadAll(resp.Body)
}

func (tk *TikTok) GetTikTokMessage(messageText string) (string, error) {
	var responseTikTok ResponseTikTok
	var message string

	if response, err := http.Get(messageText); err == nil {
		url := "https://www.tiktok.com/oembed?url=" + response.Request.URL.Scheme + "://" + response.Request.URL.Host + response.Request.URL.Path
		//message = url
		if response, err := http.Get(url); err != nil {
			return "", err
		} else {
			defer response.Body.Close()

			//Считываем ответ
			contents, err := ioutil.ReadAll(response.Body)
			if err != nil {
				return "", err
			}

			if err = json.Unmarshal(contents, &responseTikTok); err != nil {
				return "", err
			}

			message = url

			if responseTikTok.Html != "" {

				//message = responseTikTok.Html
				re := regexp.MustCompile("/music/(.*)♬")
				match := re.FindStringSubmatch(responseTikTok.Html)
				//fmt.Println(responseTikTok.Html)

				if len(match) > 1 {
					//fmt.Println(match)
					message = match[0][:len(match[0])-5]
					//inputFmt:=input[:len(input)-2]
				}
			}
		}
	}

	return message, nil
}
