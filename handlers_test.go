package main

import (
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseString(t *testing.T) {
	data := `
00:10:45,421 --> 00:10:48,606    <i>Yeah, I just
engineered that solid piece, that's all.</i>`
	res := []string{"yeah", "i", "just", "engineered", "that", "solid", "piece", "thats", "all"}

	assert.Equal(t, res, parseString(data))
}

func TestLemmatization(t *testing.T) {
	db = newDB("/tmp/test_parse.db")
	defer os.Remove("/tmp/test_parse.db")

	data := []string{"yeah", "i", "exists", "engineered", "that", "solid", "feeling", "amazing", "exist"}
	res := []string{"yes", "i", "exist", "engineer", "that", "solid", "feel", "amaze"}

	assert.Equal(t, res, lemmatization(data))

	// duplicates
	data = []string{"yes", "yes"}
	res = []string{"yes"}
	assert.Equal(t, res, lemmatization(data))

	// except
	data = []string{"yes", "already", "i"}
	res = []string{"yes"}
	db.Save(&Exceptions{"already"})
	db.Save(&Exceptions{"i"})
	assert.Equal(t, res, lemmatization(data))
}

func TestModelRefresh(t *testing.T) {
	db = newDB("/tmp/test_parse.db")
	defer os.Remove("/tmp/test_parse.db")

	sub := Subtitles{ID: 1, Name: "1", Words: []Word{{Word: "stuff"}, {Word: "i"}}}
	db.Save(&Exceptions{"i"})
	db.Save(&sub)
	sub.refresh()

	var updated Subtitles
	db.One("ID", 1, &updated)
	assert.Equal(t, 1, len(updated.Words))
	assert.Equal(t, "stuff", updated.Words[0].Word)
}
