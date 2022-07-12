package server

import (
	"DrawAndGuess/identity"
	"DrawAndGuess/storage"
	"fmt"
	"google/uuid"
)
import "errors"

func (us *UserSet) findUserByIdStr(idStr string) (*User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return us.findUserById(id)
}

func (us *UserSet) findUserById(id uuid.UUID) (*User, error) {
	for _, u := range us.users {
		if u.UserId == id {
			return u, nil
		}
	}
	return nil, errors.New("User Not Found for Uuid: " + id.String())
}

func (us *UserSet) findUserIndexById(id uuid.UUID) (int, error) {
	for i, u := range us.users {
		if u.UserId == id {
			return i, nil
		}
	}
	return -1, errors.New("User Not Found for Uuid: " + id.String())
}

func (us *UserSet) deleteUserById(id uuid.UUID) error {
	index, err := us.findUserIndexById(id)
	if err != nil {
		return err
	}
	us.users = append(us.users[:index], us.users[index+1:]...)
	return nil
}

func (us *UserSet) appendUser(u *User) error {
	if _, err := us.findUserById(u.UserId); err == nil {
		return errors.New("User with Uuid: " + u.UserId.String() + " Already Exists")
	}
	us.users = append(us.users, u)
	return nil
}

func (us *UserSet) deleteUser(u *User) {
	index, err := us.findUserIndexById(u.UserId)
	if err != nil {
		return
	}
	us.users = append(us.users[:index], us.users[index+1:]...)
}

func (us *UserSet) userReg(name string, psw string) (uuid.UUID, error) {
	rows, err := storage.NewQuery("select name from users where name = '" + name + "';")
	if err != nil {
		fmt.Println("During query")
		return uuid.Nil, err
	}

	if rows.Next() {
		return uuid.Nil, errors.New("User with name " + name + " already exists.")
	}

	id := uuid.New()
	_, err = storage.NewExec("insert into users values('" + name + "', '" + psw + "', '" + id.String() + "');")
	if err != nil {
		fmt.Println("During exec")
		return uuid.Nil, err
	}
	return id, nil
}

func (us *UserSet) userLogin(name string, psw string) (uuid.UUID, error) {
	ok, err := identity.IsUserAuthorised(name, psw)
	if err != nil {
		return uuid.Nil, err
	}

	if ok {

		rows, err := storage.NewQuery("select uuid from users where name = '" + name + "';")
		if err != nil {
			return uuid.Nil, err
		}
		if rows.Next() {
			var idStr string
			err = rows.Scan(&idStr)
			if err != nil {
				return uuid.Nil, err
			}

			id, err := uuid.Parse(idStr)
			if err != nil {
				return uuid.Nil, err
			}
			err = us.appendUser(&User{
				UserName: name,
				UserId:   id,
			})
			if err != nil {
				return uuid.Nil, err
			}
			return id, nil
		} else {
			return uuid.Nil, errors.New("no matching user")
		}

	} else {
		return uuid.Nil, errors.New("username or password is incorrect")
	}
}

func (us *UserSet) getUserNames() []string {
	names := make([]string, len(us.users))
	for _, u := range us.users {
		names = append(names, u.UserName)
	}
	return names
}
