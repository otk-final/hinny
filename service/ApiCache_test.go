package service

import (
	"testing"
	"encoding/json"
	"fmt"
	"otk-final/hinny/service/swagger"
)

func init() {
	ApiRefresh(&swagger.SwaggerHandler{}, "dev")
}

func TestGetDefinitionMap(t *testing.T) {
	out := GetDefinitionMap("dev", "Response«List«EsEntityVO»»")
	byte, err := json.Marshal(out)
	if err != nil {
		panic(err)
	}

	fmt.Println(string(byte))
}
