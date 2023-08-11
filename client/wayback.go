package client

import (
	"encoding/json"
	"fmt"
	"github.com/juju/ratelimit"
	"io/ioutil"
	"net/http"
	"sigridjineth/collection-slug-golang/utility"
	"sync"
)

type PromiseCallback struct {
	resolve func(result string)
	reject  func(error error)
}

type State map[string][]*PromiseCallback

var bottleneckWaybackMachine = ratelimit.NewBucketWithRate(4, 4) // Adjust rate as needed

// Text returns a channel with the content from the given URL using the Wayback function
func Text(url string) <-chan string {
	return Wayback()(url)
}

// JSON fetches content from the given URL using the Wayback function and decodes it into a JSON object
func JSON(url string) (interface{}, error) {
	textChannel := Wayback()(url)
	textContent := <-textChannel

	var result interface{}
	err := json.Unmarshal([]byte(textContent), &result)
	if err != nil {
		return nil, err
	}

	return result, nil
}

func Wayback() func(string) <-chan string {
	state := make(State)

	popQueue := func(key string) []*PromiseCallback {
		queue := state[key]
		delete(state, key)
		return queue
	}

	resolveQueue := func(key string, result *http.Response) error {
		text, err := ioutil.ReadAll(result.Body)
		if err != nil {
			return err
		}

		for _, callback := range popQueue(key) {
			callback.resolve(string(text))
		}
		return nil
	}

	rejectQueue := func(key string, err error) {
		for _, callback := range popQueue(key) {
			callback.reject(err)
		}
	}

	var mutex sync.Mutex

	return func(url string) <-chan string {
		ch := make(chan string)

		mutex.Lock()
		callback := &PromiseCallback{
			resolve: func(result string) { ch <- result },
			reject:  func(err error) { fmt.Println("Error:", err) },
		}

		queue := append(state[url], callback)
		state[url] = queue

		if len(queue) > 1 {
			mutex.Unlock()
			return ch
		}

		go func() {
			bottleneckWaybackMachine.Wait(1) // Wait for token
			result, err := utility.FetchWithExponentialBackoff(url)
			if err != nil {
				rejectQueue(url, err)
			} else {
				err := resolveQueue(url, result)
				if err != nil {
					rejectQueue(url, err)
				}
			}
		}()

		mutex.Unlock()
		return ch
	}
}
