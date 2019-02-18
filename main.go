package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"otk-final/hinny/web/api"

	"otk-final/hinny/module"
)

func main() {

	//数据库
	module.Install("mysql", "mycat-activeii:123qwe@(192.168.30.37:8066)/platform_behavior?charset=utf8")

	mux := mux.NewRouter()
	mux.Host("127.0.0.1").Name("业务自动化测试平台")

	/*接口服务*/
	mux.Path("/api/action/doc-fetch").Methods("POST").HandlerFunc(api.DocFetch)
	mux.Path("/api/action/list").Methods("GET").HandlerFunc(api.GetDBMetas)
	mux.Path("/api/{id}").Methods("GET")

	/*案例服务*/
	mux.Path("/case").Methods("POST")
	mux.Path("/case/action/list").Methods("GET")
	mux.Path("/case/{id}").Methods("GET")
	mux.Path("/case/{id}").Methods("PUT")

	/*调度服务*/
	mux.Path("/schedule/action/start").Methods("PUT")
	mux.Path("/schedule/action/pause").Methods("PUT")
	mux.Path("/schedule/action/stop").Methods("PUT")
	mux.Path("/schedule/action/get-process").Methods("GET")

	/*作业空间服务*/
	mux.Path("/workspace/action/list").Methods("GET")
	mux.Path("/workspace").Methods("POST")
	mux.Path("/workspace/{userId}").Methods("PUT")

	//启动服务
	err := http.ListenAndServe(":18080", mux)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
