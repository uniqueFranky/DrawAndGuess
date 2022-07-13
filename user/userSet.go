package user

import (
	"DrawAndGuess/identity"
	"DrawAndGuess/storage"
	"fmt"
	"google/uuid"
)
import "errors"

type UserSet struct {
	Users []*User
}

func (us *UserSet) findUserByIdStr(idStr string) (*User, error) {
	id, err := uuid.Parse(idStr)
	if err != nil {
		return nil, err
	}
	return us.FindUserById(id)
}

func (us *UserSet) findUserByName(name string) (*User, error) {
	for _, u := range us.Users {
		if u.UserName == name {
			return u, nil
		}
	}
	return nil, errors.New("no matching user for name: " + name)
}

func (us *UserSet) findUserIndexByName(name string) (int, error) {
	for i, u := range us.Users {
		if u.UserName == name {
			return i, nil
		}
	}
	return -1, errors.New("no matching user for name: " + name)
}

func (us *UserSet) FindUserById(id uuid.UUID) (*User, error) {
	for _, u := range us.Users {
		if u.UserId == id {
			return u, nil
		}
	}
	return nil, errors.New("User Not Found for Uuid: " + id.String())
}

func (us *UserSet) findUserIndexById(id uuid.UUID) (int, error) {
	for i, u := range us.Users {
		if u.UserId == id {
			return i, nil
		}
	}
	return -1, errors.New("User Not Found for Uuid: " + id.String())
}

func (us *UserSet) DeleteUserById(id uuid.UUID) error {
	index, err := us.findUserIndexById(id)
	if err != nil {
		return err
	}
	us.Users = append(us.Users[:index], us.Users[index+1:]...)
	return nil
}

func (us *UserSet) AppendUser(u *User) error {
	if _, err := us.FindUserById(u.UserId); err == nil {
		return errors.New("User with Uuid: " + u.UserId.String() + " Already Exists")
	}
	us.Users = append(us.Users, u)
	return nil
}

func (us *UserSet) deleteUser(u *User) {
	index, err := us.findUserIndexById(u.UserId)
	if err != nil {
		return
	}
	us.Users = append(us.Users[:index], us.Users[index+1:]...)
}

func (us *UserSet) UserReg(name string, psw string) (uuid.UUID, error) {
	rows, err := storage.NewQuery("select name from users where name = '" + name + "';")
	if err != nil {
		fmt.Println("During query")
		return uuid.Nil, err
	}
	defer rows.Close()

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

func (us *UserSet) UserLogin(name string, psw string) (uuid.UUID, error) {
	ok, err := identity.IsUserAuthorised(name, psw)
	if err != nil {
		return uuid.Nil, err
	}

	if ok {

		rows, err := storage.NewQuery("select uuid from users where name = '" + name + "';")
		if err != nil {
			return uuid.Nil, err
		}
		defer rows.Close()
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
			err = us.AppendUser(&User{
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

func (us *UserSet) GetUserNames() []string {
	names := []string{}
	for _, u := range us.Users {
		fmt.Println(u.UserName)
		names = append(names, u.UserName)
	}
	return names
}
