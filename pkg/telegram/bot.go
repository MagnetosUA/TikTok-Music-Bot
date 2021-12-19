package telegram

import (
	"bufio"
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/PuerkitoBio/goquery"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"golang.org/x/net/publicsuffix"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"os/exec"
	"regexp"
	"strconv"
	"strings"
	"time"
)

/////
///// SET COOKIEJAR
/////
func setupClient() (client *http.Client) {
	// Set cookiejar options
	options := cookiejar.Options{
		PublicSuffixList: publicsuffix.List,
	}

	// Create new cookiejar for holding cookies
	jar, _ := cookiejar.New(&options)

	// Create new http client with predefined options
	client = &http.Client{
		Jar:     jar,
		Timeout: time.Second * 60,
	}

	return client
}

func DownloadFullSongZK(message string) []string {
	client := setupClient()
	request := "https://z2.fm"
	if response, err := http.Get(request); err != nil {

	} else {
		cookie := response.Header["Set-Cookie"][2]

		keywordsList := strings.Split(message, " ")

		var keywords string

		//if strings.Contains(message, "оригинальный звук") {
		//
		//}

		for i, word := range keywordsList {
			if word == "оригинальный" || word == "звук" {
				continue
			}
			//if i > 1 {
			//	continue
			//}
			keywords += word
			if i != len(keywordsList)-1 {
				keywords += "+"
			}
		}

		request := "https://z2.fm/mp3/search?keywords=" + keywords

		reader := strings.NewReader("my_request")
		newRequest, err := http.NewRequest("GET", request, reader)
		if err != nil {
			return nil
		}

		newRequest.Header.Add("Cookie", cookie)

		// Send HTTP request and move the response to the variable
		res, err := client.Do(newRequest)
		if err != nil {
			log.Fatal(err)
		}

		defer res.Body.Close()
		if res.StatusCode != 200 {
			log.Fatalf("status code error: %d %s", res.StatusCode, res.Status)
		}

		// Load the HTML document
		doc, err := goquery.NewDocumentFromReader(res.Body)
		if err != nil {
			log.Fatal(err)
		}

		var band []string
		// Find the review items
		//doc.Find(".song-xl").Each(func(i int, s *goquery.Selection) {
		//	// For each item found, get the band and title
		//	attr, _ := s.Attr("data-play")
		//	if attr != "" {}
		//})

		songsList := doc.Find(".whb_gr_r")
		songsList.Find(".song-xl").Each(func(i int, s *goquery.Selection) {
			// For each item found, get the band and title
			attr, _ := s.Attr("data-play")
			text := "https://z2.fm/download/" + attr

			band = append(band, text)

			//if is {}

			//title := s.Find("i").Text()
			//fmt.Printf("Review %d: %s - %s\n", i, band, title)
		})

		return band
	}
	return nil
}

type Bot struct {
	bot *tgbotapi.BotAPI
}

func NewBot(bot *tgbotapi.BotAPI) *Bot {
	return &Bot{bot: bot}
}

func (b *Bot) Start() error {
	log.Printf("Authorized on account %s", b.bot.Self.UserName)

	updates, err := b.initUpdatesChannel()
	if err != nil {
		return err
	}

	b.handleUpdates(updates)

	return nil
}

type ResponseTikTok struct {
	Html string `json:"html"`
}

