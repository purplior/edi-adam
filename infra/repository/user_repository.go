package repository

import (
	"time"

	"github.com/purplior/sbec/domain/shared/inner"
	domain "github.com/purplior/sbec/domain/user"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/database/sqldb"
	"github.com/purplior/sbec/infra/entity"
)

type (
	userRepository struct {
		client *sqldb.Client
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

func (r *userRepository) FindOne_ByNickname(
	ctx inner.Context,
	nickname string,
) (
	domain.User,
	error,
) {
	var e entity.User

	db := r.client.DBWithContext(ctx)
	err := db.
		Where("nickname = ?", nickname).
		First(&e).Error
	if err != nil {
		return domain.User{}, database.ToDomainError(err)
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

func (r *userRepository) UpdateOne_InactivatedFields(
	ctx inner.Context,
	userID string,
	isInactivated bool,
	inactivatedAt time.Time,
) error {
	db := r.client.DBWithContext(ctx)
	err := db.Model(&entity.User{}).
		Where("id = ?", userID).
		Updates(map[string]interface{}{
			"is_inactivated": isInactivated,
			"inactivated_at": inactivatedAt,
		}).
		Error

	return err
}

func (r *userRepository) UpdateOne_Password_ByAccount(
	ctx inner.Context,
	joinMethod string,
	accountID string,
	newPassword string,
) error {
	db := r.client.DBWithContext(ctx)
	err := db.Model(&entity.User{}).
		Where("join_method = ?", joinMethod).
		Where("account_id = ?", accountID).
		Updates(map[string]interface{}{
			"account_password": newPassword,
		}).
		Error

	if err != nil {
		return database.ToDomainError(err)
	}

	return nil
}

func NewUserRepository(
	client *sqldb.Client,
) domain.UserRepository {
	return &userRepository{
		client: client,
	}
}
