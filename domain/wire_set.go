package domain

import (
	"github.com/google/wire"
	assistant "github.com/purplior/sbec/domain/assistant/app"
	assister "github.com/purplior/sbec/domain/assister/app"
	auth "github.com/purplior/sbec/domain/auth/app"
	bookmark "github.com/purplior/sbec/domain/bookmark/app"
	category "github.com/purplior/sbec/domain/category/app"
	challenge "github.com/purplior/sbec/domain/challenge/app"
	customervoice "github.com/purplior/sbec/domain/customervoice/app"
	ledger "github.com/purplior/sbec/domain/ledger/app"
	me "github.com/purplior/sbec/domain/me/app"
	mission "github.com/purplior/sbec/domain/mission/app"
	review "github.com/purplior/sbec/domain/review/app"
	user "github.com/purplior/sbec/domain/user/app"
	verification "github.com/purplior/sbec/domain/verification/app"
	wallet "github.com/purplior/sbec/domain/wallet/app"
)

var New = wire.NewSet(
	assistant.New,
	assister.New,
	auth.New,
	bookmark.New,
	category.New,
	challenge.New,
	customervoice.New,
	ledger.New,
	me.New,
	mission.New,
	review.New,
	user.New,
	verification.New,
	wallet.New,
)
