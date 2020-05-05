package main

import (
	"fmt"
	"io/ioutil"
	"net/http"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../safeHandler"
)

//保存草稿
func save_draft(w http.ResponseWriter, r *http.Request) {

	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}
	draft := DBModel.Draft{}
	// LoadModelWithPostData(&draft, r.Form)
	rbyte, _ := ioutil.ReadAll(r.Body)
	LoadModelWithByte(&draft, rbyte)
	draft.AuthorId = userId

	//查询创建该草稿的时间线是否存在笔记
	tmpDraft := draft
	tmpDraft.DraftNote = ""
	DBCenter.DbgetWithOneModel(&tmpDraft)

	//没有查询到笔记 -> 插入
	var res = false
	if tmpDraft.DraftNote == "" {
		fmt.Println("上传草稿-userid:", userId)
		res = DBCenter.InsertTable(draft)
	} else {
		fmt.Println("更新草稿-userid:", userId)
		//忽略创建时间的更新
		draft.CreateTime = ""
		//查询到的tmpDraft的id 给 draft
		draft.Id = tmpDraft.Id
		res = DBCenter.UpdateTable(draft)
	}
	retStr := "fail"
	if res {
		retStr = "success"
	}
	w.Write(GenericPackJson(retStr))
}
func get_draft(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	draft := DBModel.Draft{}
	LoadModelWithPostData(&draft, r.Form)
	DBCenter.DbgetWithOneModel(&draft)
	dByte := GenericPackJson(draft)
	w.Write(dByte)
}
func del_draft(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}
	draft := DBModel.Draft{}
	draft.AuthorId = userId
	draft.CreateTime = r.Form.Get("createTime")
	if draft.CreateTime == "" {
		w.Write([]byte("request forbit"))
	}
	if DBCenter.DeleteData(draft) {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("failed"))
	}
}
func get_draftlist(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}
	draft := DBModel.Draft{}
	draft.AuthorId = userId
	var allData []interface{}
	DBCenter.DbgetAllModel(&draft, &allData)
	vByte := GenericPackJson(allData)
	w.Write(vByte)
}
