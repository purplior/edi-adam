package assisterform

import "github.com/podossaem/podoroot/domain/shared/inner"

type (
	AssisterFormService interface {
		RegisterOne(
			ctx inner.Context,
			request AssisterFormRegisterRequest,
		) (
			AssisterForm,
			error,
		)

		GetOneByID(
			ctx inner.Context,
			id string,
		) (
			AssisterForm,
			error,
		)

		GetOneByAssisterID(
			ctx inner.Context,
			assistantID string,
		) (
			AssisterForm,
			error,
		)
	}
)

type (
	assisterFormService struct {
		assisterFormRepository AssisterFormRepository
	}
)

func (r *assisterFormService) RegisterOne(
	ctx inner.Context,
	request AssisterFormRegisterRequest,
) (
	AssisterForm,
	error,
) {
	return r.assisterFormRepository.InsertOne(
		ctx,
		AssisterForm{
			Fields:        request.Fields,
			QueryMessages: request.QueryMessages,
		},
	)
}

func (r *assisterFormService) GetOneByID(
	ctx inner.Context,
	id string,
) (
	AssisterForm,
	error,
) {
	return r.assisterFormRepository.FindOneByID(
		ctx,
		id,
	)
}

func (r *assisterFormService) GetOneByAssisterID(
	ctx inner.Context,
	assisterID string,
) (
	AssisterForm,
	error,
) {
	return r.assisterFormRepository.FindOneByAssisterID(
		ctx,
		assisterID,
	)
}

func NewAssisterFormService(
	assisterFormRepository AssisterFormRepository,
) AssisterFormService {
	return &assisterFormService{
		assisterFormRepository: assisterFormRepository,
	}
}
