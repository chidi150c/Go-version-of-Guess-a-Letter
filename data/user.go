package data

import (
	"fmt"
	"math"
	"time"
)

type DBType map[string]*User

type UserID string

type User struct {
	Username      string
	Password      string
	Firstname     string
	Lastname      string
	ID            UserID
	Email         string
	DisplayName   string
	ImageURL      string
	Token         string
	Url           string
	Authenticated bool
	CreatedDate   time.Time
	Expiry        int64
	Level         string
}

type UserService interface {
	AddUser(*User) error
	GetUser(username string) (*User, error)
	DeleteUser(username string) error
	ListUsers() ([]*User, error)
	UpdateUser(*User) error
}

func (b *User) MustName() string {
	if b.DisplayName == "" {
		return "UnNamed"
	}
	return b.DisplayName
}

type Mmm []*User

var cout int

func (b Mmm) Count(m int) bool {
	fmt.Println("\ncout = ", cout)
	cout++
	if math.Mod(float64(m), 5) == 0 {
		return true
	}
	return false
}
