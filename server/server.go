package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"../plugins/net/websocket"
)

type handleFunc func(http.ResponseWriter, *http.Request)
type fileFunc http.Handler
type doHttpHandler struct {
	handler    map[string]handleFunc
	fileHandle map[string]fileFunc
}

/*
 分解请求路由与参数
 白名单:[.html,.htm,.tpl]
*/
var allowType = []string{
	"html", "htm", "tpl", "xml", "map",
	"css", "js", "txt", "map",
	"jpg", "jpeg", "bmp", "png", "gif", "svg",
}

//返回通用json数据处理格式
func GenericPackJson(any interface{}) []byte {
	mapStr := map[string]interface{}{"data": any}
	jsonStr, _ := json.Marshal(&mapStr)
	return jsonStr
}

//分解请求路由
func divideRequestUrl(url string) string {
	index := strings.Index(url, "?")
	if index < 0 {
		return url
	}
	res := url[0:index]
	return res
}

//是否是静态资源
func IsResourceServ(url string) bool {
	url = strings.ToLower(url)
	for _, e := range allowType {
		if strings.Contains(url, e) {
			return true
		}
	}
	return false
}
func IsSockServ(url string) bool {
	if len(url) >= 4 {
		pt := url[len(url)-4 : len(url)]
		if pt == "sock" {
			return true
		}
	}
	return false
}
func (tx *doHttpHandler) getServ(w http.ResponseWriter, r *http.Request) {
	srcUri := r.RequestURI
	routeUrl := divideRequestUrl(srcUri)
	//取链接最后的路由地址为rpc服务函数名
	lr := strings.LastIndex(routeUrl, "/")
	if lr != -1 {
		routeUrl = routeUrl[lr:len(routeUrl)]
	}

	//静态资源服务 || 动态资源服务
	if IsResourceServ(routeUrl) {
		handle := tx.fileHandle["RscServ"]
		handle.ServeHTTP(w, r)
	} else if IsSockServ(routeUrl) {
		handle := tx.fileHandle[routeUrl]
		handle.ServeHTTP(w, r)

	} else {
		run := tx.handler[routeUrl]
		if run != nil {
			run(w, r)
		} else {
			fmt.Println("get 无法为:" + routeUrl + "服务")
		}
	}
}
func (tx *doHttpHandler) postServ(w http.ResponseWriter, r *http.Request) {
	routeUrl := r.RequestURI
	//取链接最后的路由地址为rpc服务函数名
	lr := strings.LastIndex(routeUrl, "/")
	if lr != -1 {
		routeUrl = routeUrl[lr:len(routeUrl)]
	}
	//RPC，路由请求函数
	run := tx.handler[routeUrl]
	if run != nil {
		run(w, r)
	} else {
		fmt.Println("post 无法为:" + routeUrl + "服务")
	}
}

//网络根服务函数
func (tx *doHttpHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Access-Control-Allow-Credentials", "true")
	// w.Header().Set("Access-Control-Allow-Origin", "http://localhost")
	if r.Method == "GET" {
		tx.getServ(w, r)
	} else {
		tx.postServ(w, r)
	}
}
func RPC_SERVER_REGIST() {
	//注册RPC服务，路由请求函数
	myhandler := doHttpHandler{handler: make(map[string]handleFunc, 0),
		fileHandle: make(map[string]fileFunc, 0)}
	myhandler.handler["/login"] = login
	myhandler.handler["/regist"] = regist
	myhandler.handler["/logout"] = logout
	myhandler.handler["/deleteAccount"] = deleteAccount
	myhandler.handler["/perfectPersonInfo"] = perfectPersonInfo
	myhandler.handler["/getMyUser"] = getMyUser
	myhandler.handler["/categorylist"] = CategoryList
	myhandler.handler["/getAuthorNote"] = GetAuthorNote
	myhandler.handler["/categoryMyNote"] = categoryMyNote
	myhandler.handler["/getUserConfig"] = getUserConfig
	myhandler.handler["/updateUserConfig"] = updateUserConfig
	myhandler.handler["/queryPartNote"] = queryPartNote
	myhandler.handler["/delt_note"] = delt_note
	myhandler.handler["/view_note"] = view_note
	myhandler.handler["/colcel_note"] = colcel_note

	myhandler.handler["/modifyAccount"] = modifyAccount
	myhandler.handler["/modifyPassword"] = modifyPassword
	myhandler.handler["/getCollectState"] = getCollectState
	myhandler.handler["/getCollectNotes"] = getCollectNotes
	myhandler.handler["/GiveGratuity"] = GiveGratuity
	myhandler.handler["/push_comment"] = push_comment
	myhandler.handler["/get_comment"] = get_comment
	myhandler.handler["/del_comment"] = del_comment
	myhandler.handler["/get_notify_cnt"] = get_notify_cnt
	myhandler.handler["/get_message"] = get_message
	myhandler.handler["/remove_message"] = remove_message
	//上传
	myhandler.handler["/upload"] = upload
	myhandler.handler["/save_draft"] = save_draft
	myhandler.handler["/get_draft"] = get_draft
	myhandler.handler["/get_draftlist"] = get_draftlist
	myhandler.handler["/del_draft"] = del_draft
	myhandler.handler["/xhedit_saveNote"] = xhedit_saveNote
	myhandler.handler["/xhedit_uploadImg"] = xhedit_uploadImg
	//socket
	myhandler.fileHandle["/socknet.sock"] = websocket.Handler(socknet)
	//文件资源服务
	myhandler.fileHandle["RscServ"] = http.FileServer(http.Dir("../view"))
	err := http.ListenAndServe(":9000", &myhandler)
	if err != nil {
		fmt.Println(err)
	}
}

func main() {
	RPC_SERVER_REGIST()
}
