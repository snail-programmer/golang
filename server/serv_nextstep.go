package main

import (
	"net/http"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../safeHandler"
)

func perfectPersonInfo(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	//不允许登录跳转
	if !safeHandler.AllowPass(w, r) {
		http.Redirect(w, r, "login.html", http.StatusFound)
		return
	}
	//得到当前yi用户的session
	userId := safeHandler.GetCurrentUserId(w, r)
	user := DBModel.User{}
	//装载网络请求数据至model
	LoadModelWithPostData(&user, r.Form)
	//设置当前用户ID
	user.Id = userId
	//更新
	if DBCenter.UpdateTable(user) {
		http.Redirect(w, r, "notelist.html", http.StatusFound)
	}
}
