package client

import (
	"errors"
	"regexp"
	"sigridjineth/collection-slug-golang/network"
	"sigridjineth/collection-slug-golang/utility"
)

// FetchContractAddress retrieves the contract address for a given collection slug.
// It uses a redundancy parameter to determine how many archive URLs to fetch.
func FetchContractAddress(collectionSlug string, redundancy int) (string, error) {
	if redundancy == 0 {
		redundancy = network.DEFAULT_REDUNDANCY
	}

	cdxUri := "opensea.io/collection/" + collectionSlug
	archiveUrls, err := fetchArchiveUrls(CDX{CdxURI: cdxUri, Count: redundancy})
	if err != nil {
		return "", err
	}

	for _, archiveUrl := range archiveUrls {
		maybeSnapshotUrl, err := utility.FetchMaybeAvailableSnapshotUrl(archiveUrl)
		if err != nil || maybeSnapshotUrl == "" {
			continue
		}

		htmlCh := Text(maybeSnapshotUrl)
		html := <-htmlCh

		re := regexp.MustCompile(`\b0x[a-f0-9]{40}\b`)
		addresses := re.FindAllString(html, -1)

		if winnerAddress, err := utility.Winner(addresses); err == nil {
			return winnerAddress, nil
		}
	}

	return "", errors.New("unable to determine closest snapshot url")
}
