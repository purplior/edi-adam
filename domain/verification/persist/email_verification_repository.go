package persist

import (
	"github.com/podossaem/root/domain/context"
	"github.com/podossaem/root/domain/verification"
	"github.com/podossaem/root/infra/database/mymongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

type (
	repository struct {
		mongoClient *mymongo.Client
	}
)

func (r *repository) InsertOne(
	ctx context.APIContext,
	emailVerification verification.EmailVerification,
) (
	createdEmailVerification verification.EmailVerification,
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

func (r *repository) FindOneByEmail(
	ctx context.APIContext,
	email string,
) (
	verification.EmailVerification,
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
		return verification.EmailVerification{}, err
	}

	return entity.ToModel(), nil
}

func (r *repository) UpdateOne_IsVerified(
	ctx context.APIContext,
	id string,
	isVerified bool,
) (
	verification.EmailVerification,
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
		return verification.EmailVerification{}, err
	}

	return entity.ToModel(), err
}

func (r *repository) baseCollection() *mymongo.MyMongoCollection {
	return r.mongoClient.MyCollection(mymongo.Collection_EmailVerification)
}

func NewEmailVerificationRepository(
	mongoClient *mymongo.Client,
) verification.EmailVerificationRepository {
	return &repository{
		mongoClient: mongoClient,
	}
}
