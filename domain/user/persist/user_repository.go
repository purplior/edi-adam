package persist

import (
	"github.com/podossaem/podoroot/domain/context"
	domain "github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/infra/database/mymongo"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	userRepository struct {
		mongoClient *mymongo.Client
	}
)

func (r *userRepository) InsertOne(
	ctx context.APIContext,
	user domain.User,
) (
	domain.User,
	error,
) {
	entity := MakeUser(user)

	result, err := r.baseCollection().InsertOne(ctx, entity)
	if err != nil {
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
