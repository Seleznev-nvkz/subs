package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"log"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestExceptionView(t *testing.T) {
	db = newDB("/tmp/test_parse.db")
	defer os.Remove("/tmp/test_parse.db")

	ts := httptest.NewServer(getRoute())
	defer ts.Close()

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", ts.URL+"/api/words/something", nil)
	if err != nil {
		log.Fatal(err)
	}
	req.Header.Set("Content-Type", "application/json")
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}

	assert.Equal(t, resp.StatusCode, http.StatusNoContent)
	var exception Exceptions
	err = db.One("Word", "something", &exception)
	assert.Equal(t, exception.Word, "something")
	assert.Nil(t, err)
}

func TestUploadView(t *testing.T) {
	db = newDB("/tmp/test_parse.db")
	defer os.Remove("/tmp/test_parse.db")
	ts := httptest.NewServer(getRoute())
	defer ts.Close()

	var buffer bytes.Buffer
	writer := multipart.NewWriter(&buffer)
	part, _ := writer.CreateFormFile("file", "fileName")
	part.Write([]byte(`8 
                        00:01:55,213 --> 00:01:58,133
						- Can we help you?
						- My apologies.`))
	writer.Close()

	resp, err := http.Post(ts.URL+"/api/words/upload", writer.FormDataContentType(), &buffer)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, resp.StatusCode, 200)

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, `{"id":1,"name":"fileName","words":[{"word":"can","translation":""},{"word":"we","translation":""},{"word":"help","translation":""},{"word":"you","translation":""},{"word":"my","translation":""},{"word":"apology","translation":""}]}`, string(responseData))
}

func TestTranslateView(t *testing.T) {
	tss := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"text":["staff"]}`)
	}))
	defer tss.Close()

	// global
	translator = &Translator{url: tss.URL}

	ts := httptest.NewServer(getRoute())
	defer ts.Close()

	db.Save(&Subtitles{ID: 1, Name: "1", Words: []Word{{Word: "stuff"}}})
	resp, err := http.Get(ts.URL + "/api/subtitles/1/translate")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, http.StatusOK, resp.StatusCode)

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, `{"id":1,"name":"1","words":[{"word":"stuff","translation":"staff"}]}`, string(responseData))
}

func TestSubtitlesList(t *testing.T) {
	db = newDB("/tmp/test_parse.db")
	defer os.Remove("/tmp/test_parse.db")

	ts := httptest.NewServer(getRoute())
	defer ts.Close()

	db.Save(&Subtitles{Name: "1", Words: []Word{{"q", "w"}}})
	db.Save(&Subtitles{Name: "2", Words: []Word{{"a", "s"}}})

	resp, err := http.Get(ts.URL + "/api/subtitles/")
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 200, resp.StatusCode)

	responseData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, `[{"id":1,"name":"1","words":[{"word":"q","translation":"w"}]},{"id":2,"name":"2","words":[{"word":"a","translation":"s"}]}]`, string(responseData))
}

func TestSubtitlesDelete(t *testing.T) {
	db = newDB("/tmp/test_parse.db")
	defer os.Remove("/tmp/test_parse.db")

	ts := httptest.NewServer(getRoute())
	defer ts.Close()

	db.Save(&Subtitles{Name: "1", Words: []Word{{"q", "w"}}})
	db.Save(&Subtitles{Name: "2", Words: []Word{{"a", "s"}}})

	client := &http.Client{}
	req, err := http.NewRequest("DELETE", ts.URL+"/api/subtitles/2", nil)
	if err != nil {
		log.Fatal(err)
	}
	resp, err := client.Do(req)
	if err != nil {
		log.Fatal(err)
	}
	assert.Equal(t, 204, resp.StatusCode)

	var subtitles []Subtitles
	db.AllByIndex("ID", &subtitles)
	assert.Equal(t, 1, len(subtitles))
	assert.Equal(t, 1, subtitles[0].ID)
}
