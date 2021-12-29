package elasticsearch

import (
	"time"

	"github.com/olivere/elastic"
)

func (e *Engine) GetScroll() *elastic.ScrollService {
	return e.client.Scroll(yourIndexName).Size(1000).KeepAlive("10m").Retrier(newElasticRetrier(time.Second*3, 10))
}

func (e *Engine) GetTotal() (int64, error) {
	return e.client.Count(yourIndexName).Do(e.ctx)
}
