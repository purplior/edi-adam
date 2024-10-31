package entity

import (
	"time"

	"github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	User struct {
		ID              uint        `gorm:"primaryKey;autoIncrement"`
		JoinMethod      string      `gorm:"size:255;not null;uniqueIndex:idx_join_method_account"`
		AccountID       string      `gorm:"size:255;not null;uniqueIndex:idx_join_method_account"`
		AccountPassword string      `gorm:"size:255;not null"`
		Nickname        string      `gorm:"size:100"`
		Role            int         `gorm:"default:100"`
		CreatedAt       time.Time   `gorm:"autoCreateTime"`
		Assistants      []Assistant `gorm:"foreignKey:AuthorID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
		Challenges      []Challenge `gorm:"foreignKey:UserID;references:ID;constraint:OnUpdate:CASCADE,OnDelete:SET NULL;"`
	}
)

func (e User) ToModel() user.User {
	model := user.User{
		JoinMethod:      e.JoinMethod,
		AccountID:       e.AccountID,
		AccountPassword: e.AccountPassword,
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
		Nickname:        m.Nickname,
		Role:            m.Role,
		CreatedAt:       m.CreatedAt,
	}

	if len(m.ID) > 0 {
		entity.ID = dt.UInt(m.ID)
	}

	return entity
}
