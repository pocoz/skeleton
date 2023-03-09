package elasticsearch

import (
	"encoding/json"
	"github.com/pocoz/skeleton/tools"

	"github.com/olivere/elastic"

	"github.com/pocoz/skeleton/models"
)

const (
	yourIndexName    = "your_index"
	yourIndexDoctype = "_doc" // предлагаю использовать один док тайп для индексов, как и в 7 версии
)

func (e *Engine) GetByID(id string) ([]*models.YourStruct, error) {
	query := elastic.NewTermQuery("_id", id)
	res, err := e.client.Search().
		Index(yourIndexName).
		Type(yourIndexDoctype).
		Query(query).
		Do(e.ctx)
	if err != nil {
		return nil, err
	}

	var items = make([]*models.YourStruct, 0)
	if res.TotalHits() > 0 {
		for _, hit := range res.Hits.Hits {
			var item models.YourStruct
			err = json.Unmarshal(*hit.Source, &item)
			if err != nil {
				return nil, err
			}
			items = append(items, &item)
		}
	}

	return items, nil
}

func (e *Engine) BulkUpsert(yourStructs []*models.YourStruct) error {
	var (
		maxBatchSize = 1000
		amount       = tools.GetAmount(len(yourStructs), maxBatchSize)
		start        = 0
		end          = 0
	)

	for i := 0; i < amount; i++ {
		start = maxBatchSize * i
		end = maxBatchSize * (i + 1)
		if end > len(yourStructs) {
			end = len(yourStructs)
		}

		part := yourStructs[start:end]
		backoffRetrier := elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(
			retryExponentialFirstTimeInterval,
			retryExponentialMaxTimeInterval,
		))

		bulkRequest := e.client.Bulk().Timeout("5m").Retrier(backoffRetrier)
		for _, p := range part {
			bulkIndex := elastic.NewBulkUpdateRequest().
				DocAsUpsert(true).
				Id(p.ItemID).
				Index(yourIndexName).
				Type(yourIndexDoctype).
				Doc(p)
			bulkRequest = bulkRequest.Add(bulkIndex)
		}

		_, err := bulkRequest.Do(e.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}

func (e *Engine) BulkUpdate(yourStructs []*models.YourStruct) error {
	var (
		maxBatchSize = 1000
		amount       = tools.GetAmount(len(yourStructs), maxBatchSize)
		start        = 0
		end          = 0
	)

	for i := 0; i < amount; i++ {
		start = maxBatchSize * i
		end = maxBatchSize * (i + 1)
		if end > len(yourStructs) {
			end = len(yourStructs)
		}

		part := yourStructs[start:end]
		backoffRetrier := elastic.NewBackoffRetrier(elastic.NewExponentialBackoff(
			retryExponentialFirstTimeInterval,
			retryExponentialMaxTimeInterval,
		))

		bulkRequest := e.client.Bulk().Timeout("5m").Retrier(backoffRetrier)
		for _, p := range part {
			bulkIndex := elastic.NewBulkUpdateRequest().
				Index(yourIndexName).
				Type(yourIndexDoctype).
				Id(p.ItemID).
				Doc(map[string]interface{}{
					"advanced_params":              p.AdvancedParams,
					"advanced_params_created_date": p.AdvancedParamsCreatedDate,
				})
			bulkRequest = bulkRequest.Add(bulkIndex)
		}

		_, err := bulkRequest.Do(e.ctx)
		if err != nil {
			return err
		}
	}
	return nil
}
