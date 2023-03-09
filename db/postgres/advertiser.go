package postgres

import (
	"context"
	"fmt"
	"strings"

	"github.com/lib/pq"

	"github.com/pocoz/skeleton/models"
	"github.com/pocoz/skeleton/tools"
)

func (e *Engine) StatisticsInsertOne(statistics *models.StatisticsPostgres) error {
	opts := e.queryMapAdvertiser[statisticsInsertOne]
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	_, err := e.db.ExecContext(
		ctx,
		opts.Query,
		statistics.AdvertID,
		statistics.AdvertName,
		statistics.AdvertiserName,
		statistics.CampaignName,
		statistics.Views,
		statistics.Clicks,
		statistics.CTR,
		statistics.Orders,
		statistics.CR,
		statistics.GMV,
		statistics.ActionDate,
		statistics.AdvertStartDate,
		statistics.AdvertEndDate,
	)
	if err != nil {
		return err
	}

	return nil
}

func (e *Engine) StatisticsInsertMany(statisticsList []*models.StatisticsPostgres) error {
	var (
		amount = tools.GetAmount(len(statisticsList), e.options.maxButchSize)
		start  int
		end    int
	)

	for i := 0; i < amount; i++ {
		err := func() error {
			start = e.options.maxButchSize * i
			end = e.options.maxButchSize * (i + 1)
			if end > len(statisticsList) {
				end = len(statisticsList)
			}

			part := statisticsList[start:end]
			var values string
			for _, statistics := range part {
				statistics.AdvertName = strings.Replace(statistics.AdvertName, "'", "", -1)
				statistics.AdvertiserName = strings.Replace(statistics.AdvertiserName, "'", "", -1)
				statistics.CampaignName = strings.Replace(statistics.CampaignName, "'", "", -1)

				values = fmt.Sprintf("%s ('%s', '%s', '%s', '%s', %d, %d, %f, %d, %f, %d, '%s', '%s', '%s'),",
					values,
					statistics.AdvertID,
					statistics.AdvertName,
					statistics.AdvertiserName,
					statistics.CampaignName,
					statistics.Views,
					statistics.Clicks,
					statistics.CTR,
					statistics.Orders,
					statistics.CR,
					statistics.GMV,
					statistics.ActionDate.String(),
					statistics.AdvertStartDate.String(),
					statistics.AdvertEndDate.String(),
				)
			}
			values = strings.TrimSuffix(values, ",")

			opts := e.queryMapAdvertiser[statisticsInsertMany]
			opts.Query = fmt.Sprintf(opts.Query, values)
			ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
			defer cancel()

			_, err := e.db.ExecContext(
				ctx,
				opts.Query,
			)
			if err != nil {
				fmt.Println(opts.Query)
				return err
			}

			return nil
		}()
		if err != nil {
			return err
		}
	}

	return nil
}

func (e *Engine) GetAdvertsStatistic(
	ids []string,
	startDate string,
	endDate string,
	orderBy string,
	sortingDirection string,
	limit int,
	offset int,
) (mds []*models.StatisticsPostgres, err error) {
	opts := e.queryMapAdvertiser[statisticsGetMany]
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	q, attrs := prepareGetAdvertsStatisticQuery(opts.Query, ids, startDate, endDate, orderBy, sortingDirection, limit, offset, true)

	rows, err := e.db.QueryContext(ctx, q, attrs...)
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		model := new(models.StatisticsPostgres)
		err := rows.Scan(
			&model.AdvertID,
			&model.Views,
			&model.Clicks,
			&model.CTR,
			&model.Orders,
			&model.CR,
			&model.GMV,
		)
		if err != nil {
			return nil, err
		}

		mds = append(mds, model)
	}

	return
}

func (e *Engine) GetAdvertsStatisticTotal(
	ids []string,
	startDate string,
	endDate string,
	orderBy string,
	sortingDirection string,
	limit int,
	offset int,
) (int64, error) {
	opts := e.queryMapAdvertiser[statisticsGetManyTotal]
	ctx, cancel := context.WithTimeout(context.Background(), opts.Timeout)
	defer cancel()

	q, attrs := prepareGetAdvertsStatisticQuery(opts.Query, ids, startDate, endDate, orderBy, sortingDirection, limit, offset, false)

	var total int64
	err := e.db.QueryRowContext(ctx, q, attrs...).Scan(&total)
	if err != nil {
		return 0, err
	}

	return total, nil
}

func prepareGetAdvertsStatisticQuery(
	sq string,
	ids []string,
	startDate string,
	endDate string,
	orderBy string,
	sortingDirection string,
	limit int,
	offset int,
	withOrderBy bool,
) (string, []interface{}) {
	q, attrs, i := sq, make([]interface{}, 0, 6), 0
	if len(ids) > 0 {
		q, attrs, i = addWhereClause(q, attrs, ` %s advert_id = ANY($%d)`, pq.Array(ids), i)
	}

	if len(startDate) > 0 {
		q, attrs, i = addWhereClause(q, attrs, ` %s advert_start_date >= $%d::date`, startDate, i)
	}

	if len(endDate) > 0 {
		q, attrs, i = addWhereClause(q, attrs, ` %s advert_end_date <= $%d::date`, endDate, i)
	}

	q += ` GROUP BY advert_id`

	if withOrderBy {
		q += fmt.Sprintf(` ORDER BY %s %s`, orderBy, sortingDirection)
	}

	if limit > 0 {
		i++
		q += fmt.Sprintf(` LIMIT $%d`, i)
		attrs = append(attrs, limit)
	}

	i++
	q += fmt.Sprintf(` OFFSET $%d`, i)
	attrs = append(attrs, offset)

	return q, attrs
}

func addWhereClause(
	q string,
	attrs []interface{},
	clause string,
	attr interface{},
	it int,
) (string, []interface{}, int) {
	op := "WHERE"
	if it > 0 {
		op = "AND"
	}

	it++
	q += fmt.Sprintf(clause, op, it)
	attrs = append(attrs, attr)

	return q, attrs, it
}
