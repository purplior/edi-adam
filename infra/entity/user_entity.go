package entity

import (
	"time"

	"github.com/purplior/podoroot/domain/user"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	User struct {
		ID              uint            `gorm:"primaryKey;autoIncrement"`
		JoinMethod      string          `gorm:"size:255;not null;uniqueIndex:idx_join_method_account"`
		AccountID       string          `gorm:"size:255;not null;uniqueIndex:idx_join_method_account"`
		AccountPassword string          `gorm:"size:255;not null"`
		AvatarTheme     int             `gorm:"default:1"`
		AvatarText      string          `gorm:"size:10"`
		Nickname        string          `gorm:"size:100"`
		Role            int             `gorm:"default:100"`
		CreatedAt       time.Time       `gorm:"autoCreateTime"`
		Wallet          Wallet          `gorm:"foreignKey:OwnerID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Categories      []Category      `gorm:"foreignKey:CreatorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Assistants      []Assistant     `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Challenges      []Challenge     `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		CustomerVoices  []CustomerVoice `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e User) ToModel() user.User {
	model := user.User{
		JoinMethod:      e.JoinMethod,
		AccountID:       e.AccountID,
		AccountPassword: e.AccountPassword,
		AvatarTheme:     e.AvatarTheme,
		AvatarText:      e.AvatarText,
		Nickname:        e.Nickname,
		Role:            e.Role,
		CreatedAt:       e.CreatedAt,
	}

	if e.ID > 0 {
		model.ID = dt.Str(e.ID)
	}

	return model
}

func MakeUser(m user.User) User {
	entity := User{
		JoinMethod:      m.JoinMethod,
		AccountID:       m.AccountID,
		AccountPassword: m.AccountPassword,
		AvatarTheme:     m.AvatarTheme,
		AvatarText:      m.AvatarText,
		Nickname:        m.Nickname,
		Role:            m.Role,
		CreatedAt:       m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}

	return entity
}
