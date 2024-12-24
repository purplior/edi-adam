package category

import "github.com/purplior/podoroot/domain/shared/inner"

type (
	CategoryService interface {
		GetMainCategoryList(
			ctx inner.Context,
		) (
			[]Category,
			error,
		)
	}
)

type (
	categoryService struct {
		categoryRepository CategoryRepository
	}
)

func (s *categoryService) GetMainCategoryList(
	ctx inner.Context,
) (
	[]Category,
	error,
) {
	return s.categoryRepository.FindListByIDs(ctx, []string{
		"1",
	})
}

func NewCategoryService(
	categoryRepository CategoryRepository,
) CategoryService {
	return &categoryService{
		categoryRepository: categoryRepository,
	}
}
