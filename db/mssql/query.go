package mssql

import (
	"context"
	"time"
)

type options struct {
	query   string
	timeout time.Duration
}

func (e *Engine) getQuery(ctx context.Context, name uint8) (context.Context, string) {
	opts := e.queryMap[name]
	ctx, _ = context.WithTimeout(ctx, opts.timeout)
	return ctx, opts.query
}

const (
	getOffersByItems = uint8(iota)
)

func (e *Engine) queryInit() {

	e.queryMap[getOffersByItems] = options{
		query: `
			SELECT i.item_id,
				   i.short_web_name,
				   t.level,
				   t.level1_category_name,
				   t.level2_category_name,
				   t.level3_category_name,
				   t.level4_category_name,
				   t.level5_category_name,
				   t.level6_category_name,
				   t.level7_category_name,
				   t.level8_category_name,
				   t.level9_category_name,
				   t.level10_category_name
			FROM OpenJson(@values) WITH (item_id NVARCHAR(250)) AS vi
					 JOIN Malibu.items AS i
						 ON i.item_id = vi.item_id
								AND i.display_on_web = 1
					 LEFT JOIN Malibu.prices p
							   ON p.partition_id = i.partition_id
								   AND p.item_id = i.item_id
					 JOIN Malibu.item_collections icw
						  ON icw.item_id = i.item_id
							  AND icw.collection_type_id = 2
							  AND icw.is_main = 1
					 LEFT JOIN Malibu.item_collections icwgod
							   ON icwgod.item_id = i.item_id
								   AND icwgod.collection_id = 682854
					 JOIN Malibu.collections c
						  ON c.collection_id = icw.collection_id
							  AND c.is_active = 1
							  AND CHARINDEX('Link', c.collection_name) = 0
							  AND CHARINDEX('Link', c.identifier) = 0
					 JOIN Malibu.tree t
						  ON t.category_id = c.identifier
							  AND t.structure_id = 10001
			WHERE (isNull(p.is_available, 0) = 0 OR p.is_available = 1)
		`,
		timeout: time.Second * 360,
	}
}
