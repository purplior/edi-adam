package repository

import (
	"github.com/podossaem/podoroot/domain/shared/context"
	"github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/infra/entity"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	emailVerificationRepository struct {
		client *podosql.Client
	}
)

func (r *emailVerificationRepository) InsertOne(
	ctx context.APIContext,
	emailVerification verification.EmailVerification,
) (
	verification.EmailVerification,
	error,
) {
	e := entity.MakeEmailVerification(emailVerification)

	result := r.client.DB.WithContext(ctx).
		Create(&e)
	if result.Error != nil {
		return verification.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *emailVerificationRepository) FindOneById(
	ctx context.APIContext,
	id string,
) (
	verification.EmailVerification,
	error,
) {
	eid := dt.UInt(id)

	var e entity.EmailVerification
	result := r.client.DB.WithContext(ctx).
		First(&e, eid)
	if result.Error != nil {
		return verification.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *emailVerificationRepository) FindRecentOneByEmail(
	ctx context.APIContext,
	email string,
) (
	verification.EmailVerification,
	error,
) {
	var e entity.EmailVerification

	result := r.client.DB.WithContext(ctx).
		Where("email = ?", email).
		Order("created_at DESC").
		First(&e)

	if result.Error != nil {
		return verification.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return e.ToModel(), nil
}

func (r *emailVerificationRepository) UpdateOne_IsVerified(
	ctx context.APIContext,
	id string,
	isVerified bool,
) error {
	eid := dt.UInt(id)

	result := r.client.DB.WithContext(ctx).
		Model(&entity.EmailVerification{}).
		Where("id = ?", eid).
		Update("is_verified", isVerified)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func (r *emailVerificationRepository) UpdateOne_isConsumed(
	ctx context.APIContext,
	id string,
	isConsumed bool,
) error {
	eid := dt.UInt(id)

	result := r.client.DB.WithContext(ctx).
		Model(&entity.EmailVerification{}).
		Where("id = ?", eid).
		Update("is_consumed", isConsumed)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func NewEmailVerificationRepository(
	client *podosql.Client,
) verification.EmailVerificationRepository {
	return &emailVerificationRepository{
		client: client,
	}
}
