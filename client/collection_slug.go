package client

import (
	"errors"
	"regexp"
	slugNetwork "sigridjineth/collection-slug-golang/network"
	"sigridjineth/collection-slug-golang/utility"
	"strings"
)

// FetchCollectionSlug retrieves the collection slug for a given contract address, network, token ID, and redundancy.
func FetchCollectionSlug(contractAddress string, tokenId *string, network slugNetwork.Network, redundancy *int) (string, error) {
	if redundancy == nil {
		defaultRedundancy := slugNetwork.DEFAULT_REDUNDANCY
		redundancy = &defaultRedundancy
	}

	if network == "" {
		network = slugNetwork.ETHEREUM
	}

	isOpenStore := utility.IsOpenStoreDeployment(contractAddress, utility.Network(network))
	if isOpenStore && tokenId == nil {
		return "", errors.New("to find the collectionSlug on the OPENSTORE contract, you must specify a tokenId")
	}

	tokenIdStr := "0"
	if tokenId != nil {
		tokenIdStr = *tokenId
	}

	cdxUri := "opensea.io/assets/" + string(network) + "/" + contractAddress + "/" + tokenIdStr
	archiveUrls, err := fetchArchiveUrls(CDX{CdxURI: cdxUri, Count: *redundancy})
	if err != nil {
		return "", err
	}

	for _, archiveUrl := range archiveUrls {
		maybeSnapshotUrl, err := utility.FetchMaybeAvailableSnapshotUrl(archiveUrl)
		if err != nil || maybeSnapshotUrl == "" {
			continue
		}

		someTextCh := Text(maybeSnapshotUrl)
		someText := <-someTextCh

		re := regexp.MustCompile(`(https):\/\/([\w_-]+(?:(?:\.[\w_-]+)+))([\w.,@?^=%&:\/~+#-]*[\w@?^=%&\/~+#-])`)
		slugs := re.FindAllString(someText, -1)

		var filteredSlugs []string
		for _, slug := range slugs {
			if strings.Contains(slug, slugNetwork.BASE_COLLECTION_URL) {
				slug = strings.Split(strings.Split(slug, slugNetwork.BASE_COLLECTION_URL)[1], "/")[0]
				slug = strings.Split(slug, "?")[0]
				filteredSlugs = append(filteredSlugs, slug)
			}
		}

		if winnerSlug, err := utility.Winner(filteredSlugs); err == nil {
			return winnerSlug, nil
		}
	}

	return "", errors.New("unable to determine closest snapshot url")
}
