package persist

import (
	"github.com/podossaem/podoroot/domain/context"
	domain "github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/infra/database"
	"github.com/podossaem/podoroot/infra/database/podosql"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	emailVerificationRepository struct {
		client *podosql.Client
	}
)

func (r *emailVerificationRepository) InsertOne(
	ctx context.APIContext,
	emailVerification domain.EmailVerification,
) (
	domain.EmailVerification,
	error,
) {
	entity := MakeEmailVerification(emailVerification)

	result := r.client.DB.WithContext(ctx).Create(&entity)
	if result.Error != nil {
		return domain.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return entity.ToModel(), nil
}

func (r *emailVerificationRepository) FindOneById(
	ctx context.APIContext,
	id string,
) (
	domain.EmailVerification,
	error,
) {
	eid := dt.UInt(id)

	var entity EmailVerification
	result := r.client.DB.First(&entity, eid)
	if result.Error != nil {
		return domain.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return entity.ToModel(), nil
}

func (r *emailVerificationRepository) FindRecentOneByEmail(
	ctx context.APIContext,
	email string,
) (
	domain.EmailVerification,
	error,
) {
	var entity EmailVerification

	result := r.client.DB.WithContext(ctx).
		Where("email = ?", email).
		Order("created_at DESC").
		First(&entity)

	if result.Error != nil {
		return domain.EmailVerification{}, database.ToDomainError(result.Error)
	}

	return entity.ToModel(), nil
}

func (r *emailVerificationRepository) UpdateOne_IsVerified(
	ctx context.APIContext,
	id string,
	isVerified bool,
) error {
	eid := dt.UInt(id)

	result := r.client.DB.Model(&EmailVerification{}).
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

	result := r.client.DB.Model(&EmailVerification{}).
		Where("id = ?", eid).
		Update("is_consumed", isConsumed)
	if result.Error != nil {
		return database.ToDomainError(result.Error)
	}

	return nil
}

func NewEmailVerificationRepository(
	client *podosql.Client,
) domain.EmailVerificationRepository {
	client.AddModel(&EmailVerification{})

	return &emailVerificationRepository{
		client: client,
	}
}
