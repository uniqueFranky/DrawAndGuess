package server

import (
	"fmt"
	"strings"
)

func (mb *MyBytes) MarshalJSON() ([]byte, error) {
	var array string
	if mb == nil {
		array = "null"
	} else {
		array = strings.Join(strings.Fields(fmt.Sprintf("%d", *mb)), ",")
	}
	jsonResult := fmt.Sprintf(`%s`, array)
	return []byte(jsonResult), nil
}
