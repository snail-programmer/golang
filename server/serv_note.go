package main

import (
	"io/ioutil"
	"net/http"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../Utils"
	"../safeHandler"
)

//保存笔记
func xhedit_saveNote(w http.ResponseWriter, r *http.Request) {
	//得到当前用户ID
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}

	article := DBModel.Article{}
	rbyte, _ := ioutil.ReadAll(r.Body)
	LoadModelWithByte(&article, rbyte)

	cate := DBModel.Category{CategoryName: article.CategoryName, CategoryContain: article.CategoryContain}
	DBCenter.DbgetWithOneModel(&cate)

	article.AuthorId = userId
	article.CategoryId = cate.Id
	article.CategoryName = ""
	article.CategoryContain = ""
	res := false
	if article.ArticleId == "" {
		res = DBCenter.InsertTable(article)
	} else {
		//重新发布笔记-时间不变
		article.CreateTime = ""
		res = DBCenter.UpdateTable(article)
	}

	retStr := "fail"
	if res {
		retStr = "success"
	}
	w.Write(GenericPackJson(retStr))

}

//删除
func delt_note(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	noteid := r.Form.Get("noteId")
	if userId == "" {
		jsonStr := GenericPackJson("未登录呢")
		w.Write(jsonStr)
		return
	}
	if noteid == "" {
		jsonStr := GenericPackJson("无法删除")
		w.Write(jsonStr)
		return
	}
	//删除笔记，同时删除所有对该笔记的收藏，评论
	article := DBModel.Article{ArticleId: noteid, AuthorId: userId}
	collect := DBModel.Collect_article{ArticleId: noteid}
	if DBCenter.DeleteData(article) {
		DBCenter.DeleteData(collect)
		jsonStr := GenericPackJson("success")
		w.Write(jsonStr)
	} else {
		jsonStr := GenericPackJson("请求被拒绝!")
		w.Write(jsonStr)
	}

}

//浏览一条笔记
func view_note(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	noteid := r.Form.Get("noteId")
	if noteid == "" {
		jsonStr := GenericPackJson("未查询到数据!")
		w.Write(jsonStr)
		return
	}
	//查询一条笔记
	article := DBModel.Article{ArticleId: noteid}
	DBCenter.DbgetWithOneModel(&article)
	w.Write(GenericPackJson(article))

	if article.CategoryId != "" && userId != "" && userId != article.AuthorId {
		//更新访问量
		v_n := Utils.AddNumString(article.View_num, "1")
		udnote := DBModel.Article{ArticleId: article.ArticleId, View_num: v_n}
		DBCenter.UpdateTable(udnote)
		//访问日志记录
		record_visit_log(userId, article.CategoryId)
	}

}
