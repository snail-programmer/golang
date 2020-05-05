package safeHandler

import (
	"fmt"
	"net/http"
	"strconv"
	"time"
)

var sessionMap map[string]Session

type Session struct {
	sessionid string
	oldtime   int
	UserId    string
}

func init() {
	sessionMap = make(map[string]Session, 0)
	go sessionManage()
}
func sessionManage() {
	for {
		for k, v := range sessionMap {
			v.oldtime = v.oldtime - 1
			if v.oldtime == 0 {
				delete(sessionMap, k)
				fmt.Println("sessionId:" + k + "已过期")
			} else {
				sessionMap[k] = v
			}
		}
		time.Sleep(time.Second)
	}
}

//设置或从cookies得到sessionId
func SGCookieSession(w http.ResponseWriter, r *http.Request, method string) Session {
	if method == "set" {
		nano := time.Now().UnixNano()
		str := strconv.FormatInt(nano, 10)
		cookie := http.Cookie{Name: "gosessionid", Value: str, Path: "/", HttpOnly: true, MaxAge: 5000}
		http.SetCookie(w, &cookie)
		session := Session{sessionid: str, oldtime: 5000}
		sessionMap[str] = session
		return sessionMap[str]
	} else {
		cookie, err := r.Cookie("gosessionid")
		if err != nil {
			fmt.Println("cookie-error:", err)
			return Session{}
		}
		//fmt.Println("cookie:", cookie.Value, sessionMap[cookie.Value])
		return sessionMap[cookie.Value]
	}
}

//sessionid 得到session数据
func GetCurrentUserId(w http.ResponseWriter, r *http.Request) string {
	session := SGCookieSession(w, r, "get")
	return session.UserId
}
func UpdateSession(session Session) {
	sessionMap[session.sessionid] = session
	fmt.Println("UpdateSession:", sessionMap[session.sessionid])
}
func AllowPass(w http.ResponseWriter, r *http.Request) bool {
	session := SGCookieSession(w, r, "get")
	if session.oldtime > 0 {
		return true
	}
	return false
}
func RemoveSession(w http.ResponseWriter, r *http.Request) {
	session := SGCookieSession(w, r, "get")
	delete(sessionMap, session.sessionid)
}
