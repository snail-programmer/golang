package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"time"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../Utils"
	"../safeHandler"
)

func push_comment(w http.ResponseWriter, r *http.Request) {
	var userId = safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		PackMsgAndSend("400", "请先登录", w)
		return
	}
	comment := DBModel.Article_comment{}
	bte, err := ioutil.ReadAll(r.Body)

	if err != nil {
		PackMsgAndSend("400", "系统错误", w)
		return
	}
	LoadModelWithByte(&comment, bte)
	if comment.Comment == "" {
		PackMsgAndSend("400", "评论不能为空", w)
		return
	}
	comment.MyId = userId
	comment.Cmt_time = time.Now().Format("2006-01-02 15:04:05")
	fmt.Println(comment)
	if DBCenter.InsertTable(comment) {
		PackMsgAndSend("100", "发布成功", w)
		//如果在自己文章底下评论不发布消息通知
		if userId == comment.AuthorId {
			fmt.Println("en?")
			return
		}
		if comment.ReplyId != "" {
			//通知被回复的用户
			cmt := DBModel.Article_comment{Id: comment.ReplyId}
			DBCenter.DbgetWithOneModel(&cmt)
			if cmt.MyId != "" {
				cnt := get_mynotify_cnt(cmt.MyId, "0")
				comment_notify(cnt, cmt.MyId)
			}
		} else {
			//通知文章作者
			if comment.AuthorId != "" {
				cnt := get_mynotify_cnt(comment.AuthorId, "0")
				comment_notify(cnt, comment.AuthorId)
			}
		}
	} else {
		PackMsgAndSend("400", "失败，未知错误", w)
	}
}
func get_comment(w http.ResponseWriter, r *http.Request) {
	unid := r.URL.Query().Get("id")
	aid := r.URL.Query().Get("ArticleId")
	scnt := r.URL.Query().Get("cur_comment_cnt")
	cnt := Utils.StringToInt(scnt)
	comment := DBModel.Article_comment{}
	if unid != "" {
		comment.Id = unid
	} else {
		comment.ArticleId = aid
	}
	comment_cnt := DBCenter.DbgetCount(&comment)
	if comment_cnt > 10 {
		comment_cnt = 10
	}
	var comments = make([]interface{}, comment_cnt)
	for i := 0; i < comment_cnt; i++ {
		if unid != "" {
			comment = DBModel.Article_comment{Id: unid}
		} else {
			comment = DBModel.Article_comment{ArticleId: aid}
		}
		DBCenter.DbgetWithModel(&comment, nil, cnt+i, "id desc")
		if comment.MyId == "" {
			break
		}
		user := DBModel.User{Id: comment.MyId}
		DBCenter.DbgetWithOneModel(&user)
		comment.NickName = user.NickName
		comment.HeadImg = user.HeadImg
		comments[i] = comment
	}
	if len(comments) > 0 {
		w.Write(GenericPackJson(comments))
	} else {
		if unid != "" {
			PackMsgAndSend("400", "此评论被作者删除", w)
		} else {
			PackMsgAndSend("tip", "没有更多数据了U•ェ•*U", w)
		}
	}
}
func del_comment(w http.ResponseWriter, r *http.Request) {
	var userId = safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		PackMsgAndSend("400", "请先登录", w)
		return
	}
	var cmt_time = r.URL.Query().Get("cmt_time")
	comment := DBModel.Article_comment{MyId: userId, Cmt_time: cmt_time}
	cmt := comment
	DBCenter.DbgetWithOneModel(&cmt)
	replyId := cmt.ReplyId
	if DBCenter.DeleteData(comment) {
		PackMsgAndSend("400", "删除成功", w)
		//如果是删除回复别人的评论，则通知对方更新消息
		if replyId != "" {
			//查找被回复用户的id
			comment := DBModel.Article_comment{Id: replyId}
			DBCenter.DbgetWithOneModel(&comment)
			cnt := get_mynotify_cnt(comment.MyId, "0")
			comment_notify(cnt, comment.MyId)
		}
		if cmt.AuthorId != "" {
			fmt.Println("删除评论 - 通知作者")
			cnt := get_mynotify_cnt(cmt.AuthorId, "0")
			comment_notify(cnt, cmt.AuthorId)
		}
	} else {
		PackMsgAndSend("400", "删除失败", w)
	}
}
