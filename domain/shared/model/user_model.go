package model

import (
	"time"

	"github.com/purplior/edi-adam/application/config"
	"github.com/purplior/edi-adam/domain/shared/exception"
	"github.com/purplior/edi-adam/lib/dt"
	"github.com/purplior/edi-adam/lib/security"
)

type (
	// @ENUM {가입방법}
	UserJoinMethod string

	// @ENUM {권한}
	UserRole int
)

const (
	UserJoinMethod_Phone UserJoinMethod = "phone"

	UserRole_Member UserRole = 100
	UserRole_Admin  UserRole = 10000
)

type (
	// @DTO {아바타 JSON 객체 문자열}
	UserAvatar struct {
		Theme string `json:"theme"`
		Text  string `json:"text"`
	}

	// @MODEL {사용자 모델}
	User struct {
		ID uint `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		// 가입방법 & 전화번호 등이 암호화 되서 ID로 저장됨
		// { "method": "phone", "phone_number": "01012345678" }
		AccountID        string     `gorm:"size:255" json:"accountId"`
		Nickname         string     `gorm:"size:255" json:"nickname"`
		Avatar           UserAvatar `gorm:"serializer:json" json:"avatar"`
		Role             UserRole   `gorm:"default:100" json:"role"`
		IsMarketingAgree bool       `json:"isMarketingAgree"`
		IsInactivated    bool       `json:"isInactivated"`

		CreatedAt     time.Time  `gorm:"autoCreateTime" json:"createdAt"`
		InactivatedAt *time.Time `json:"inactivatedAt"`

		Assistants     []Assistant     `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;" json:"assistants,omitempty"`
		Bookmarks      []Bookmark      `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"bookmarks,omitempty"`
		MissionLogs    []MissionLog    `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"challenges,omitempty"`
		CustomerVoices []CustomerVoice `gorm:"foreignKey:UserID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"customervoices,omitempty"`
		Reviews        []Review        `gorm:"foreignKey:AuthorID;constraint:OnUpdate:CASCADE,OnDelete:CASCADE;" json:"reviews,omitempty"`
		Wallet         Wallet          `gorm:"foreignKey:OwnerID;constraint:OnUpdate:CASCADE,OnDelete:RESTRICT;" json:"wallet,omitempty"`
	}

	// @MODEL {사용자 모델 : 프로필}
	UserProfile struct {
		ID       uint       `json:"id"`
		Avatar   UserAvatar `json:"avatar"`
		Nickname string     `json:"nickname"`
	}

	// @MODEL {사용자 모델 : 상세정보}
	UserDetail struct {
		ID          uint       `json:"id"`
		Avatar      UserAvatar `json:"avatar"`
		Nickname    string     `json:"nickname"`
		PhoneNumber string     `json:"phoneNumber"`

		CreatedAt time.Time `json:"createdAt"`
	}
)

func (m User) ToProfile() UserProfile {
	return UserProfile{
		ID:       m.ID,
		Avatar:   m.Avatar,
		Nickname: m.Nickname,
	}
}

func (m User) ToDetail() (UserDetail, error) {
	data, err := security.DecryptMapDataWithAESGCM(m.AccountID, config.SymmetricKey())
	if err != nil {
		return UserDetail{}, err
	}
	phoneNumber := dt.Str(data["phone_number"])
	if len(phoneNumber) == 0 {
		return UserDetail{}, exception.ErrBadRequest
	}

	return UserDetail{
		ID:          m.ID,
		Avatar:      m.Avatar,
		Nickname:    m.Nickname,
		PhoneNumber: phoneNumber,

		CreatedAt: m.CreatedAt,
	}, nil
}
