package main

import (
	"html/template"
	"net/http"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../safeHandler"
)

func getUserConfig(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}
	user := DBModel.User{Id: userId}
	DBCenter.DbgetWithOneModel(&user)
	t, _ := template.ParseFiles("../view/template/setting.html")
	t.Execute(w, user)
}
func updateUserConfig(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	var jsonStr []byte
	if userId == "" {
		jsonStr = GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}
	r.ParseForm()
	user := DBModel.User{}
	LoadModelWithPostData(&user, r.Form)
	user.Id = userId
	if DBCenter.UpdateTable(user) {
		w.Write([]byte("success"))
		return
	}
	jsonStr = GenericPackJson("更新设置失败")
	w.Write(jsonStr)
}
func modifyAccount(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}
	newPhone := r.Form.Get("newPhone")
	verifyCode := r.Form.Get("verifyCode")
	if len(newPhone) < 4 {
		w.Write([]byte("手机号不正确"))
		return
	}
	if newPhone[len(newPhone)-4:len(newPhone)] != verifyCode {
		w.Write([]byte("验证码不正确"))
		return
	}
	user := DBModel.User{Id: userId, PhoneNumber: newPhone}
	if DBCenter.UpdateTable(user) {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("failed"))
	}
}
func modifyPassword(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录")
		w.Write(jsonStr)
		return
	}
	r.ParseForm()
	oldPassword := r.Form.Get("Password")
	newPassword := r.Form.Get("newPassword")
	if oldPassword == "" || newPassword == "" {
		w.Write([]byte("密码不能为空"))
		return
	}
	//源密码是否有效
	user := DBModel.User{Id: userId, Password: oldPassword}
	DBCenter.DbgetWithOneModel(&user)
	if user.Id == "" {
		w.Write([]byte("源密码错误"))
		return
	}
	//更新密码
	user = DBModel.User{Id: userId, Password: newPassword}
	if DBCenter.UpdateTable(user) {
		w.Write([]byte("success"))
	} else {
		w.Write([]byte("修改密码失败!"))
	}
}
