package repository

import (
	"github.com/purplior/sbec/domain/shared/exception"
	"github.com/purplior/sbec/domain/shared/inner"
	domain "github.com/purplior/sbec/domain/verification"
	"github.com/purplior/sbec/infra/database"
	"github.com/purplior/sbec/infra/database/sqldb"
	"github.com/purplior/sbec/infra/entity"
	"github.com/purplior/sbec/lib/dt"
	"github.com/purplior/sbec/lib/mydate"
)

type (
	phoneVerificationRepository struct {
		client *sqldb.Client
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

func (r *phoneVerificationRepository) FindRecentOne_ByPhoneNumber(
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

func (r *phoneVerificationRepository) FindCount_ByPhoneNumber(
	ctx inner.Context,
	phoneNumber string,
) (
	int,
	error,
) {
	startAt := mydate.DayStartFromNow(0)
	endAt := mydate.DayEndFromNow(0)

	var count int64
	db := r.client.DBWithContext(ctx)
	err := db.Model(&entity.PhoneVerification{}).
		Where("phone_number = ? AND created_at >= ? AND created_at < ?", phoneNumber, startAt, endAt).
		Count(&count).
		Error

	if err != nil {
		domainErr := database.ToDomainError(err)
		if domainErr == exception.ErrNoRecord {
			return 0, nil
		}

		return 0, domainErr
	}

	return int(count), nil
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
	client *sqldb.Client,
) domain.PhoneVerificationRepository {
	return &phoneVerificationRepository{
		client: client,
	}
}
