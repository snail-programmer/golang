package main

import (
	"net/http"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../Utils"
	"../safeHandler"
)

//收藏/取消收藏笔记
func colcel_note(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	noteid := r.Form.Get("noteId")
	ope := r.Form.Get("ope")
	ret := map[string]string{}
	if userId == "" {
		ret := map[string]string{"state": "400", "msg": "请先登录!"}
		w.Write(GenericPackJson(ret))
		return
	}
	if noteid == "" {
		ret := map[string]string{"state": "400", "msg": "收藏的笔记无效!"}
		w.Write(GenericPackJson(ret))
		return
	}
	//原始笔记是否存在
	srcnote := DBModel.Article{ArticleId: noteid}
	cnt := DBCenter.DbgetCount(&srcnote)
	if cnt == 0 {
		PackMsgAndSend("400", "该笔记已被删除或不可用!", w)
		return
	}
	colnote := DBModel.Collect_article{MyId: userId, ArticleId: noteid}
	isExitcolNote := colnote
	//查询是否已经收藏
	cnt = DBCenter.DbgetCount(&isExitcolNote)
	var suc = "fail"
	var msg = "失败，网络故障!"
	//该文章被收藏的次数
	note := DBModel.Article{ArticleId: noteid}
	DBCenter.DbgetWithOneModel(&note)
	colnum := note.Collection
	//添加/取消收藏
	if ope == "delt_collect" {
		if DBCenter.DeleteData(colnote) {
			//更新该笔记收藏量
			note := DBModel.Article{ArticleId: colnote.ArticleId}
			if colnum != "0" {
				note.Collection = Utils.RdcNumString(colnum, "1")
				DBCenter.UpdateTable(note)
			}
			suc = "success"
		} else {
			msg = "删除失败，该笔记不存在!"
		}
	} else {
		if cnt != 0 {
			msg = "此笔记以被收藏"
		} else {
			if DBCenter.InsertTable(colnote) {
				//更新该笔记收藏量
				note := DBModel.Article{ArticleId: colnote.ArticleId}
				note.Collection = Utils.AddNumString(colnum, "1")
				DBCenter.UpdateTable(note)
				suc = "success"
			} else {
				msg = "添加收藏失败，系统异常"
			}
		}

	}
	ret = map[string]string{"state": suc, "msg": msg}
	w.Write(GenericPackJson(ret))
}

//获取收藏状态
func getCollectState(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	noteid := r.Form.Get("noteId")
	ret := map[string]string{}
	if userId == "" {
		ret = map[string]string{"state": "400", "msg": "请先登录!"}
		w.Write(GenericPackJson(ret))
		return
	}
	if noteid == "" {
		ret = map[string]string{"state": "400", "msg": "该笔记无效!"}
		w.Write(GenericPackJson(ret))
		return
	}
	col_note := DBModel.Collect_article{MyId: userId, ArticleId: noteid}
	DBCenter.DbgetWithOneModel(&col_note)
	if col_note.Id != "" {
		ret = map[string]string{"state": "true"}
	} else {
		ret = map[string]string{"state": "false"}
	}
	w.Write(GenericPackJson(ret))
}
func getCollectNotes(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var userId = safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		w.Write([]byte("未登录呢"))
		return
	}
	catename := r.URL.Query().Get("catename")
	var collects []interface{}
	collect := DBModel.Collect_article{MyId: userId}
	DBCenter.DbgetAllModel(&collect, &collects)
	var notelist []interface{}
	var catelist []interface{}
	for i := 0; i < len(collects); i++ {
		if nil == collects[i] {
			continue
		}
		cole := collects[i].(DBModel.Collect_article)
		//得到对应笔记
		note := DBModel.Article{ArticleId: cole.ArticleId}
		DBCenter.DbgetWithOneModel(&note)
		if note.ArticleId != "" {
			//得到当前笔记的作者
			user := DBModel.User{Id: note.AuthorId}
			DBCenter.DbgetWithOneModel(&user)
			note.HeadImg = user.HeadImg //传入头像
			note.NickName = user.NickName
			//得到当前笔记的分类
			cate := DBModel.Category{Id: note.CategoryId}
			DBCenter.DbgetWithOneModel(&cate)
			if cate.CategoryName == catename || catename == "全部" {
				catelist = append(catelist, cate)
				notelist = append(notelist, note)
			}
		}
	}
	myUser := DBModel.User{Id: userId}
	cnt := len(notelist)
	// DBCenter.DbgetCount(&DBModel.Collect_article{MyId: userId})
	myUser.Note_sum = Utils.IntToString(cnt) //我的收藏笔记数量
	//[ user  notelist<>catelist ]
	var final = []interface{}{myUser, notelist, catelist}
	w.Write(GenericPackJson(final))
}
