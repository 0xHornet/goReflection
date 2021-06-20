package urlsource

import (
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/deitrix/xss-scraper/scraper"
)

type File struct {
	Filename string
}

func (s *File) GetURLS() (urls []string, err error) {
	data, err := ioutil.ReadFile(s.Filename)
	if err != nil {
		return nil, fmt.Errorf("error reading file: %v", err)
	}

	if err := json.Unmarshal(data, &urls); err != nil {
		return nil, fmt.Errorf("error decoding file contents: %v", err)
	}

	return urls, nil
}

func (s *File) StoreURLs(src scraper.URLSource) (err error) {
	urls, err := src.GetURLS()
	if err != nil {
		return fmt.Errorf("getting urls from source: %v", err)
	}

	bs, err := json.MarshalIndent(urls, "", "  ")
	if err != nil {
		return fmt.Errorf("failed to encode urls to json: %v", err)
	}

	if err := ioutil.WriteFile(s.Filename, bs, 0644); err != nil {
		return fmt.Errorf("failed to write to file: %v", err)
	}

	return nil
}
