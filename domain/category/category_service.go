package category

import (
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
)

type (
	CategoryService interface {
		Get(
			session inner.Session,
			queryOption QueryOption,
		) (
			model.Category,
			error,
		)

		GetMainList(
			session inner.Session,
		) (
			[]model.Category,
			error,
		)
	}
)

type (
	categoryService struct {
		categoryRepository CategoryRepository
	}
)

func (s *categoryService) Get(
	session inner.Session,
	queryOption QueryOption,
) (
	model.Category,
	error,
) {
	return s.categoryRepository.Read(
		session,
		queryOption,
	)
}

func (s *categoryService) GetMainList(
	session inner.Session,
) (
	[]model.Category,
	error,
) {
	return s.categoryRepository.ReadList(
		session,
		QueryOption{},
	)
}

func NewCategoryService(
	categoryRepository CategoryRepository,
) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
