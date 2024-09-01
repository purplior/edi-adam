package domain

import (
	"github.com/google/wire"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
)

var New = wire.NewSet(
	user.New,
	verification.New,
)
