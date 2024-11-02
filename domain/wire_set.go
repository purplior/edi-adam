package domain

import (
	"github.com/google/wire"
	assistant "github.com/podossaem/podoroot/domain/assistant/app"
	assister "github.com/podossaem/podoroot/domain/assister/app"
	assisterform "github.com/podossaem/podoroot/domain/assisterform/app"
	auth "github.com/podossaem/podoroot/domain/auth/app"
	challenge "github.com/podossaem/podoroot/domain/challenge/app"
	ledger "github.com/podossaem/podoroot/domain/ledger/app"
	me "github.com/podossaem/podoroot/domain/me/app"
	mission "github.com/podossaem/podoroot/domain/mission/app"
	user "github.com/podossaem/podoroot/domain/user/app"
	verification "github.com/podossaem/podoroot/domain/verification/app"
	wallet "github.com/podossaem/podoroot/domain/wallet/app"
)

var New = wire.NewSet(
	assistant.New,
	assister.New,
	assisterform.New,
	auth.New,
	challenge.New,
	ledger.New,
	me.New,
	mission.New,
	user.New,
	verification.New,
	wallet.New,
)
