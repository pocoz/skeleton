package mgo

import (
	"context"

	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/mongo/options"

	"github.com/pocoz/skeleton/models"
)

const collectionIndexScrollIdent = "index_scroll_ident"

func (e *Engine) ScrollIdentUpsert(ctx context.Context, scrollIdent *models.ScrollIdent) error {
	filter := bson.M{
		"index_name_from": scrollIdent.IndexNameFrom,
		"server_from":     scrollIdent.ServerFrom,
		"server_to":       scrollIdent.ServerTo,
	}
	upd := bson.M{"$set": &scrollIdent}
	opts := options.Update().SetUpsert(true)

	_, err := e.db.Collection(collectionIndexScrollIdent).UpdateOne(ctx, filter, upd, opts)
	return err
}

func (e *Engine) ScrollIdentFind(
	ctx context.Context,
	scrollIdentForFind *models.ScrollIdent,
) (*models.ScrollIdent, error) {
	filter := bson.M{
		"index_name_from": scrollIdentForFind.IndexNameFrom,
		"server_from":     scrollIdentForFind.ServerFrom,
		"server_to":       scrollIdentForFind.ServerTo,
	}
	scrollIdent := &models.ScrollIdent{}

	err := e.db.Collection(collectionIndexScrollIdent).FindOne(ctx, filter).Decode(&scrollIdent)
	return scrollIdent, err
}
