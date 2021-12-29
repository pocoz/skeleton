package elasticsearch

import (
	"context"
	"fmt"
	"net/http"
	"time"

	"github.com/olivere/elastic"
)

const (
	retryExponentialFirstTimeInterval = 10 * time.Millisecond
	retryExponentialMaxTimeInterval   = 32 * time.Second
)

type retrier struct {
	backoff    elastic.Backoff
	retryCount int
}

func (e *retrier) Retry(ctx context.Context, retry int, req *http.Request, resp *http.Response, err error) (time.Duration, bool, error) {
	fmt.Printf("elastic Retry %d, err: %q\n", retry, err)

	// Stop after 5 retries
	if retry >= e.retryCount {
		return 0, false, nil
	}

	// Let the backoff strategy decide how long to wait and whether to stop
	wait, stop := e.backoff.Next(retry)

	return wait, stop, nil
}

func newElasticRetrier(interval time.Duration, count int) *retrier {
	return &retrier{
		backoff:    elastic.NewConstantBackoff(interval),
		retryCount: count,
	}
}
