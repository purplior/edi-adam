package persist

import (
	"log"

	"github.com/podossaem/podoroot/domain/context"
	"github.com/podossaem/podoroot/domain/exception"
	domain "github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/infra/database/mymongo"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

type (
	userRepository struct {
		mongoClient *mymongo.Client
	}
)

func (r *userRepository) FindByAccount(
	ctx context.APIContext,
	joinMethod string,
	accountID string,
) (
	domain.User,
	error,
) {
	var entity User
	if err := r.
		baseCollection().
		FindOne(
			ctx,
			bson.M{
				"join_method": joinMethod,
				"account_id":  accountID,
			},
		).
		Decode(&entity); err != nil {
		return domain.User{}, err
	}

	return entity.ToModel(), nil
}

func (r *userRepository) InsertOne(
	ctx context.APIContext,
	user domain.User,
) (
	domain.User,
	error,
) {
	entity := MakeUser(user).BeforeInsert()

	result, err := r.baseCollection().InsertOne(ctx, entity)
	if err != nil {
		if mongo.IsDuplicateKeyError(err) {
			log.Println("here..")
			return domain.User{}, exception.ErrAlreadySignedUp
		}
		return domain.User{}, err
	}

	entity.ID = result.InsertedID.(primitive.ObjectID)

	return entity.ToModel(), nil
}

func (r *userRepository) baseCollection() *mymongo.MyMongoCollection {
	return r.mongoClient.MyCollection(mymongo.Collection_User)
}

func NewUserRepository(
	mongoClient *mymongo.Client,
) domain.UserRepository {
	return &userRepository{
		mongoClient: mongoClient,
	}
}
