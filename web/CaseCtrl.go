package web

import (
	"net/http"
	"otk-final/hinny/module/db"
	"strings"
)



type CaseModuleGroup struct {
	Module string   `json:"module"`
	Groups []string `json:"groups"`
}

func GetCaseModuleCroups(response http.ResponseWriter, request *http.Request) {
	db.Conn.Cols("").GroupBy("module,group").Find(db.CaseTemplate{})

	rows, err := db.Conn.Query("select module,group_concat(group_name) groups from hinny_case_template group by module")
	if err != nil {
		panic(err)
	}

	out := make([]CaseModuleGroup, 0)
	for _, row := range rows {
		item := &CaseModuleGroup{
			Module: string(row["module"]),
		}
		groups, ok := row["groups"]
		if !ok {
			item.Groups = []string{}
		} else {
			item.Groups = strings.Split(string(groups), ",")
		}
		out = append(out, *item)
	}
	view.JSON(response, 200, out)
}

func GetCaseGroups(response http.ResponseWriter, request *http.Request) {

}
