package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"otk-final/hinny/web"
	"otk-final/hinny/module/db"
	"time"
)

func main() {

	//数据库
	db.Install("mysql", "dev62:dev62.123456@(192.168.30.62:3306)/platform_behavior?charset=utf8")
	//ID生成规则
	db.InstallIDGeneral(time.Now(), 11)

	router := mux.NewRouter()
	router.Host("127.0.0.1").Name("业务自动化测试平台")

	router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	/*服务接口*/
	router.Path("/service/action/list").Methods("GET").HandlerFunc(web.GetServices)

	/*路径服务*/
	router.Path("/path/action/list").Methods("GET").HandlerFunc(web.GetPaths)
	router.Path("/path/action/primary").Methods("GET").HandlerFunc(web.GetPrimaryPath)
	router.Path("/path/action/execute").Methods("POST")
	router.Path("/path/action/submit").Methods("POST")

	/*作业空间服务*/
	router.Path("/workspace/action/list").Methods("GET").HandlerFunc(web.GetWorkspaces)
	router.Path("/workspace").Methods("POST").HandlerFunc(web.CreateWorkspace)
	router.Path("/workspace/{kid}").Methods("DELETE").HandlerFunc(web.RemoveWorkspace)
	router.Path("/workspace/action/refresh").Methods("POST").HandlerFunc(web.RefreshWorkspace)
	router.Path("/workspace/action/update-script").Methods("POST").HandlerFunc(web.UpdateScript)

	router.Path("/case/action/get-modules").Methods("GET").HandlerFunc(web.GetCaseModules)
	router.Path("/case/action/execute").Methods("POST").HandlerFunc(web.CaseExecute)
	router.Path("/case/action/save").Methods("POST").HandlerFunc(web.CaseSave)
	router.Path("/case/action/list").Methods("GET").HandlerFunc(web.GetCases)
	router.Path("/case/action/get-logs").Methods("GET").HandlerFunc(web.GetCaseLogs)
	router.Path("/case/{kid}").Methods("GET").HandlerFunc(web.GetCaseLog)

	//TODO 验证脚本，默认参数

	//全部支持跨域
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "content-type,workspace,application")
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
