package postgres

import (
	"time"

	"github.com/pocoz/skeleton/models"
)

const (
	statisticsInsertOne = uint8(iota)
	statisticsInsertMany
	statisticsGetOne
	statisticsGetMany
	statisticsGetManyTotal
)

func (e *Engine) initQueryAdvertister() {
	e.queryMapAdvertiser[statisticsInsertOne] = models.QueryOptions{
		Query: `
			INSERT INTO statistics
			(
				advert_id,
				advert_name,
				advertiser_name,
				campaign_name,
				views,
				clicks,
				ctr,
				orders,
				cr,
				gmv,
				action_date,
				advert_start_date,
				advert_end_date
			)
			VALUES ($1, $2, $3, $4, $5, $6, $7, $8, $9)
			ON CONFLICT (advert_id, advert_name, action_date) DO UPDATE
				SET clicks            = EXCLUDED.clicks,
					ctr               = EXCLUDED.ctr,
					orders            = EXCLUDED.orders,
					cr                = EXCLUDED.cr,
					gmv               = EXCLUDED.gmv,
					advert_start_date = EXCLUDED.advert_start_date,
					advert_end_date   = EXCLUDED.advert_end_date
		`,
		Timeout: time.Second * 5,
	}

	e.queryMapAdvertiser[statisticsInsertMany] = models.QueryOptions{
		Query: `
			INSERT INTO statistics
			(
				advert_id,
				advert_name,
				advertiser_name,
				campaign_name,
				views,
				clicks,
				ctr,
				orders,
				cr,
				gmv,
				action_date,
				advert_start_date,
				advert_end_date
			)
			VALUES %s
			ON CONFLICT (advert_id, advert_name, action_date) DO UPDATE
				SET views             = EXCLUDED.views,
					clicks            = EXCLUDED.clicks,
					ctr               = EXCLUDED.ctr,
					orders            = EXCLUDED.orders,
					cr                = EXCLUDED.cr,
					gmv               = EXCLUDED.gmv,
					advert_start_date = EXCLUDED.advert_start_date,
					advert_end_date   = EXCLUDED.advert_end_date
		`,
		Timeout: time.Minute * 5,
	}

	e.queryMapAdvertiser[statisticsGetMany] = models.QueryOptions{
		Query: `
		SELECT 
			advert_id,
			SUM(views) AS views,
			SUM(clicks) AS clicks, 
			CASE
				WHEN SUM(views) > 0 THEN ROUND(SUM(clicks) / SUM(views), 4)::FLOAT 
				ELSE 0::FLOAT
			END AS ctr, 
			SUM(orders) AS orders,
			CASE
				WHEN SUM(clicks) > 0 THEN ROUND(SUM(orders) / SUM(clicks), 4)::FLOAT
				ELSE 0::FLOAT
			END AS cr, 
			SUM(gmv) AS gmv
		FROM statistics
		`,
		Timeout: time.Minute * 5,
	}

	e.queryMapAdvertiser[statisticsGetManyTotal] = models.QueryOptions{
		Query: `
		SELECT 
			COUNT(advert_id) OVER() AS total
		FROM statistics
		`,
		Timeout: time.Second * 10,
	}
}
