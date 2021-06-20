package urlsource

import (
	"encoding/json"
	"fmt"
	"net/http"
	"net/url"
	"path"
	"strings"
)

var excludedPathExtensions = []string{
	".css",
	".eot",
	".gif",
	".ico",
	".jpg",
	".pdf",
	".png",
	".svg",
	".ttf",
	".txt",
	".woff",
	".woff2",
	".jsp",
	".js",
}

var excludedPathExtensionsMap map[string]struct{}

const archiveURLTemplate = "http://web.archive.org/cdx/search/cdx?url=*.%s/*&output=json&fl=original&collapse=urlkey&filter=mimetype:text/html"

type ArchiveOrg struct {
	Domain string
}

func (s *ArchiveOrg) GetURLS() (urls []string, err error) {
	u := fmt.Sprintf(archiveURLTemplate, s.Domain)

	resp, err := http.Get(u)
	if err != nil {
		return nil, fmt.Errorf("failed to decode json body: %v", err)
	}
	defer resp.Body.Close()

	var decodedURLs [][]string
	if err := json.NewDecoder(resp.Body).Decode(&decodedURLs); err != nil {
		return nil, fmt.Errorf("failed to decode json body: %v", err)
	}

	urls = make([]string, 0, len(decodedURLs))
	for _, v := range decodedURLs {
		p, err := url.Parse(v[0])
		if err != nil {
			//log.Printf("error parsing url '%s': %v", v[0], err)
			continue
		}

		ext := strings.ToLower(path.Ext(p.Path))
		if _, ok := excludedPathExtensionsMap[ext]; ok {
			continue
		}

		urls = append(urls, v[0])
	}

	return urls[1:], nil
}

func init() {
	excludedPathExtensionsMap = make(map[string]struct{}, len(excludedPathExtensions))
	for _, ext := range excludedPathExtensions {
		excludedPathExtensionsMap[ext] = struct{}{}
	}
}
