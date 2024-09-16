package persist

import (
	"github.com/podossaem/podoroot/domain/context"
	domain "github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
)

type (
	userRepository struct {
		client *podosql.Client
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
	result := r.client.DB.WithContext(ctx).
		Where("join_method = ?", joinMethod).
		Where("account_id = ?", accountID).
		First(&entity)
	if result.Error != nil {
		return domain.User{}, database.ToDomainError(result.Error)
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
	entity := MakeUser(user)
	result := r.client.DB.WithContext(ctx).
		Create(&entity)

	if result.Error != nil {
		return domain.User{}, database.ToDomainError(result.Error)
	}

	model := entity.ToModel()

	return model, nil
}

func NewUserRepository(
	client *podosql.Client,
) domain.UserRepository {
	client.AddModel(&User{})

	return &userRepository{
		client: client,
	}
}
