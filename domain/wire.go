package domain

import (
	"github.com/google/wire"
	"github.com/purplior/edi-adam/domain/assistant"
	"github.com/purplior/edi-adam/domain/assister"
	"github.com/purplior/edi-adam/domain/auth"
	"github.com/purplior/edi-adam/domain/bookmark"
	"github.com/purplior/edi-adam/domain/category"
	"github.com/purplior/edi-adam/domain/customervoice"
	"github.com/purplior/edi-adam/domain/mission"
	"github.com/purplior/edi-adam/domain/missionlog"
	"github.com/purplior/edi-adam/domain/review"
	"github.com/purplior/edi-adam/domain/user"
	"github.com/purplior/edi-adam/domain/verification"
	"github.com/purplior/edi-adam/domain/wallet"
	"github.com/purplior/edi-adam/domain/walletlog"
)

var New = wire.NewSet(
	assistant.NewAssistantService,
	assister.NewAssisterService,
	auth.NewAuthService,
	bookmark.NewBookmarkService,
	category.NewCategoryService,
	customervoice.NewCustomerVoiceService,
	mission.NewMissionService,
	missionlog.NewMissionLogService,
	review.NewReviewService,
	user.NewUserService,
	verification.NewVerificationService,
	wallet.NewWalletService,
	walletlog.NewWalletLogService,
)
