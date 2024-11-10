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

		CreateOne(
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

func (s *assisterFormService) RegisterOne(
	ctx inner.Context,
	request AssisterFormRegisterRequest,
) (
	AssisterForm,
	error,
) {
	return s.assisterFormRepository.InsertOne(
		ctx,
		AssisterForm{
			Fields:        request.Fields,
			QueryMessages: request.QueryMessages,
		},
	)
}

func (s *assisterFormService) GetOne_ByID(
	ctx inner.Context,
	id string,
) (
	AssisterForm,
	error,
) {
	return s.assisterFormRepository.FindOne_ByID(
		ctx,
		id,
	)
}

func (s *assisterFormService) GetOne_ByAssisterID(
	ctx inner.Context,
	assisterID string,
) (
	AssisterForm,
	error,
) {
	return s.assisterFormRepository.FindOne_ByAssisterID(
		ctx,
		assisterID,
	)
}

func (s *assisterFormService) GetViewOne_ByAssister(
	ctx inner.Context,
	assisterID string,
) (
	AssisterFormView,
	error,
) {
	assisterForm, err := s.assisterFormRepository.FindOne_ByAssisterID(
		ctx,
		assisterID,
	)
	if err != nil {
		return AssisterFormView{}, err
	}

	return assisterForm.ToView(), err
}

func (s *assisterFormService) PutOne(
	ctx inner.Context,
	assisterForm AssisterForm,
) error {
	return s.assisterFormRepository.UpdateOne(ctx, assisterForm)
}

func (s *assisterFormService) CreateOne(
	ctx inner.Context,
	assisterForm AssisterForm,
) error {
	_, err := s.assisterFormRepository.InsertOne(ctx, assisterForm)

	return err
}

func NewAssisterFormService(
	assisterFormRepository AssisterFormRepository,
) AssisterFormService {
	return &assisterFormService{
		assisterFormRepository: assisterFormRepository,
	}
}
