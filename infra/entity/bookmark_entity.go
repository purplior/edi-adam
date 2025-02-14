package entity

import (
	"time"

	domain "github.com/purplior/sbec/domain/bookmark"
	"github.com/purplior/sbec/lib/dt"
)

type (
	Bookmark struct {
		ID          uint      `gorm:"primaryKey;autoIncrement"`
		UserID      uint      `gorm:"uniqueIndex:idx_user_assistant"`
		AssistantID uint      `gorm:"uniqueIndex:idx_user_assistant"`
		CreatedAt   time.Time `gorm:"autoCreateTime"`
		Assistant   Assistant
	}
)

func (e Bookmark) ToModel() domain.Bookmark {
	m := domain.Bookmark{
		CreatedAt:     e.CreatedAt,
		AssistantInfo: e.Assistant.ToModel().ToInfo(),
	}

	if e.ID > 0 {
		m.ID = dt.Str(e.ID)
	}
	if e.UserID > 0 {
		m.UserID = dt.Str(e.UserID)
	}
	if e.AssistantID > 0 {
		m.AssistantID = dt.Str(e.AssistantID)
	}

	return m
}

func MakeBookmark(m domain.Bookmark) Bookmark {
	e := Bookmark{
		CreatedAt: m.CreatedAt,
	}

	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}
	if len(m.UserID) > 0 {
		e.UserID = dt.UInt(m.UserID)
	}
	if len(m.AssistantID) > 0 {
		e.AssistantID = dt.UInt(m.AssistantID)
	}

	return e
}
