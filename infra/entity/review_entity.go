package entity

import (
	"time"

	domain "github.com/purplior/podoroot/domain/review"
	"github.com/purplior/podoroot/lib/dt"
)

type (
	Review struct {
		ID          uint `gorm:"primaryKey;autoIncrement"`
		AuthorID    uint
		Author      User
		AssistantID uint
		Content     string `gorm:"size:1500"`
		Score       float64
		CreatedAt   time.Time `gorm:"autoCreateTime"`
	}
)

func (e Review) ToModel() domain.Review {
	m := domain.Review{
		Content:   e.Content,
		Score:     e.Score,
		CreatedAt: e.CreatedAt,
	}
	if e.ID > 0 {
		m.ID = dt.Str(e.ID)
	}
	if e.AuthorID > 0 {
		m.AuthorID = dt.Str(e.AuthorID)
		m.Author = e.Author.ToModel()
	}
	if e.AssistantID > 0 {
		m.AssistantID = dt.Str(e.AssistantID)
	}

	return m
}

func MakeReview(m domain.Review) Review {
	e := Review{
		Content:   m.Content,
		Score:     m.Score,
		CreatedAt: m.CreatedAt,
	}
	if len(m.ID) > 0 {
		e.ID = dt.UInt(m.ID)
	}
	if len(m.AuthorID) > 0 {
		e.AuthorID = dt.UInt(m.AuthorID)
	}
	if len(m.AssistantID) > 0 {
		e.AssistantID = dt.UInt(m.AssistantID)
	}

	return e
}
