package scrollsvc

import (
	"encoding/json"
	"io"

	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/skeleton/models"
)

func (s *Scroller) processBuffer() {
	defer s.wg.Done()

	level.Info(s.logger).Log("msg", "transport start", "indexFrom")

	var (
		size            = 1000
		scroll          = s.dbES.GetScroll()
		processedAmount = 0
		yourStructs     = make([]*models.YourStruct, 0)
	)

	total, err := s.dbES.GetTotal()
	if err != nil {
		level.Error(s.logger).Log("msg", "get total failure", "err", err)
		return
	}

	if total == 0 {
		level.Warn(s.logger).Log("msg", "get total returned zero")
		return
	}

	level.Info(s.logger).Log("msg", "total elements to transfer", "total", total)

	for {
		results, err := scroll.Do(s.ctx)
		if err != nil {
			if err == io.EOF {
				s.eofChan <- true
			} else {
				s.errChan <- err
			}
		}

		select {
		default:
			for _, hit := range results.Hits.Hits {

				var yourStruct models.YourStruct
				err = json.Unmarshal(*hit.Source, &yourStruct)
				if err != nil {
					// Внимание!!! Возможно для вас критично обработать все записи без исключений
					level.Error(s.logger).Log("msg", "unmarshall error", "indexNameFrom", "err", err)
					continue
				}

				yourStructs = append(yourStructs, &yourStruct)

				if len(yourStructs) >= size {
					err = s.dbES.BulkUpsert(yourStructs)
					if err != nil {
						s.errChan <- err
					}

					yourStructs = yourStructs[:0]
				}

				processedAmount++
				if processedAmount%10000 == 0 {
					level.Info(s.logger).Log("msg", "processed", "amount", processedAmount, "total", total)
				}

			}
		case err = <-s.errChan:
			level.Error(s.logger).Log("msg", "transport index failure", "err", err)

			err = s.dbES.Reconnect()
			if err != nil {
				level.Error(s.logger).Log("msg", "reconnect failure", "err", err)
				return
			}

		case <-s.eofChan:
			err = s.dbES.BulkUpsert(yourStructs)
			if err != nil {
				s.errChan <- err
			}

			yourStructs = yourStructs[:0]

			level.Info(s.logger).Log("msg", "transport success end", "indexNameFrom")

			return
		}
	}
}
