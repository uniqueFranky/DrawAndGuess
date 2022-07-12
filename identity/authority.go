package identity

import (
	"DrawAndGuess/storage"
	"errors"
	"fmt"
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
