package server

import "google/uuid"
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

func (us *UserSet) appendUserWithName(name string) error {
	id := uuid.New()
	us.users = append(us.users, &User{
		UserId:   id,
		UserName: name,
	})
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
