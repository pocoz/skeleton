package mssql

import (
	"database/sql"
	"encoding/json"

	"github.com/go-kit/kit/log/level"

	"github.com/pocoz/skeleton/models"
)

func (e *Engine) GetOffersByItems(items []*models.YourStruct) (*models.OffersSyncMap, []*models.Offer, error) {
	values, err := json.Marshal(items)
	if err != nil {
		return nil, nil, err
	}

	ctx, query := e.getQuery(e.ctx, getOffersByItems)

	rows, err := e.db.QueryContext(
		ctx,
		query,
		sql.Named("values", string(values)),
	)
	if err != nil {
		return nil, nil, err
	}
	defer rows.Close()

	index := 0
	offers := make([]*models.Offer, 0)
	offersSyncMap := &models.OffersSyncMap{
		M: make(map[string]*models.Offer),
	}
	for rows.Next() {
		index++
		if index%100000 == 0 {
			level.Info(e.logger).Log("msg", "getOffersByItems in work", "amount", index)
		}

		offer := new(models.Offer)
		err = rows.Scan(
			&offer.ItemID,
			&offer.ShortWebName,
			&offer.Level,
			&offer.Category1,
			&offer.Category2,
			&offer.Category3,
			&offer.Category4,
			&offer.Category5,
			&offer.Category6,
			&offer.Category7,
			&offer.Category8,
			&offer.Category9,
			&offer.Category10,
		)
		if err != nil {
			level.Info(e.logger).Log("msg", "getOffersByItems scan rows failed", "err", err)
			continue
		}
		offers = append(offers, offer)
	}

	level.Info(e.logger).Log("msg", "getOffersByItems done", "amount", index)

	return offersSyncMap, offers, nil
}
