package apiuser

import (
	"fmt"
	"time"
	"user-apiv2/data"

	"github.com/pkg/errors"
)

type UserService struct {
	session *Session
}

// func NewUserService(dbi interface{}) *UserService {
// 	return &UserService{
// 		db: dbi,
// 	}
// }

var _ data.UserService = &UserService{}

func (u *UserService) AddUser(usr *data.User) error {

	usr.CreatedDate = time.Now()
	db, ok := u.session.db.(data.DBType)
	if !ok {
		return data.ErrUsrDbUnreachable
	}
	db[usr.Username] = usr
	return nil
}

func (u *UserService) GetUser(usrname string) (*data.User, error) {
	db, ok := u.session.db.(data.DBType)
	if !ok {
		return nil, data.ErrUsrDbUnreachable
	}
	if usrname != "" {
		usr := db[usrname]
		if usr == nil {
			return nil, data.ErrUserNotFound
		}
		return usr, nil
	}
	return nil, data.ErrUserNameEmpty
}

func (u *UserService) DeleteUser(usrname string) error {

	db, ok := u.session.db.(data.DBType)
	if !ok {
		return data.ErrUsrDbUnreachable
	}
	usr, ok := db[usrname]
	if !ok {
		return data.ErrUserNotFound
	} else if usr.Username != usrname {
		return data.ErrUnauthorized
	}
	delete(db, usrname)
	return nil
}

func (u *UserService) ListUsers() ([]*data.User, error) {
	udb, ok := u.session.db.(data.DBType)
	if !ok {
		return nil, data.ErrUsrDbUnreachable
	}
	var db []*data.User
	for _, b := range udb {
		db = append(db, b)
	}
	return db, nil
}

func (u *UserService) UpdateUser(usr *data.User) error {

	db, ok := u.session.db.(data.DBType)
	if !ok {
		return data.ErrUsrDbUnreachable
	}
	// Only allow owner to update Product.
	userInDB, ok := db[usr.Username]
	if !ok {
		return fmt.Errorf("memory db: product not found with ID %v", usr.ID)
	} else if userInDB.Username != string(usr.Username) {
		return data.ErrUnauthorized
		//return fmt.Errorf("memory db: Non owner not allowed to update Product %v", b.ID)
	}
	if usr.Username == "" {
		return errors.New("memory db: product with unassigned ID passed into updateProduct")
	}
	db[usr.Username] = usr
	return nil
}
