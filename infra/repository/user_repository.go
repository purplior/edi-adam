package repository

import (
	"github.com/purplior/edi-adam/domain/shared/model"
	domain "github.com/purplior/edi-adam/domain/user"
	"github.com/purplior/edi-adam/infra/database/postgre"
)

type (
	userRepository struct {
		postgreRepository[model.User, domain.QueryOption]
	}
)

func NewUserRepository(
	client *postgre.Client,
) domain.UserRepository {
	var repo postgreRepository[model.User, domain.QueryOption]

	repo.client = client
	repo.applyQueryOption = func(db *postgre.DB, opt domain.QueryOption) (*postgre.DB, error) {
		return db, nil
	}

	return &userRepository{
		postgreRepository: repo,
	}
}
