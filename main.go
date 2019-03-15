package main

import (
	"github.com/gorilla/mux"
	"net/http"
	"log"
	"otk-final/hinny/web"
	"fmt"
	"otk-final/hinny/module/db"
	"otk-final/hinny/module/global"
)

func init() {

	application := global.Conf.GetString("application")
	fmt.Println(application)

	initDB()

	initIDGeneral()
}

//加载数据库配置文件
func initDB() {
	dbConf := global.Conf.GetStringMapString("db")
	dbUrl := fmt.Sprintf("%s:%s@(%s)/%s", dbConf["username"], dbConf["password"], dbConf["host"], dbConf["url"])
	fmt.Printf("数据库:%s", dbUrl)
	//数据库
	db.Install(dbConf["driver"], dbUrl)

}

//加载雪花算法配置
func initIDGeneral() {

	machineID := global.Conf.GetInt("snowflake.machineID")
	startTime := global.Conf.GetTime("snowflake.startTime")

	//ID生成规则
	db.InstallIDGeneral(startTime, uint16(machineID))
}

//加载web请求路径配置
func initWebCtrl() *mux.Router {

	router := mux.NewRouter().StrictSlash(true)

	router.PathPrefix("/").Methods("OPTIONS").HandlerFunc(func(writer http.ResponseWriter, request *http.Request) {
		writer.WriteHeader(200)
	})

	//静态资源
	router.PathPrefix("/static/").Handler(http.StripPrefix("/static/", http.FileServer(http.Dir("./static"))))
	router.Path("/").Methods("GET").HandlerFunc(web.Index)

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
	/*案例服务*/
	router.Path("/case/action/get-modules").Methods("GET").HandlerFunc(web.GetCaseModules)
	router.Path("/case/action/execute").Methods("POST").HandlerFunc(web.CaseExecute)
	router.Path("/case/action/save").Methods("POST").HandlerFunc(web.CaseSave)
	router.Path("/case/action/list").Methods("GET").HandlerFunc(web.GetCases)
	router.Path("/case/action/get-logs").Methods("GET").HandlerFunc(web.GetCaseLogs)
	router.Path("/case/{kid}").Methods("GET").HandlerFunc(web.GetCaseLog)

	//全部支持跨域
	router.Use(func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, req *http.Request) {
			w.Header().Set("Access-Control-Allow-Origin", "*")
			w.Header().Add("Access-Control-Allow-Headers", "content-type,workspace,application")
			next.ServeHTTP(w, req)
		})
	})
	router.Use(mux.CORSMethodMiddleware(router))

	return router
}

func main() {



	//地址端口
	addr := fmt.Sprintf("%s:%s", global.Conf.GetString("server.host"), global.Conf.GetString("server.port"))

	fmt.Printf("服务地址:%s", addr)
	//启动服务
	err := http.ListenAndServe(addr, initWebCtrl())
	if err != nil {
		log.Fatal("ListenAndServe", err)
	}
}
