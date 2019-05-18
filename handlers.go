package main

import (
	"io"
	"io/ioutil"
	"regexp"
	"strings"

	"github.com/asdine/storm"
)

func parseString(input string) []string {
	res := strings.Replace(input, "\n", " ", -1)
	res = regexp.MustCompile(`([^a-zA-Z ]+)|(\s\s+)`).ReplaceAllString(
		regexp.MustCompile(`</?[^>]+(>|$)`).ReplaceAllString(res, ""), "")

	return strings.Fields(strings.ToLower(res))
}

func isException(word string) bool {
	_, err := db.GetBytes("Exceptions", word)
	if err == storm.ErrNotFound {
		return false
	}
	return true
}

func lemmatization(input []string) []string {
	result := make([]string, 0, len(input))
	keys := make(map[string]struct{})

	for _, word := range input {
		if _, value := keys[word]; !value {
			keys[word] = struct{}{}
			word = lemm.Lemma(word)
			if !isException(word) {
				result = append(result, word)
			}
		}
	}
	return result
}

func handleFile(file io.Reader) ([]string, error) {
	data, err := ioutil.ReadAll(file)
	if err != nil {
		return nil, err
	}

	arr := parseString(string(data))
	return lemmatization(arr), nil
}
