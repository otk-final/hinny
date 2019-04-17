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
	"strconv"
	"errors"
)

func GetWorkspaceFromHeader(request *http.Request) (*db.Workspace, error) {
	wsKid := request.Header.Get("workspace")
	if wsKid == "" {
		return nil, errors.New("工作空间不能为nil")
	}

	ws := &db.Workspace{}
	ok, err := db.Conn.ID(wsKid).Get(ws)
	if !ok || err != nil {
		return nil, errors.New("工作空间查询异常")
	}
	return ws, nil
}

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

	//唯一标识
	wsKid, _ := db.GetNextKid()
	ws.Kid = wsKid

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
	kid := vars["kid"]
	//删除
	_, err := db.Conn.Id(kid).Delete(&db.Workspace{})
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	// 清除缓存
	uintKid, _ := strconv.ParseUint(kid, 10, 64)
	service.ApiRemove(uintKid)

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

	err = service.ApiRefresh(&swagger.SwaggerHandler{}, ws)
	if err != nil {
		view.JSON(response, 500, err)
	}

	view.JSON(response, 200, true)
}

func UpdateScript(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)
	ws := &db.Workspace{}
	err := json.Unmarshal(body, ws)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	count, err := db.Conn.Id(ws.Kid).Update(db.Workspace{Script: ws.Script})
	if count != 1 || err != nil {
		panic(err)
	}

	view.JSON(response, 200, "成功")
}
