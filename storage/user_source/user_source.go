package user_source

type UserSource struct {
	MemID   string
	Source  string
	Enabled bool
}

func (UserSource) TableName() string {
	return "users_sources"
}

func New(memID string, source string) *UserSource {
	return &UserSource{
		MemID:   memID,
		Source:  source,
		Enabled: true,
	}
}
