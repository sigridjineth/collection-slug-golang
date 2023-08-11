package client

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type CDX struct {
	CdxURI string
	Count  int
}

func cdx(params CDX) string {
	return fmt.Sprintf("http://web.archive.org/cdx/search/cdx?url=%s&limit=%d&filter=mimetype:text/html&fl=original&output=json&status=200", params.CdxURI, params.Count)
}

func fetchArchiveUrls(params CDX) ([]string, error) {
	url := cdx(struct {
		CdxURI string
		Count  int
	}{params.CdxURI, params.Count})

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	var data []interface{}
	err = json.NewDecoder(resp.Body).Decode(&data)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		return nil, fmt.Errorf(`Expected data array, encountered "%v".`, data)
	}

	archiveURLs := make([]string, 0)
	uniqueURLs := make(map[string]bool)

	for _, e := range data {
		if url, ok := e.(string); ok {
			if !uniqueURLs[url] {
				archiveURLs = append(archiveURLs, url)
				uniqueURLs[url] = true
			}
		}
	}

	if len(archiveURLs) == 0 {
		return nil, fmt.Errorf(`unable to find an attempted archive url for "%s"`, params.CdxURI)
	}

	// Skip the first element if needed
	return archiveURLs[1:], nil
}
