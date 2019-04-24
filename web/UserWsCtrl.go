package web

import (
	"encoding/json"
	"errors"
	"fmt"
	"github.com/gorilla/mux"
	"github.com/otk-final/hinny/module"
	"github.com/otk-final/hinny/service"
	"github.com/otk-final/hinny/swagger"
	"io/ioutil"
	"net/http"
	"strconv"
)

func GetWorkspaceFromHeader(request *http.Request) (*module.Workspace, error) {
	wsKid := request.Header.Get("workspace")
	if wsKid == "" {
		return nil, errors.New("工作空间不能为nil")
	}

	ws := &module.Workspace{}
	ok, err := module.Conn.ID(wsKid).Get(ws)
	if !ok || err != nil {
		return nil, errors.New("工作空间查询异常")
	}
	return ws, nil
}

func GetWorkspaces(response http.ResponseWriter, request *http.Request) {
	allWs := make([]module.Workspace, 0)

	//查询数据库
	err := module.Conn.Find(&allWs)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	view.JSON(response, 200, allWs)
}

func CreateWorkspace(response http.ResponseWriter, request *http.Request) {
	body, _ := ioutil.ReadAll(request.Body)

	ws := &module.Workspace{}
	err := json.Unmarshal(body, ws)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	//唯一标识
	wsKid, _ := module.GetNextKid()
	ws.Kid = wsKid

	//新增数据库
	count, err := module.Conn.Insert(ws)
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
	_, err := module.Conn.Id(kid).Delete(&module.Workspace{})
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
	ws := &module.Workspace{}
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
	ws := &module.Workspace{}
	err := json.Unmarshal(body, ws)
	if err != nil {
		view.JSON(response, 500, err)
		return
	}

	count, err := module.Conn.Id(ws.Kid).Update(module.Workspace{Script: ws.Script})
	if count != 1 || err != nil {
		panic(err)
	}

	view.JSON(response, 200, "成功")
}
