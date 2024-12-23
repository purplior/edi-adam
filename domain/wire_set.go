package domain

import (
	"github.com/google/wire"
	assistant "github.com/purplior/podoroot/domain/assistant/app"
	assister "github.com/purplior/podoroot/domain/assister/app"
	assisterform "github.com/purplior/podoroot/domain/assisterform/app"
	auth "github.com/purplior/podoroot/domain/auth/app"
	challenge "github.com/purplior/podoroot/domain/challenge/app"
	customervoice "github.com/purplior/podoroot/domain/customervoice/app"
	ledger "github.com/purplior/podoroot/domain/ledger/app"
	me "github.com/purplior/podoroot/domain/me/app"
	mission "github.com/purplior/podoroot/domain/mission/app"
	user "github.com/purplior/podoroot/domain/user/app"
	verification "github.com/purplior/podoroot/domain/verification/app"
	wallet "github.com/purplior/podoroot/domain/wallet/app"
)

var New = wire.NewSet(
	assistant.New,
	assister.New,
	assisterform.New,
	auth.New,
	challenge.New,
	customervoice.New,
	ledger.New,
	me.New,
	mission.New,
	user.New,
	verification.New,
	wallet.New,
)