func (b *Bot) handleUpdates(updates tgbotapi.UpdatesChannel) {
	for update := range updates {
		if update.Message == nil { // ignore any non-Message Updates
			continue
		}

		var responseTikTok ResponseTikTok

		var message string

		var url string

		if response, err := http.Get(update.Message.Text); err == nil {
			url = "https://www.tiktok.com/oembed?url=" + response.Request.URL.Scheme + "://" + response.Request.URL.Host + response.Request.URL.Path
			//message = url
			if response, err := http.Get(url); err != nil {

			} else {
				defer response.Body.Close()

				//Считываем ответ
				contents, err := ioutil.ReadAll(response.Body)
				if err != nil {
					log.Fatal(err)
				}

				if err = json.Unmarshal(contents, &responseTikTok); err != nil {

				}

				message = url
				//if responseTikTok.Html != "" {
				//
				//	//message = responseTikTok.Html
				//	re := regexp.MustCompile("♬(.*)(\" )")
				//	match := re.FindStringSubmatch(responseTikTok.Html)
				//	//fmt.Println(responseTikTok.Html)
				//
				//	if len(match) > 1 {
				//		message = match[1]
				//	}
				//}

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

		songName, err := getSongName(message)
		if err != nil {
			return
		}

		if songName == "" {
			songName = "Bad request!"
		} else {
			//fmt.Println("\n\n\n", "Song name: ", message, "\n\n\n")
		}

		var messages []string

		if len(songName) < 2 {
			songName = "Song is not found"
		} else {
			messages = DownloadFullSongZK(songName)
		}

		fmt.Println("SOUND NAME: ", songName)

		//DownloadFile(songName+"mp3", messages[0])

		//messages = append(messages, message)

		fmt.Sprintf("MESSAGES.len: %v", len(messages))

		if messages != nil {
			for i, mes := range messages {
				if i >= 2 {
					break
				}

				//func NewAudioShare(chatID int64, fileID string) AudioConfig {
				//	return AudioConfig{
				//	BaseFile: BaseFile{
				//	BaseChat:    BaseChat{ChatID: chatID},
				//	FileID:      fileID,
				//	UseExisting: true,
				//},
				//}
				//}

				fmt.Println("SOUND URL: ", mes)

				//response, err := http.Get(mes)
				//if err != nil {
				//
				//}

				//link := "<a href=\"google.com\">Link</a>"
				//audio := tgbotapi.NewInlineQueryResultAudio("123", mes, songName)
				textMessage := tgbotapi.NewMessage(update.Message.Chat.ID, songName)
				msg := tgbotapi.NewMessage(update.Message.Chat.ID, mes+"\n")
				//msg := tgbotapi.NewAudioUpload(update.Message.Chat.ID, mes)
				//msg := tgbotapi.NewMessage(update.Message.Chat.ID, mes)

				b.bot.Send(textMessage)
				b.bot.Send(msg)

				break
			}
		} else {
			textMessage := tgbotapi.NewMessage(update.Message.Chat.ID, "Song is not found")
			//msg := tgbotapi.NewMessage(update.Message.Chat.ID, message)
			b.bot.Send(textMessage)
			//b.bot.Send(msg)
		}

	}
}

func (b *Bot) initUpdatesChannel() (tgbotapi.UpdatesChannel, error) {
	u := tgbotapi.NewUpdate(0)
	u.Timeout = 60

	return b.bot.GetUpdatesChan(u)
}

func getSongName(originalSoundURL string) (songName string, err error) {

	playURL, err := getPlayURL(originalSoundURL)
	if err != nil {
		fmt.Println(err.Error())
	}

	filePath := "music.mp3"

	_, err = DownloadFile(filePath, playURL)
	if err != nil {
		panic(err)
	}

	fmt.Println("Downloaded: " + playURL)

	sourceSong := "music"
	sourceSongName := fmt.Sprintf("%v.mp3", sourceSong)
	outputShortSongName := fmt.Sprintf("%v_short.mp3", sourceSong)
	outputSongName := fmt.Sprintf("%v.wav", sourceSong)

	exec.Command("ffmpeg", "-y", "-ss", "3", "-i", sourceSongName, "-t", "5", outputShortSongName).Output()

	exec.Command("ffmpeg", "-y", "-i", outputShortSongName, "-acodec", "pcm_s16le", "-ac", "1", "-ar", "48000", outputSongName).Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	filePath = outputSongName

	//Open file on disk.
	f, _ := os.Open(filePath)

	// Read entire JPG into byte slice.
	reader := bufio.NewReader(f)
	content, _ := ioutil.ReadAll(reader)

	// Encode as base64.
	encoded := base64.StdEncoding.EncodeToString(content)

	song, err := shazamDetect(encoded)
	if err != nil {
		fmt.Println(err.Error())
		return
	}

	fmt.Println(song.Track.Subtitle, " ", song.Track.Title)

	return song.Track.Subtitle + " " + song.Track.Title, nil
}

type ShazamResponse struct {
	Track ShazamTrack `json:"track"`
}

type ShazamTrack struct {
	Title    string `json:"title"`
	Subtitle string `json:"subtitle"`
}

func shazamDetect(songBytes string) (ShazamResponse, error) {

	url := "https://shazam.p.rapidapi.com/songs/detect"

	payload := strings.NewReader(songBytes)

	req, _ := http.NewRequest("POST", url, payload)

	req.Header.Add("x-rapidapi-host", "shazam.p.rapidapi.com")
	req.Header.Add("x-rapidapi-key", "0dcc409e08msh55fe6be19bff0bcp192cf3jsn0d091d334fda")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	var response ShazamResponse

	err := json.Unmarshal(body, &response)
	if err != nil {
		return response, err
	}

	return response, nil
}

// referer header is not necessary while accessing API, but is must when
// downloading videos. The same headers and cookies must be used both when
// access API and downloading videos.
var defaultHeaders = map[string]string{
	"Accept":     "text/html,application/xhtml+xml,application/xml;q=0.9,image/avif,image/webp,image/apng,*/*;q=0.8,application/signed-exchange;v=b3;q=0.9",
	"User-Agent": "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/88.0.4324.182 Safari/537.36",
}

var defaultCookies = make(map[string]string)

func GetFirstTimeCookies() (c map[string]string, err error) {
	url := "https://www.tiktok.com/"
	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return
	}

	for name, value := range defaultHeaders {
		req.Header.Set(name, value)
	}

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	c = make(map[string]string)
	// Ref: https://stackoverflow.com/a/53023010
	// Ref: https://gist.github.com/rowland/984989
	for _, cookie := range resp.Cookies() {
		c[cookie.Name] = cookie.Value
	}
	return
}

// SendHttpRequest sends HTTP request.
func SendHttpRequest(url, method string, cookies, headers map[string]string) (b []byte, err error) {
	req, err := http.NewRequest(method, url, nil)
	if err != nil {
		return
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
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode != 200 {
		err = errors.New(url +
			"\nresp.StatusCode: " + strconv.Itoa(resp.StatusCode))
		return
	}

	return ioutil.ReadAll(resp.Body)
}

// GetHeaders returns HTTP request headers.
func GetHeaders() map[string]string {
	return defaultHeaders
}

func getPlayURL(originalSoundURL string) (string, error) {
	var playURL string

	fmt.Println("\n\n", originalSoundURL, "\n\n")
	defaultCookies, err := GetFirstTimeCookies()

	//url := "https://www.tiktok.com/music/%D0%BE%D1%80%D0%B8%D0%B3%D0%B8%D0%BD%D0%B0%D0%BB%D1%8C%D0%BD%D1%8B%D0%B9-%D0%B7%D0%B2%D1%83%D0%BA-7021957590359886593"
	url := "https://www.tiktok.com" + originalSoundURL
	b, err := SendHttpRequest(url, http.MethodGet, defaultCookies, GetHeaders())
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

	println(playURL)

	return playURL, nil
}

// DownloadFile will download a url to a local file. It's efficient because it will
// write as it downloads and not load the whole file into memory.
func DownloadFile(filepath string, url string) (*os.File, error) {
	var file *os.File
	// Get the data
	resp, err := http.Get(url)
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
