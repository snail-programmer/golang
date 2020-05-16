package main

import (
	"fmt"
	"net/http"
	"reflect"
	"time"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../Utils"
	"../plugins/net/websocket"
	"../safeHandler"
)

type SockServ struct {
	usrSock map[string]*websocket.Conn
}

var _sockserv = SockServ{usrSock: make(map[string]*websocket.Conn)}

func get_mynotify_cnt(userId string, readed string) int {
	sql := "select count(*) from article_comment where (readed=" + readed + " and  replyId in (select id from article_comment where myId =" + userId + ")) or (authorId =" + userId + " and readed=" + readed + ")  ORDER BY id desc"
	cnt := DBCenter.DbgetCountWithSql(sql)
	return cnt
}

//得到所有通知数量
func get_notify_cnt(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录该账户!")
		w.Write(jsonStr)
		return
	}
	//查询消息数量
	cnt := get_mynotify_cnt(userId, "0")
	comment_notify(cnt, userId)
}
func _get_message(userId string, readed string) []interface{} {
	if readed == "" {
		readed = "0"
	}
	cnt := get_mynotify_cnt(userId, readed)
	sql := "select * from article_comment where (readed=" + readed + " and  replyId in (select id from article_comment where myId =" + userId + ")) or (authorId =" + userId + " and readed=" + readed + ")  ORDER BY id desc"
	comment := DBModel.Article_comment{}
	store := make([]interface{}, cnt)
	for i := 0; i < cnt; i++ {
		DBCenter.DbgetWithSql(sql, 1, i, &comment, nil)
		if comment.ReplyId != "" {
			cmt := DBModel.Article_comment{Id: comment.ReplyId}
			DBCenter.DbgetWithOneModel(&cmt)
			comment.RefCmt = cmt.Comment
		}
		if comment.MyId != "" {
			user := DBModel.User{Id: comment.MyId}
			DBCenter.DbgetWithOneModel(&user)
			comment.HeadImg = user.HeadImg
			comment.NickName = user.NickName
		}
		store[i] = comment
	}
	return store
}

//评论通知
func comment_notify(cnt int, id string) {
	fmt.Println("消息通知:", id, cnt)
	if cnt > 0 {
		sock_send(_sockserv.usrSock[id], Utils.IntToString(cnt))
	}
}
func socknet(ws *websocket.Conn) {
	var reply string
	for {
		if err := websocket.Message.Receive(ws, &reply); err != nil {
		} else {
			fmt.Println("rec:", reply)
			_sockserv.usrSock[reply] = ws
			fmt.Printf("为" + reply + "建立")
			fmt.Println(ws)
		}
		time.Sleep(500)
	}
}
func sock_send(ws *websocket.Conn, data string) {
	fmt.Println("send-ws:", ws)
	if reflect.ValueOf(ws).IsNil() {
		return
	}
	if err := websocket.Message.Send(ws, data); err != nil {
		fmt.Println(err)
	}
}

//读取所有消息
func get_message(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	readed := r.URL.Query().Get("readed")
	if userId == "" {
		jsonStr := GenericPackJson("请先登录该账户!")
		w.Write(jsonStr)
		return
	}
	store := _get_message(userId, readed)
	w.Write(GenericPackJson(store))
}

//已读
func remove_message(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	id := r.URL.Query().Get("id")
	if userId == "" {
		jsonStr := GenericPackJson("请先登录该账户!")
		w.Write(jsonStr)
		return
	}
	if id == "" {
		return
	}
	if id == "all" {
		store := _get_message(userId, "0")
		for _, v := range store {
			cmt := DBModel.Article_comment{Id: v.(DBModel.Article_comment).Id, Readed: "1"}
			DBCenter.UpdateTable(cmt)
		}
		return
	}
	comment := DBModel.Article_comment{Id: id, Readed: "1"}
	DBCenter.UpdateTable(comment)
}
