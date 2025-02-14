package repository

import (
	"github.com/purplior/sbec/domain/shared/inner"
	"github.com/purplior/sbec/domain/verification"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/database/sqldb"
	"github.com/purplior/sbec/infra/entity"
	"github.com/purplior/sbec/lib/dt"
)

type (
	emailVerificationRepository struct {
		client *sqldb.Client
	}
)

func (r *emailVerificationRepository) InsertOne(
	ctx inner.Context,
	emailVerification verification.EmailVerification,
) (
	verification.EmailVerification,
	error,
) {
	db := r.client.DBWithContext(ctx)
	e := entity.MakeEmailVerification(emailVerification)

	result := db.Create(&e)
	if result.Error != nil {
		return verification.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *emailVerificationRepository) FindOneById(
	ctx inner.Context,
	id string,
) (
	verification.EmailVerification,
	error,
) {
	db := r.client.DBWithContext(ctx)
	eid := dt.UInt(id)

	var e entity.EmailVerification
	result := db.First(&e, eid)
	if result.Error != nil {
		return verification.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *emailVerificationRepository) FindRecentOneByEmail(
	ctx inner.Context,
	email string,
) (
	verification.EmailVerification,
	error,
) {
	var e entity.EmailVerification

	db := r.client.DBWithContext(ctx)
	result := db.
		Where("email = ?", email).
		Order("created_at DESC").
		First(&e)

	if result.Error != nil {
		return verification.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *emailVerificationRepository) UpdateOne_IsVerified(
	ctx inner.Context,
	id string,
	isVerified bool,
) error {
	eid := dt.UInt(id)
	db := r.client.DBWithContext(ctx)

	result := db.
		Model(&entity.EmailVerification{}).
		Where("id = ?", eid).
		Update("is_verified", isVerified)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func (r *emailVerificationRepository) UpdateOne_isConsumed(
	ctx inner.Context,
	id string,
	isConsumed bool,
) error {
	db := r.client.DBWithContext(ctx)
	eid := dt.UInt(id)

	result := db.
		Model(&entity.EmailVerification{}).
		Where("id = ?", eid).
		Update("is_consumed", isConsumed)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func NewEmailVerificationRepository(
	client *sqldb.Client,
) verification.EmailVerificationRepository {
	return &emailVerificationRepository{
		client: client,
	}
}
