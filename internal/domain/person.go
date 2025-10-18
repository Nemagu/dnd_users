package domain

import "fmt"

type Person struct {
	firstName string
	lastName  string
}

func NewPerson(firstName, lastName string) (Person, error) {
	return Person{firstName: firstName, lastName: lastName}, nil
}

func (p Person) FirstName() string {
	return p.firstName
}

func (p Person) LastName() string {
	return p.lastName
}

func (p Person) String() string {
	return fmt.Sprintf("first name: %s\nlast name: %s", p.firstName, p.lastName)
}
