package me

type (
	MeService interface {
	}
)

type (
	meService struct {
	}
)

func NewMeService() MeService {
	return &meService{}
}
