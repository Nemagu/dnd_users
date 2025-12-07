package app

type emailValidator interface {
	Validate(email string) error
}

type passwordValidator interface {
	Validate(password, email string) error
}

type codeGenerator interface {
	Generate() string
}
