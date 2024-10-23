package repository

import (
	domain "github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/infra/database/podomongo"
	"github.com/podossaem/podoroot/infra/entity"
	"github.com/podossaem/podoroot/lib/dt"
	"github.com/podossaem/podoroot/lib/mydate"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	assisterFormRepository struct {
		client *podomongo.Client
	}
)

func (r *assisterFormRepository) InsertOne(
	ctx context.APIContext,
	assisterForm domain.AssisterForm,
) (
	domain.AssisterForm,
	error,
) {
	e := entity.MakeAssisterForm(assisterForm)
	e.CreatedAt = mydate.Now()

	result, err := r.client.MyCollection(podomongo.Collection_AssisterForm).InsertOne(ctx, e)
	if err != nil {
		return domain.AssisterForm{}, err
	}

	e.ID = result.InsertedID.(primitive.ObjectID)

	return e.ToModel(), nil
}

func (r *assisterFormRepository) FindOneByID(
	ctx context.APIContext,
	id string,
) (
	domain.AssisterForm,
	error,
) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var e entity.AssisterForm
	if err := r.client.MyCollection(podomongo.Collection_AssisterForm).FindOne(
		ctx,
		bson.M{
			"_id": oid,
		},
	).Decode(&e); err != nil {
		return domain.AssisterForm{}, err
	}

	return e.ToModel(), nil
}

func (r *assisterFormRepository) FindOneByAssisterID(
	ctx context.APIContext,
	assisterID string,
) (
	domain.AssisterForm,
	error,
) {
	assisterEntityID := dt.UInt(assisterID)

	var e entity.AssisterForm
	if err := r.client.MyCollection(podomongo.Collection_AssisterForm).FindOne(
		ctx,
		bson.M{
			"assisterId": assisterEntityID,
		},
	).Decode(&e); err != nil {
		return domain.AssisterForm{}, err
	}

	return e.ToModel(), nil
}

func NewAssisterFormRepository(
	client *podomongo.Client,
) domain.AssisterFormRepository {
	return &assisterFormRepository{
		client: client,
	}
}
