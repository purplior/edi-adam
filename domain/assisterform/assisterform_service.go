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

		GetOne_ByID(
			ctx inner.Context,
			id string,
		) (
			AssisterForm,
			error,
		)

		GetOne_ByAssisterID(
			ctx inner.Context,
			id string,
		) (
			AssisterForm,
			error,
		)

		GetViewOne_ByAssister(
			ctx inner.Context,
			assistantID string,
		) (
			AssisterFormView,
			error,
		)

		PutOne(
			ctx inner.Context,
			assisterForm AssisterForm,
		) error
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

func (r *assisterFormService) GetOne_ByID(
	ctx inner.Context,
	id string,
) (
	AssisterForm,
	error,
) {
	return r.assisterFormRepository.FindOne_ByID(
		ctx,
		id,
	)
}

func (r *assisterFormService) GetOne_ByAssisterID(
	ctx inner.Context,
	assisterID string,
) (
	AssisterForm,
	error,
) {
	return r.assisterFormRepository.FindOne_ByAssisterID(
		ctx,
		assisterID,
	)
}

func (r *assisterFormService) GetViewOne_ByAssister(
	ctx inner.Context,
	assisterID string,
) (
	AssisterFormView,
	error,
) {
	assisterForm, err := r.assisterFormRepository.FindOne_ByAssisterID(
		ctx,
		assisterID,
	)
	if err != nil {
		return AssisterFormView{}, err
	}

	return assisterForm.ToView(), err
}

func (r *assisterFormService) PutOne(
	ctx inner.Context,
	assisterForm AssisterForm,
) error {
	return r.assisterFormRepository.UpdateOne(ctx, assisterForm)
}

func NewAssisterFormService(
	assisterFormRepository AssisterFormRepository,
) AssisterFormService {
	return &assisterFormService{
		assisterFormRepository: assisterFormRepository,
	}
}
