package scraper

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"strings"
	"sync"

	"github.com/cheggaaa/pb"
)

const concurrency = 32

type URLSource interface {
	GetURLS() (urls []string, err error)
}

func Scrape(urlSource URLSource, params []string) (err error) {
	urls, err := urlSource.GetURLS()
	if err != nil {
		return fmt.Errorf("failed to get urls from url source: %v", err)
	}

	fudgedURLs := FudgeURLs(urls, params)

	// Randomly shuffle the list of urls.
	rand.Shuffle(len(fudgedURLs), func(i, j int) {
		fudgedURLs[i], fudgedURLs[j] = fudgedURLs[j], fudgedURLs[i]
	})

	jobs := makeFudgedURLChan(fudgedURLs)
	res := make(chan ScrapeResult)

	var wg sync.WaitGroup
	wg.Add(concurrency)
	for i := 0; i < concurrency; i++ {
		go func() {
			defer wg.Done()
			for r := range processFudgedURLs(jobs) {
				res <- r
			}
		}()
	}

	go func() {
		wg.Wait()
		close(res)
	}()

	progress := pb.StartNew(len(fudgedURLs))

	for r := range res {
		progress.Increment()
		if len(r.FoundParams) != 0 {
			fmt.Printf("Scraped url: %s, found reflected params: %s\n", r.URL.URL, strings.Join(r.FoundParams, ", "))
		}
	}

	progress.Finish()

	return nil
}

type ScrapeResult struct {
	URL         FudgedURL
	Err         error
	Status      int
	FoundParams []string
}

func GetResponse(url string) (body []byte, status int, err error) {
	resp, err := http.Get(url)
	if err != nil {
		return nil, 0, fmt.Errorf("error performing get request: %v", err)
	}
	defer resp.Body.Close()

	body, err = ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, 0, fmt.Errorf("error reading response body: %v", err)
	}

	return body, resp.StatusCode, nil
}

func SearchResponse(body []byte, params map[string]string) (foundParams []string) {
	for param, value := range params {
		if i := bytes.Index(body, []byte(value)); i != -1 {
			foundParams = append(foundParams, param)
		}
	}
	return foundParams
}

func makeFudgedURLChan(urls []FudgedURL) <-chan FudgedURL {
	c := make(chan FudgedURL)
	go func() {
		for _, url := range urls {
			c <- url
		}
		close(c)
	}()
	return c
}

func processFudgedURLs(recvURL <-chan FudgedURL) <-chan ScrapeResult {
	c := make(chan ScrapeResult)
	go func() {
		for {
			u, ok := <-recvURL
			if !ok {
				break
			}

			res := ScrapeResult{
				URL: u,
			}

			body, status, err := GetResponse(u.URL)
			if err != nil {
				res.Err = err
				continue
			}

			res.Status = status
			res.FoundParams = SearchResponse(body, u.Params)
			c <- res
		}
		close(c)
	}()
	return c
}
