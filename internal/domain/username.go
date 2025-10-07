package domain

type Username struct {
	username string
}

func NewUsername(username string) (Username, error) {
	// TODO: add validation username
	return Username{username: username}, nil
}

func (un Username) Username() string {
	return un.username
}
