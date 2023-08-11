package tests

import (
	"sigridjineth/collection-slug-golang/client"
	"sigridjineth/collection-slug-golang/network"
	"testing"
)

func Test_CollectionSlug(t *testing.T) {
	t.Run("ethereum (default)", func(t *testing.T) {
		t.Run("boredapeyachtclub", func(t *testing.T) {
			collectionSlug, err := client.FetchCollectionSlug("0xbc4ca0eda7647a8ab7c2061c2e118a18a936f13d", nil, network.ETHEREUM, nil)
			if err != nil {
				t.Fatalf("Expected no error, got %v", err)
			}
			if collectionSlug != "boredapeyachtclub" {
				t.Fatalf("Expected 'boredapeyachtclub', got %v", collectionSlug)
			}
		})
	})
}
