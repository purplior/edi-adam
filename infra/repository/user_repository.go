package repository

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	domain "github.com/purplior/podoroot/domain/user"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
)

type (
	userRepository struct {
		client *podosql.Client
	}
)

func (r *userRepository) FindOne_ByID(
	ctx inner.Context,
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

func (r *userRepository) FindOne_ByAccount(
	ctx inner.Context,
	joinMethod string,
	accountID string,
) (
	domain.User,
	error,
) {
	var e entity.User

	db := r.client.DBWithContext(ctx)
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
	ctx inner.Context,
	userForInsert domain.User,
) (
	domain.User,
	error,
) {
	e := entity.MakeUser(userForInsert)

	db := r.client.DBWithContext(ctx)
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
