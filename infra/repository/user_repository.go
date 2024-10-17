package repository

import (
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
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
	user.User,
	error,
) {
	var e entity.User
	result := r.client.DB.WithContext(ctx).
		Where("join_method = ?", joinMethod).
		Where("account_id = ?", accountID).
		First(&e)
	if result.Error != nil {
		return user.User{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *userRepository) InsertOne(
	ctx context.APIContext,
	userForInsert user.User,
) (
	user.User,
	error,
) {
	e := entity.MakeUser(userForInsert)
	result := r.client.DB.WithContext(ctx).
		Create(&e)

	if result.Error != nil {
		return user.User{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func NewUserRepository(
	client *podosql.Client,
) user.UserRepository {
	return &userRepository{
		client: client,
	}
}
