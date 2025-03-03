package model

import "time"

type (
	// @MODEL {지갑}
	Wallet struct {
		ID      uint  `gorm:"primaryKey;autoIncrement" json:"id,omitempty"`
		OwnerID uint  `gorm:"unique" json:"ownerId"`
		Coin    int64 `gorm:"not null" json:"coin"`

		CreatedAt time.Time `gorm:"autoCreateTime" json:"createdAt"`
		UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updatedAt"`
	}
)
