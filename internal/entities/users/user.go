package users

import "fmt"

type User struct {
	ID    int64
	Email string `validate:"min:3; max:320"`
	Name  string `validate:"min:2; max:99"`
}

func (u User) String() string {
	return fmt.Sprintf(
		"<User id=%d email=%s name=`%s`>",
		u.ID,
		u.Email,
		u.Name,
	)
}

func New(email string, name string) User {
	return User{
		Email: email,
		Name:  name,
	}
}
