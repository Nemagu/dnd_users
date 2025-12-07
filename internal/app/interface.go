package app

type emailValidator interface {
	Validate(email string) error
}

type passwordValidator interface {
	Validate(password, email string) error
}

type passwordHasher interface {
	Hash(password string) (string, error)
}

type passwordComparer interface {
	Compare(password, hash string) (bool, error)
}

type codeGenerator interface {
	Generate() string
}
