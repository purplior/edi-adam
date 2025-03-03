package repository

import (
	"github.com/purplior/edi-adam/domain/shared/model"
	domain "github.com/purplior/edi-adam/domain/verification"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	verificationRepository struct {
		postgreRepository[model.Verification, domain.QueryOption]
	}
)

// func (r *verificationRepository) applyQueryOption(query *gorm.DB, opt domain.QueryOption) *gorm.DB {
// 	if opt.ID > 0 {
// 		query = query.Where("id = ?", opt.ID)
// 	}
// 	if len(opt.Hash) > 0 {
// 		query = query.Where("hash = ?", opt.Hash)
// 	}
// 	if !opt.CreatedAtStart.IsZero() {
// 		query = query.Where("created_at >= ?", opt.CreatedAtStart)
// 	}
// 	if !opt.CreatedAtEnd.IsZero() {
// 		query = query.Where("created_at < ?", opt.CreatedAtEnd)
// 	}

// 	return query
// }

func NewVerificationRepository(
	client *postgre.Client,
) domain.VerificationRepository {
	var repo postgreRepository[model.Verification, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		query := db
		if opt.ID > 0 {
			query = query.Where("id = ?", opt.ID)
		}
		if len(opt.Hash) > 0 {
			query = query.Where("hash = ?", opt.Hash)
		}
		if !opt.CreatedAtStart.IsZero() {
			query = query.Where("created_at >= ?", opt.CreatedAtStart)
		}
		if !opt.CreatedAtEnd.IsZero() {
			query = query.Where("created_at < ?", opt.CreatedAtEnd)
		}
		return query, nil
	}

	return &verificationRepository{
		postgreRepository: repo,
	}
}
