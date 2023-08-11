package utility

import (
	"errors"
	"math"
	"net"
	"net/http"
	"time"
)

const (
	InitialBackoffTime = 5000
	BackoffMultiplier  = 1.25
)

var RetrievableFetchErrors = map[error]bool{
	net.UnknownNetworkError("ENOTFOUND"):    true,
	net.UnknownNetworkError("ECONNREFUSED"): true,
	net.UnknownNetworkError("ECONNRESET"):   true,
}

type ExponentialBackoffState struct {
	delay         int
	initializedAt int64
}

func initialState() ExponentialBackoffState {
	return ExponentialBackoffState{
		delay: 0,
	}
}

func stateWithBackoff(state ExponentialBackoffState, initialBackoffTime int, backoffMultiplier float64) ExponentialBackoffState {
	if state.delay == 0 {
		return ExponentialBackoffState{
			delay:         initialBackoffTime,
			initializedAt: time.Now().Unix(),
		}
	}
	return ExponentialBackoffState{
		delay: int(math.Ceil(float64(state.delay) * backoffMultiplier)),
	}
}

func FetchWithExponentialBackoff(url string) (*http.Response, error) {
	state := initialState()

	for {
		startedAt := time.Now().Unix()
		for now := startedAt; now-startedAt < int64(state.delay); now = time.Now().Unix() {
			time.Sleep(time.Duration(int64(state.delay)-(now-startedAt)) * time.Millisecond)
		}

		result, err := http.Get(url)

		if err != nil {
			// Check if error is in RetrievableFetchErrors
			var netErr net.Error
			if errors.As(err, &netErr) && RetrievableFetchErrors[netErr] {
				state = stateWithBackoff(state, InitialBackoffTime, BackoffMultiplier)
				continue
			}
			return nil, err
		}

		if result.StatusCode == 429 {
			state = stateWithBackoff(state, InitialBackoffTime, BackoffMultiplier)
			continue
		}

		return result, nil
	}
}
