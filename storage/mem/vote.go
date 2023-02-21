package mem

import (
	"time"
)

type Vote struct {
	ID        string
	MemID     string
	UserID    string
	Vote      int
	CreatedAt time.Time
}
