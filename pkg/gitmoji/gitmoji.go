package gitmoji

import (
	"encoding/json"
	"io/ioutil"
	"net/http"
	"sort"
)

func Fetch() (gitmojis []string, err error) {
	url := `https://raw.githubusercontent.com/carloscuesta/gitmoji/master/src/data/gitmojis.json`
	res, getErr := http.Get(url)
	if getErr != nil {
		return nil, getErr
	}
	body, readErr := ioutil.ReadAll(res.Body)
	if readErr != nil {
		return nil, readErr
	}

	gitmojiResponse := struct {
		Gitmojis []struct {
			Emoji string `json:"emoji"`
		} `json:"gitmojis"`
	}{}

	jsonErr := json.Unmarshal(body, &gitmojiResponse)
	if jsonErr != nil {
		return nil, jsonErr
	}
	for _, gm := range gitmojiResponse.Gitmojis {
		gitmojis = append(gitmojis, gm.Emoji)
	}
	sort.Strings(gitmojis)
	return
}
