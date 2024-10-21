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
	return r.assisterFormRepository.InsertOne(ctx, AssisterForm{
		Fields: request.Fields,
	})
}

func (r *assisterFormService) GetOneByID(
	ctx context.APIContext,
	id string,
) (
	AssisterForm,
	error,
) {
	return r.assisterFormRepository.FindOneByID(ctx, id)
}

func NewAssisterFormService(
	assisterFormRepository AssisterFormRepository,
) AssisterFormService {
	return &assisterFormService{
		assisterFormRepository: assisterFormRepository,
	}
}
