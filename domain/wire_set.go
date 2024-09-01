package domain

import (
	"github.com/google/wire"
	auth "github.com/podossaem/podoroot/domain/auth/app"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

var New = wire.NewSet(
	auth.New,
	user.New,
	verification.New,
)
