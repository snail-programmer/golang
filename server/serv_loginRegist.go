package main

import (
	"fmt"
	"html/template"
	"net/http"
	"strings"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../safeHandler"
)

func login(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "POST" {
		username := strings.TrimSpace(r.Form.Get("username"))
		password := strings.TrimSpace(r.Form.Get("password"))
		username = template.HTMLEscapeString(username)
		password = template.HTMLEscapeString(password)
		if username == "" || password == "" {
			t, err := template.ParseFiles("../view/error.html")
			fmt.Println(err)
			e := map[string]string{"error": "失败:用户名或密码不能为空!", "jumpUrl": "login.html"}
			t.Execute(w, e)
			return
		}
		user := DBModel.User{}
		user.PhoneNumber = username
		user.Password = password
		DBCenter.DbgetWithOneModel(&user)
		fmt.Println("userModel:", user)
		if user.Id != "" {
			session := safeHandler.SGCookieSession(w, r, "set")
			session.UserId = user.Id
			//更新session,存入用户id
			safeHandler.UpdateSession(session)
			var jumpUrl = "notelist.html"
			//完善个人信息
			if user.Education == "" || user.NickName == "" {
				//jumpUrl = "nextstep.html"
				t, _ := template.ParseFiles("../view/nextStep.html")
				t.Execute(w, user)
			} else {
				http.Redirect(w, r, jumpUrl, http.StatusFound)
			}

		} else {
			t, err := template.ParseFiles("../view/error.html")
			fmt.Println(err)
			e := map[string]string{"error": "失败:用户名或密码不正确!", "jumpUrl": "login.html"}
			t.Execute(w, e)
			return
		}
	}
}
func regist(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	if r.Method == "POST" {
		username := strings.TrimSpace(r.Form.Get("username"))
		password := strings.TrimSpace(r.Form.Get("password"))
		fmt.Println(username, password)
		if username == "" || password == "" {
			t, err := template.ParseFiles("../view/error.html")
			fmt.Println(err)
			e := map[string]string{"error": "失败:用户名或密码不正确!", "jumpUrl": "login.html"}
			t.Execute(w, e)
			return
		}
		userModel := DBModel.User{}
		userModel.PhoneNumber = username
		userModel.Password = password
		userModel.HeadImg = "image/headimg_init.jpg"
		if !DBCenter.InsertTable(userModel) {
			t, err := template.ParseFiles("../view/error.html")
			fmt.Println(err)
			e := map[string]string{"error": "失败:此手机号已注册!", "jumpUrl": "login.html"}
			t.Execute(w, e)
		} else {
			http.Redirect(w, r, "login.html", http.StatusFound)
		}
	}
}
func logout(w http.ResponseWriter, r *http.Request) {
	safeHandler.RemoveSession(w, r)
	http.Redirect(w, r, "notelist.html", http.StatusFound)
}
