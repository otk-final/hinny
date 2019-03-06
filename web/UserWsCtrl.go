package web

import (
	"net/http"
	"otk-final/hinny/module/db"
	"io/ioutil"
	"encoding/json"
	"fmt"
	"github.com/gorilla/mux"
	"otk-final/hinny/service"
	"otk-final/hinny/service/swagger"
)

func GetWorkspaces(response http.ResponseWriter, request *http.Request) {
	allWs := make([]db.Workspace, 0)

	//查询数据库
	err := db.Conn.Find(&allWs)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	view.JSON(response, 200, allWs)
}

func CreateWorkspace(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	ws := &db.Workspace{}
	err := json.Unmarshal(body, ws)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	//新增数据库
	count, err := db.Conn.Insert(ws)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	fmt.Println("新增条目数:", count)
	view.JSON(response, 200, true)
}

func RemoveWorkspace(response http.ResponseWriter, request *http.Request) {
	vars := mux.Vars(request)
	//删除
	_, err := db.Conn.Id(vars["workspaceId"]).Delete(&db.Workspace{})
	if err != nil {
		view.JSON(response, 500, err)
		return
	}
	view.JSON(response, 200, true)
}

func RefreshWorkspace(response http.ResponseWriter, request *http.Request) {

	body, _ := ioutil.ReadAll(request.Body)
	ws := &db.Workspace{}
	err := json.Unmarshal(body, ws)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	err = service.ApiRefresh(&swagger.SwaggerHandler{}, ws.WsKey)
	if err != nil {
		view.JSON(response, 500, err)
	}

	view.JSON(response, 200, true)
}

func ChangeWorkspace(response http.ResponseWriter, request *http.Request) {

}
