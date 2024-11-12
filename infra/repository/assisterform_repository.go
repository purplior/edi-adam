package repository

import (
	"fmt"

	domain "github.com/podossaem/podoroot/domain/assisterform"
	"github.com/podossaem/podoroot/domain/shared/inner"
	"github.com/podossaem/podoroot/infra/database"
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
		fmt.Println("여기 아니야?")
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

func NewAssisterFormRepository(
	client *podomongo.Client,
) domain.AssisterFormRepository {
	return &assisterFormRepository{
		client: client,
	}
}
