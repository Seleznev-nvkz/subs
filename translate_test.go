package main

import (
	"fmt"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestTranslate(t *testing.T) {
	ts := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "application/json")
		fmt.Fprintln(w, `{"text":["staff"]}`)
	}))

	translator = &Translator{url: ts.URL}

	word := &Word{Word: "stuff"}
	word.translate()
	assert.Equal(t, word.Translation, "staff")
}
