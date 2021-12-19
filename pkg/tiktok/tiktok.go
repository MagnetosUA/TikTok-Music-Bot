package tiktok

import (
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
)

func WikipediaAPI(request string) (answer []string) {
	//Отправляем запрос
	if response, err := http.Get(request); err != nil {
		return nil
	} else {
		defer response.Body.Close()

		//Считываем ответ
		contents, err := ioutil.ReadAll(response.Body)
		if err != nil {
			log.Fatal(err)
		}

		if len(contents) != 0 {
		}
	}

	return nil
}

//Конвертируем запрос для использование в качестве части URL
func UrlEncoded(str string) (string, error) {
	u, err := url.Parse(str)
	if err != nil {
		return "", err
	}
	return u.String(), nil
}
