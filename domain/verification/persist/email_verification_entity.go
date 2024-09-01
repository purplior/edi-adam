package persist

import (
	"time"

	domain "github.com/podossaem/podoroot/domain/verification"
	"github.com/podossaem/podoroot/lib/mydate"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type (
	EmailVerification struct {
		ID         primitive.ObjectID `bson:"_id,omitempty"`
		Email      string             `bson:"email"`
		Code       string             `bson:"code"`
		IsVerified bool               `bson:"is_verified"`
		IsConsumed bool               `bson:"is_consumed"`
		ExpiredAt  time.Time          `bson:"expired_at"`
		CreatedAt  time.Time          `bson:"created_at"`
	}
)

func (e EmailVerification) BeforeInsert() EmailVerification {
	e.CreatedAt = mydate.Now()

	return e
}

func (e EmailVerification) ToModel() domain.EmailVerification {
	id := ""
	if !e.ID.IsZero() {
		id = e.ID.Hex()
	}

	return domain.EmailVerification{
		ID:         id,
		Email:      e.Email,
		Code:       e.Code,
		IsVerified: e.IsVerified,
		IsConsumed: e.IsConsumed,
		ExpiredAt:  e.ExpiredAt,
		CreatedAt:  e.CreatedAt,
	}
}

func MakeEmailVerification(m domain.EmailVerification) EmailVerification {
	oid, _ := primitive.ObjectIDFromHex(m.ID)

	return EmailVerification{
		ID:         oid,
		Email:      m.Email,
		Code:       m.Code,
		IsVerified: m.IsVerified,
		IsConsumed: m.IsConsumed,
		ExpiredAt:  m.ExpiredAt,
		CreatedAt:  m.CreatedAt,
	}
}
