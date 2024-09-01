package persist

import (
	"github.com/podossaem/podoroot/domain/context"
	domain "github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/infra/database/mymongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	emailVerificationRepository struct {
		mongoClient *mymongo.Client
	}
)

func (r *emailVerificationRepository) InsertOne(
	ctx context.APIContext,
	emailVerification domain.EmailVerification,
) (
	createdEmailVerification domain.EmailVerification,
	err error,
) {
	entity := MakeEmailVerification(emailVerification).BeforeInsert()

	if result, err := r.baseCollection().InsertOne(ctx, entity); err != nil {
		return createdEmailVerification, err
	} else {
		entity.ID = result.InsertedID.(primitive.ObjectID)
		createdEmailVerification = entity.ToModel()
	}

	return createdEmailVerification, nil
}

func (r *emailVerificationRepository) FindOneById(
	ctx context.APIContext,
	id string,
) (
	domain.EmailVerification,
	error,
) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var entity EmailVerification
	err := r.
		baseCollection().
		FindOne(
			ctx,
			bson.M{
				"_id": oid,
			},
			options.
				FindOne().
				SetSort(bson.M{"_id": -1}),
		).
		Decode(&entity)

	if err != nil {
		return domain.EmailVerification{}, err
	}

	return entity.ToModel(), nil
}

func (r *emailVerificationRepository) FindOneByEmail(
	ctx context.APIContext,
	email string,
) (
	domain.EmailVerification,
	error,
) {
	var entity EmailVerification
	err := r.
		baseCollection().
		FindOne(
			ctx,
			bson.M{
				"email": email,
			},
			options.
				FindOne().
				SetSort(bson.M{"_id": -1}),
		).
		Decode(&entity)

	if err != nil {
		return domain.EmailVerification{}, err
	}

	return entity.ToModel(), nil
}

func (r *emailVerificationRepository) UpdateOne_IsVerified(
	ctx context.APIContext,
	id string,
	isVerified bool,
) (
	domain.EmailVerification,
	error,
) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var entity EmailVerification
	err := r.
		baseCollection().
		FindOneAndUpdate(
			ctx,
			bson.M{
				"_id": oid,
			},
			bson.M{
				"$set": bson.M{
					"is_verified": isVerified,
				},
			},
			options.
				FindOneAndUpdate().
				SetReturnDocument(options.After),
		).
		Decode(&entity)

	if err != nil {
		return domain.EmailVerification{}, err
	}

	return entity.ToModel(), err
}

func (r *emailVerificationRepository) UpdateOne_isConsumed(
	ctx context.APIContext,
	id string,
	isConsumed bool,
) (
	domain.EmailVerification,
	error,
) {
	oid, _ := primitive.ObjectIDFromHex(id)

	var entity EmailVerification
	err := r.
		baseCollection().
		FindOneAndUpdate(
			ctx,
			bson.M{
				"_id": oid,
			},
			bson.M{
				"$set": bson.M{
					"is_consumed": isConsumed,
				},
			},
			options.
				FindOneAndUpdate().
				SetReturnDocument(options.After),
		).
		Decode(&entity)

	if err != nil {
		return domain.EmailVerification{}, err
	}

	return entity.ToModel(), err
}

func (r *emailVerificationRepository) baseCollection() *mymongo.MyMongoCollection {
	return r.mongoClient.MyCollection(mymongo.Collection_EmailVerification)
}

func NewEmailVerificationRepository(
	mongoClient *mymongo.Client,
) domain.EmailVerificationRepository {
	return &emailVerificationRepository{
		mongoClient: mongoClient,
	}
}
