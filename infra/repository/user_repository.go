package repository

import (
	"github.com/podossaem/podoroot/domain/shared/context"
	domain "github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
)

type (
	userRepository struct {
		client *podosql.Client
	}
)

func (r *userRepository) FindOneByID(
	ctx context.APIContext,
	id string,
) (
	domain.User,
	error,
) {
	var e entity.User
	db := r.client.DBWithContext(ctx)

	result := db.
		Where("id = ?", id).
		First(&e)
	if result.Error != nil {
		return domain.User{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *userRepository) FindOneByAccount(
	ctx context.APIContext,
	joinMethod string,
	accountID string,
) (
	domain.User,
	error,
) {
	var e entity.User

	db := r.client.DB.WithContext(ctx)
	result := db.
		Where("join_method = ?", joinMethod).
		Where("account_id = ?", accountID).
		First(&e)
	if result.Error != nil {
		return domain.User{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *userRepository) InsertOne(
	ctx context.APIContext,
	userForInsert domain.User,
) (
	domain.User,
	error,
) {
	e := entity.MakeUser(userForInsert)

	db := r.client.DB.WithContext(ctx)
	result := db.Create(&e)

	if result.Error != nil {
		return domain.User{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func NewUserRepository(
	client *podosql.Client,
) domain.UserRepository {
	return &userRepository{
		client: client,
	}
}
