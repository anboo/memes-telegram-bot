package mem

import (
	"time"

	"github.com/google/uuid"
)

type Mem struct {
	ID         string
	Text       string
	ExternalID string
	Source     string
	SourceFrom string
	Img        string
	Rating     int
	CreatedAt  time.Time
}

func (Mem) TableName() string {
	return "memes"
}

func NewMem(externalID string, text string, source string, sourceFrom string, img string) *Mem {
	return &Mem{
		ID:         uuid.New().String(),
		Text:       text,
		ExternalID: externalID,
		Source:     source,
		SourceFrom: sourceFrom,
		Img:        img,
		Rating:     0,
		CreatedAt:  time.Now(),
	}
}
