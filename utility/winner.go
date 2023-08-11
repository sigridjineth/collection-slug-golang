package utility

import (
	"errors"
)

// Winner takes a slice of strings and returns the string that appears most frequently.
// It returns an error if there's a tie for the most frequent string or if no strings are found.
func Winner(str []string) (string, error) {
	slugPoints := make(map[string]int)

	for _, slug := range str {
		slugPoints[slug]++
	}

	var maxScore int
	for _, score := range slugPoints {
		if score > maxScore {
			maxScore = score
		}
	}

	var winners []string
	for slug, score := range slugPoints {
		if score == maxScore {
			winners = append(winners, slug)
		}
	}

	if len(winners) > 1 {
		return "", errors.New("unable to resolve conclusively")
	}

	if len(winners) == 0 {
		return "", errors.New("not found")
	}

	return winners[0], nil
}
