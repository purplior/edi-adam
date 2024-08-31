package persist

import (
	"github.com/podossaem/root/domain/context"
	"github.com/podossaem/root/domain/verification"
	"github.com/podossaem/root/infra/database/mymongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	repository struct {
		mongoClient *mymongo.Client
	}
)

func (r *repository) Create(
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
