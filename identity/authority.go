package identity

import (
	"DrawAndGuess/storage"
	"errors"
	"fmt"
	"google/uuid"
)

func IsUserAuthorised(name string, psw string) (bool, error) {

	rows, err := storage.NewQuery("select psw from users where name = '" + name + "';")
	if err != nil {
		fmt.Println("During query")
		return false, err
	}
	if rows.Next() {
		var realPsw string
		err = rows.Scan(&realPsw)
		fmt.Println(psw, realPsw)
		if err != nil {
			return false, err
		}

		if realPsw != psw {
			return false, errors.New("password is incorrect")
		}
	} else {
		return false, errors.New("no matching user with name " + name)
	}
	return true, nil
}

func IsIdStrValid(name string, idStr string) (bool, error) {
	rows, err := storage.NewQuery("select uuid from users where name = '" + name + "';")
	if err != nil {
		return false, err
	}
	if rows.Next() {
		var realIdStr string
		err = rows.Scan(&realIdStr)
		if err != nil {
			return false, err
		}

		if realIdStr != idStr {
			return false, errors.New("uuid not match with real uuid")
		}

		return true, nil
	} else {
		return false, errors.New("no matching user for username: " + name)
	}
}

func IsIdValid(name string, id uuid.UUID) (bool, error) {
	return IsIdStrValid(name, id.String())
}