package repository

import (
	domain "github.com/purplior/podoroot/domain/assisterform"
	"github.com/purplior/podoroot/domain/shared/inner"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podomongo"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/lib/dt"
	"github.com/purplior/podoroot/lib/mydate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	assisterFormRepository struct {
		client *podomongo.Client
	}
)

func (r *assisterFormRepository) InsertOne(
	ctx inner.Context,
	assisterForm domain.AssisterForm,
) (
	domain.AssisterForm,
	error,
) {
	e := entity.MakeAssisterForm(assisterForm)
	e.CreatedAt = mydate.Now()

	result, err := r.client.
		MyCollection(podomongo.Collection_AssisterForm).
		InsertOne(ctx.Value(), e)
	if err != nil {
		return domain.AssisterForm{}, database.ToDomainError(err)
	}

	e.ID = result.InsertedID.(primitive.ObjectID)

	return e.ToModel(), nil
}

func (r *assisterFormRepository) FindOne_ByID(
	ctx inner.Context,
	id string,
) (
	domain.AssisterForm,
	error,
) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var e entity.AssisterForm
	if err := r.client.
		MyCollection(podomongo.Collection_AssisterForm).
		FindOne(
			ctx.Value(),
			bson.M{
				"_id": oid,
			},
		).Decode(&e); err != nil {
		return domain.AssisterForm{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *assisterFormRepository) FindOne_ByAssisterID(
	ctx inner.Context,
	assisterID string,
) (
	domain.AssisterForm,
	error,
) {
	assisterEntityID := dt.UInt(assisterID)

	var e entity.AssisterForm
	if err := r.client.
		MyCollection(podomongo.Collection_AssisterForm).
		FindOne(
			ctx.Value(),
			bson.M{
				"assisterId": assisterEntityID,
			},
		).Decode(&e); err != nil {
		return domain.AssisterForm{}, database.ToDomainError(err)
	}

	return e.ToModel(), nil
}

func (r *assisterFormRepository) UpdateOne(
	ctx inner.Context,
	assisterForm domain.AssisterForm,
) error {
	oid, _ := primitive.ObjectIDFromHex(assisterForm.ID)
	e := entity.MakeAssisterForm(assisterForm)

	_, err := r.client.
		MyCollection(podomongo.Collection_AssisterForm).
		UpdateOne(
			ctx.Value(),
			bson.M{
				"_id": oid,
			},
			bson.M{
				"$set": e,
			},
		)

	return database.ToDomainError(err)
}

func (r *assisterFormRepository) DeleteAll_ByAssisterIDs(
	ctx inner.Context,
	assisterIDs []string,
) error {
	eIDs := make([]int, len(assisterIDs))
	for i, assisterID := range assisterIDs {
		eIDs[i] = dt.Int(assisterID)
	}

	_, err := r.client.
		MyCollection(podomongo.Collection_AssisterForm).
		DeleteMany(
			ctx.Value(),
			bson.M{
				"assisterId": bson.M{
					"$in": eIDs,
				},
			},
		)

	if err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func NewAssisterFormRepository(
	client *podomongo.Client,
) domain.AssisterFormRepository {
	return &assisterFormRepository{
		client: client,
	}
}
