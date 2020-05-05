package main

import (
	"encoding/json"
	"fmt"
	"html/template"
	"net/http"

	"../Utils"
)

const (
	toRoot   = "../"
	viewPath = "../view/"
)

func JumpPage(w http.ResponseWriter, page string) {
	dispatch := map[string]string{"dispatch": page}
	jsonStr, _ := json.Marshal(&dispatch)
	w.Write(jsonStr)
}
func Dispatcher(w http.ResponseWriter, page string, info interface{}) {
	if page != "" && page[0:1] == "/" {
		page = page[1:len(page)]
	}
	path := toRoot + page
	res := Utils.IsfileExists(path)
	if !res {
		path = viewPath + page
		res = Utils.IsfileExists(path)
		if !res {
			if info != nil {
				w.Write(info.([]byte))
			}
			fmt.Println("404:", path, res)
		}
	}
	fmt.Println("dispatcher:", path)
	if res {

		t, err := template.ParseFiles(path)
		if w != nil {
			if err != nil {
				fmt.Println("err:", err)
			} else {
				t.Execute(w, info)
			}
		}
	}

}
