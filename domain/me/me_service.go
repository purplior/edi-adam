package me

import (
	"github.com/podossaem/podoroot/domain/user"
)

type (
	MeService interface {
	}
)

type (
	meService struct {
		userRepository user.UserRepository
	}
)

func NewMeService(
	userRepository user.UserRepository,
) MeService {
	return &meService{
		userRepository: userRepository,
	}
}
