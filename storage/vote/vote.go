package vote

import (
	"time"

	"github.com/google/uuid"
)

type Vote struct {
	ID        string
	MemID     string
	UserID    string
	Vote      int
	CreatedAt time.Time
}

func NewVote(memID, UserID string, vote int) *Vote {
	return &Vote{
		ID:        uuid.New().String(),
		MemID:     memID,
		UserID:    UserID,
		Vote:      vote,
		CreatedAt: time.Now(),
	}
}
