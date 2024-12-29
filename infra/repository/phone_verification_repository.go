package repository

import (
	"github.com/purplior/podoroot/domain/shared/inner"
	domain "github.com/purplior/podoroot/domain/verification"
	"github.com/purplior/podoroot/infra/database"
	"github.com/purplior/podoroot/infra/database/podosql"
	"github.com/purplior/podoroot/infra/entity"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	phoneVerificationRepository struct {
		client *podosql.Client
	}
)

func (r *phoneVerificationRepository) InsertOne(
	ctx inner.Context,
	phoneVerification domain.PhoneVerification,
) (
	domain.PhoneVerification,
	error,
) {
	db := r.client.DBWithContext(ctx)
	e := entity.MakePhoneVerification(phoneVerification)

	result := db.Create(&e)
	if result.Error != nil {
		return domain.PhoneVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *phoneVerificationRepository) FindOneById(
	ctx inner.Context,
	id string,
) (
	domain.PhoneVerification,
	error,
) {
	db := r.client.DBWithContext(ctx)
	eid := dt.UInt(id)

	var e entity.PhoneVerification
	result := db.First(&e, eid)
	if result.Error != nil {
		return domain.PhoneVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *phoneVerificationRepository) FindRecentOneByPhoneNumber(
	ctx inner.Context,
	phoneNumber string,
) (
	domain.PhoneVerification,
	error,
) {
	var e entity.PhoneVerification

	db := r.client.DBWithContext(ctx)
	result := db.
		Where("phone_number = ?", phoneNumber).
		Order("created_at DESC").
		First(&e)

	if result.Error != nil {
		return domain.PhoneVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *phoneVerificationRepository) UpdateOne_IsVerified(
	ctx inner.Context,
	id string,
	isVerified bool,
) error {
	eid := dt.UInt(id)
	db := r.client.DBWithContext(ctx)

	result := db.
		Model(&entity.PhoneVerification{}).
		Where("id = ?", eid).
		Update("is_verified", isVerified)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func (r *phoneVerificationRepository) UpdateOne_isConsumed(
	ctx inner.Context,
	id string,
	isConsumed bool,
) error {
	db := r.client.DBWithContext(ctx)
	eid := dt.UInt(id)

	result := db.
		Model(&entity.PhoneVerification{}).
		Where("id = ?", eid).
		Update("is_consumed", isConsumed)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func NewPhoneVerificationRepository(
	client *podosql.Client,
) domain.PhoneVerificationRepository {
	return &phoneVerificationRepository{
		client: client,
	}
}
