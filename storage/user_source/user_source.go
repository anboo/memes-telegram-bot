package user_source

type UserSource struct {
	UserID  string
	Source  string
	Enabled bool
}

func (UserSource) TableName() string {
	return "users_sources"
}

func New(userID string, source string) *UserSource {
	return &UserSource{
		UserID:  userID,
		Source:  source,
		Enabled: true,
	}
}
