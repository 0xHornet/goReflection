package scraper

import (
	"encoding/hex"
	"math/rand"
	"net/url"
)

const fudgeURLParamsPerURL = 10

func FudgeURLs(urls, queryParams []string) (fudgedURLs []FudgedURL) {
	for _, u := range urls {
		fudgedURLs = append(fudgedURLs, FudgeURL(u, queryParams)...)
	}
	return fudgedURLs
}

func FudgeURL(rawURL string, queryParams []string) (fudgedURLs []FudgedURL) {
	u, _ := url.Parse(rawURL)

	// Save the query before we make any modifications so that we can revert back to it after each iteration.
	savedQuery := u.RawQuery

	// In batches, append a list of query params to the url, and then push the fudged url to fudgedURLs.
	for low := 0; low < len(queryParams); low += fudgeURLParamsPerURL {
		high := low + fudgeURLParamsPerURL
		if high > len(queryParams) {
			high = len(queryParams)
		}

		q := u.Query()
		params := make(map[string]string, high-low)
		for _, param := range queryParams[low:high] {
			value := randomHex(8)
			q.Set(param, value)
			params[param] = value
		}

		u.RawQuery = q.Encode()
		fudgedURLs = append(fudgedURLs, FudgedURL{
			URL:    u.String(),
			Params: params,
		})
		u.RawQuery = savedQuery
	}

	return fudgedURLs
}

type FudgedURL struct {
	URL    string            `json:"url"`
	Params map[string]string `json:"params"`
}

func randomHex(n int) string {
	bytes := make([]byte, n)
	rand.Read(bytes)
	return hex.EncodeToString(bytes)
}
