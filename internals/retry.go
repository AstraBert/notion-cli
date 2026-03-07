package internals

import (
	"fmt"
	"log"
	"net/http"
	"slices"
	"strconv"
	"time"
)

const MaxRetries int = 3
const RetryAfterHeader string = "Retry-After"
const DefaultRetryTime = 1

var RetriableErrorCodes [3]int = [3]int{429, 500, 503}

func RequestWithRetries(client *http.Client, request *http.Request) (*http.Response, error) {
	for range MaxRetries {
		response, err := client.Do(request)
		if err != nil {
			return nil, err
		}
		if !slices.Contains(RetriableErrorCodes[:], response.StatusCode) {
			return response, nil
		}
		retryAfter := response.Header.Get(RetryAfterHeader)
		var retryAfterTime time.Duration
		seconds, err := strconv.Atoi(retryAfter)
		switch err {
		case nil:
			retryAfterTime = time.Duration(seconds) * time.Second
		default:
			dt, errDt := time.Parse(time.RFC1123, retryAfter)
			if errDt != nil {
				retryAfterTime = time.Duration(1) * time.Second
			} else {
				retryAfterTime = time.Until(dt)
			}
		}
		log.Printf("Will retry after sleeping for %d seconds...\n", retryAfterTime/time.Second)
		time.Sleep(retryAfterTime)
	}
	return nil, fmt.Errorf("exceeded maximum number of retries: %d", MaxRetries)
}
