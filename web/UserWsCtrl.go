package web

import (
	"net/http"
)

type Workspace struct {
	Name string `json:"name"`
	Key  string `json:"key"`
	Host string `json:"host"`
}

var WsCaches = make([]Workspace, 0)

func init() {
	WsCaches = append(WsCaches, Workspace{
		Name: "黄祥云Dev",
		Key:  "huangxy-dev",
		Host: "192.168.30.61:18080",
	}, Workspace{
		Name: "开发环境",
		Key:  "dev",
		Host: "http://api-dev.yryz.com/gateway/lovelorn",
	}, Workspace{
		Name: "李凡",
		Key:  "lifan-dev",
		Host: "192.168.30.23:19000",
	})
}



func GetWorkspaces(response http.ResponseWriter, request *http.Request) {
	view.JSON(response, 200, WsCaches)
}

func CreateWorkspace(response http.ResponseWriter, request *http.Request) {

}

func RemoveWorkspace(response http.ResponseWriter, request *http.Request) {

}

func ChangeWorkspace(response http.ResponseWriter, request *http.Request) {

}

func FindWorkspace(key string) *Workspace {
	for _, e := range WsCaches {
		if e.Key == key {
			return &e
		}
	}
	return nil
}
