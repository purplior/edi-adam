package persist

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/user"
	"github.com/podossaem/podoroot/lib/mydate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	User struct {
		ID              primitive.ObjectID `bson:"_id,omitempty"`
		JoinMethod      string             `bson:"join_method"`
		AccountID       string             `bson:"account_id"`
		AccountPassword string             `bson:"account_password"`
		Nickname        string             `bson:"nickname"`
		Role            int                `bson:"role"`
		CreatedAt       time.Time          `bson:"created_at"`
	}
)

func (e User) BeforeInsert() User {
	e.CreatedAt = mydate.Now()

	return e
}

func (e User) ToModel() domain.User {
	id := ""
	if !e.ID.IsZero() {
		id = e.ID.Hex()
	}

	return domain.User{
		ID:              id,
		JoinMethod:      e.JoinMethod,
		AccountID:       e.AccountID,
		AccountPassword: e.AccountPassword,
		Nickname:        e.Nickname,
		Role:            e.Role,
		CreatedAt:       e.CreatedAt,
	}
}

func MakeUser(m domain.User) User {
	oid, _ := primitive.ObjectIDFromHex(m.ID)

	return User{
		ID:              oid,
		JoinMethod:      m.JoinMethod,
		AccountID:       m.AccountID,
		AccountPassword: m.AccountPassword,
		Nickname:        m.Nickname,
		Role:            m.Role,
		CreatedAt:       m.CreatedAt,
	}
}
