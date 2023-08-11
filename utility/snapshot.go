package utility

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
)

type ArchivedSnapshots struct {
	Closest struct {
		Available bool   `json:"available"`
		URL       string `json:"url"`
	} `json:"closest"`
}

type AvailabilityResponse struct {
	ArchivedSnapshots *ArchivedSnapshots `json:"archived_snapshots"`
}

// FetchMaybeAvailableSnapshotUrl fetches the available snapshot URL for a given archive URL
func FetchMaybeAvailableSnapshotUrl(maybeArchiveUrl string) (string, error) {
	url := fmt.Sprintf("https://archive.org/wayback/available?url=%s", maybeArchiveUrl)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}

	var maybeAvailability AvailabilityResponse
	err = json.Unmarshal(body, &maybeAvailability)
	if err != nil {
		return "", err
	}

	if maybeAvailability.ArchivedSnapshots == nil || !maybeAvailability.ArchivedSnapshots.Closest.Available {
		return "", nil
	}

	return maybeAvailability.ArchivedSnapshots.Closest.URL, nil
}
