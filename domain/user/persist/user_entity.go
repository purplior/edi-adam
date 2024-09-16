package persist

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/lib/dt"
)

type (
	User struct {
		ID              uint      `gorm:"primaryKey;autoIncrement"`
		JoinMethod      string    `gorm:"type:varchar(255);uniqueIndex:idx_join_method_account"`
		AccountID       string    `gorm:"type:varchar(255);not null;uniqueIndex:idx_join_method_account"`
		AccountPassword string    `gorm:"type:varchar(255);not null"`
		Nickname        string    `gorm:"type:varchar(100)"`
		Role            int       `gorm:"default:100"`
		CreatedAt       time.Time `gorm:"autoCreateTime"`
	}
)

func (e User) ToModel() domain.User {
	model := domain.User{
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

func MakeUser(m domain.User) User {
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
