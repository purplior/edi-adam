package domain

import (
	"github.com/google/wire"
	assistant "github.com/podossaem/podoroot/domain/assistant/app"
	auth "github.com/podossaem/podoroot/domain/auth/app"
	me "github.com/podossaem/podoroot/domain/me/app"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

var New = wire.NewSet(
	assistant.New,
	auth.New,
	me.New,
	user.New,
	verification.New,
)
