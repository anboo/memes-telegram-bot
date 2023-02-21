package mem

type ReservedMemUser struct {
	MemID  string
	UserID string
}

func (ReservedMemUser) TableName() string {
	return "reserved_memes"
}
