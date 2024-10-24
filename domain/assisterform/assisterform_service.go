package assisterform

import "github.com/podossaem/podoroot/domain/shared/context"

type (
	AssisterFormService interface {
		RegisterOne(
			ctx context.APIContext,
			request AssisterFormRegisterRequest,
		) (
			AssisterForm,
			error,
		)

		GetOneByID(
			ctx context.APIContext,
			id string,
		) (
			AssisterForm,
			error,
		)

		GetOneByAssisterID(
			ctx context.APIContext,
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
	ctx context.APIContext,
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
	ctx context.APIContext,
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
	ctx context.APIContext,
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
