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
		phone := strings.TrimSpace(r.Form.Get("phone"))
		password := strings.TrimSpace(r.Form.Get("password"))
		phone = template.HTMLEscapeString(phone)
		password = template.HTMLEscapeString(password)
		if phone == "" || password == "" {
			t, err := template.ParseFiles("../view/error.html")
			fmt.Println(err)
			e := map[string]string{"error": "失败:用户名或密码不能为空!", "jumpUrl": "login.html"}
			t.Execute(w, e)
			return
		}
		user := DBModel.User{}
		user.PhoneNumber = phone
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
		phone := strings.TrimSpace(r.Form.Get("phone"))
		verifycode := strings.TrimSpace(r.Form.Get("verifycode"))
		password := strings.TrimSpace(r.Form.Get("password"))
		retry_password := strings.TrimSpace(r.Form.Get("retry_password"))
		if len(phone) < 4 {
			JumpErrorPage(w, "失败:手机号格式不正确!", "regist.html")
			return
		}
		if verifycode != phone[len(phone)-4:len(phone)] {
			JumpErrorPage(w, "失败:验证码不正确!", "regist.html")
			return
		}
		if password != retry_password {
			JumpErrorPage(w, "失败:两次密码不一致!", "regist.html")
			return
		}
		userModel := DBModel.User{}
		userModel.PhoneNumber = phone
		userModel.Password = password
		userModel.HeadImg = "image/headimg_init.jpg"
		if !DBCenter.InsertTable(userModel) {
			t, err := template.ParseFiles("../view/error.html")
			fmt.Println(err)
			e := map[string]string{"error": "失败:此手机号已注册!", "jumpUrl": "regist.html"}
			t.Execute(w, e)
		} else {
			http.Redirect(w, r, "login.html?state=注册成功,请登录&phone="+phone, http.StatusFound)
		}
	}
}

//登出
func logout(w http.ResponseWriter, r *http.Request) {
	safeHandler.RemoveSession(w, r)
	http.Redirect(w, r, "notelist.html", http.StatusFound)
}

//永久注销
func deleteAccount(w http.ResponseWriter, r *http.Request) {
	//得到当前用户ID
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("请先登录该账户!")
		w.Write(jsonStr)
		return
	}
	user := DBModel.User{Id: userId}
	note := DBModel.Article{AuthorId: userId}
	draft := DBModel.Draft{AuthorId: userId}
	collect := DBModel.Collect_article{MyId: userId}
	log := DBModel.Visit_log{VisitId: userId}

	tpusr := user
	DBCenter.DbgetWithOneModel(&tpusr)
	getUserMoney(&tpusr)
	if tpusr.Money != "0.00" {
		PackMsgAndSend("400", "请先结算账户收益!", w)
		return
	}
	DBCenter.DeleteData(user)
	DBCenter.DeleteData(note)
	DBCenter.DeleteData(draft)
	DBCenter.DeleteData(collect)
	DBCenter.DeleteData(log)
	safeHandler.RemoveSession(w, r)
	PackMsgAndSend("100", "所有账户数据已被删除!", w)
}
