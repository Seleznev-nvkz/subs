package main

import (
	"encoding/json"
	"log"
	"net/http"
	"net/url"
)

type Translator struct {
	apiKey string
	lang   string
	url    string
}

func newTranslator(key, lang string) *Translator {
	return &Translator{
		key,
		lang,
		"https://translate.yandex.net/api/v1.5/tr.json/translate",
	}
}

func (tr *Translator) translate(text string) string {
	builtParams := url.Values{"key": {tr.apiKey}, "lang": {tr.lang}, "text": {text}, "options": {"1"}}
	resp, err := http.PostForm(tr.url, builtParams)
	if err != nil {
		log.Println(err)
		return text
	}
	defer resp.Body.Close()

	data := make(map[string]interface{})
	if err = json.NewDecoder(resp.Body).Decode(&data); err != nil {
		log.Println(err)
		return text
	}
	words := data["text"].([]interface{})
	return words[0].(string)
}
