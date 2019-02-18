package api

import (
	"fmt"
	"net/http"
	"otk-final/hinny/module"
)


func GetDBMetas(response http.ResponseWriter, request *http.Request) {

	data, err := module.DB.DBMetas()

	if err != nil {

	}

	fmt.Println(&data)
}
