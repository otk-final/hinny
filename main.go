package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"otk-final/hinny/web"
	"otk-final/hinny/module/db"
)

func main() {

	//数据库
	db.Install("mysql", "dev62:dev62.123456@(192.168.30.62:3306)/platform_behavior?charset=utf8")

	router := mux.NewRouter()
	router.Host("127.0.0.1").Name("业务自动化测试平台")

	/*服务接口*/

	/*路径服务*/
	router.Path("/service/action/list").Methods("OPTIONS", "GET").HandlerFunc(web.GetServices)
	router.Path("/path/action/list").Methods("OPTIONS", "GET").HandlerFunc(web.GetPaths)
	router.Path("/path/action/primary").Methods("OPTIONS", "GET").HandlerFunc(web.GetPrimaryPath)
	router.Path("/path/action/execute").Methods("POST")
	router.Path("/path/action/submit").Methods("POST")

	/*作业空间服务*/
	router.Path("/workspace/action/list").Methods("OPTIONS", "GET").HandlerFunc(web.GetWorkspaces)
	router.Path("/workspace").Methods("POST").HandlerFunc(web.CreateWorkspace)
	//router.Path("/workspace/{key}").Methods("DELETE").HandlerFunc(web.RemoveWorkspace)
	router.Path("/workspace/action/change").Methods("POST").HandlerFunc(web.ChangeWorkspace)

	//全部支持跨域
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "workspace")
			next.ServeHTTP(w, req)
		})
	})
	router.Use(mux.CORSMethodMiddleware(router))

	//启动服务
	err := http.ListenAndServe(":18080", router)
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
