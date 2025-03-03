package session

import (
	"context"

	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/domain/shared/inner"
	"github.com/purplior/edi-adam/domain/shared/model"
	"github.com/purplior/edi-adam/infra/database/postgre"
	"gorm.io/gorm"
)

type (
	_session struct {
		identity      *inner.Identity
		innerContext  context.Context
		postgreClient *postgre.Client
		postgreTX     *gorm.DB
	}
)

func (s *_session) Clear() {
	s.identity = nil
	s.innerContext = nil
	s.postgreClient = nil
	s.postgreTX = nil
}

func (s *_session) Identity() *inner.Identity {
	return s.identity
}

func (s *_session) SetIdentity(identity *inner.Identity) {
	s.identity = identity
}

func (s *_session) Context() context.Context {
	return s.innerContext
}

func (s *_session) Transaction() *gorm.DB {
	return s.postgreTX
}

func (s *_session) BeginTransaction() error {
	if s.postgreTX != nil {
		return exception.ErrInTransaction
	}

	s.postgreTX = s.postgreClient.WithContext(s.innerContext).Begin()

	return nil
}

func (s *_session) CommitTransaction() error {
	if s.postgreTX == nil {
		return exception.ErrInvalidTransaction
	}

	err := s.postgreTX.Commit().Error
	if err != nil {
		s.postgreTX.Rollback()
	}
	s.postgreTX = nil

	return err
}

func (s *_session) RollbackTransaction() {
	if s.postgreTX != nil {
		s.postgreTX.Rollback()
		s.postgreTX = nil
	}
}

func (s *_session) ReleaseTransaction() {
	if s.postgreTX != nil {
		s.postgreTX.Rollback()
		s.postgreTX = nil
	}
}

func (s *_session) GuardMemberAuth() error {
	if s.identity == nil || s.identity.ID == 0 {
		return exception.ErrUnauthorized
	}

	return nil
}

func (s *_session) GuardAdminAuth() error {
	if s.identity == nil || s.identity.ID == 0 || s.identity.Role != model.UserRole_Admin {
		return exception.ErrUnauthorized
	}

	return nil
}
