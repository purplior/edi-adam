package repository

import (
	domain "github.com/purplior/podoroot/domain/assister"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podomongo"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/lib/mydate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	assisterRepository struct {
		client *podomongo.Client
	}
)

func (r *assisterRepository) FindOne_ByID(
	ctx inner.Context,
	id string,
) (
	domain.Assister,
	error,
) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var e entity.Assister
	if err := r.client.
		MyCollection(podomongo.Collection_Assister).
		FindOne(
			ctx.Value(),
			bson.M{
				"_id": oid,
			},
		).Decode(&e); err != nil {
		return domain.Assister{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *assisterRepository) InsertOne(
	ctx inner.Context,
	assister domain.Assister,
) (
	domain.Assister,
	error,
) {
	e := entity.MakeAssister(assister)
	e.CreatedAt = mydate.Now()

	result, err := r.client.
		MyCollection(podomongo.Collection_Assister).
		InsertOne(ctx.Value(), e)
	if err != nil {
		return domain.Assister{}, database.ToDomainError(err)
	}

	e.ID = result.InsertedID.(primitive.ObjectID)

	return e.ToModel(), nil
}

func (r *assisterRepository) UpdateOne(
	ctx inner.Context,
	assister domain.Assister,
) error {
	e := entity.MakeAssister(assister)

	_, err := r.client.
		MyCollection(podomongo.Collection_Assister).
		UpdateOne(
			ctx.Value(),
			bson.M{
				"_id": e.ID,
			},
			bson.M{
				"$set": e,
			},
		)
	if err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func (r *assisterRepository) DeleteOne_ByID(
	ctx inner.Context,
	id string,
) error {
	eID, err := primitive.ObjectIDFromHex(id)
	if err != nil {
		return database.ToDomainError(err)
	}

	_, err = r.client.
		MyCollection(podomongo.Collection_Assister).
		DeleteOne(
			ctx.Value(),
			bson.M{
				"_id": eID,
			},
		)
	if err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func NewAssisterRepository(
	client *podomongo.Client,
) domain.AssisterRepository {
	return &assisterRepository{
		client: client,
	}
}
