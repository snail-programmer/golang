package main

import (
	"fmt"
	"net/http"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../Utils"
	"../safeHandler"
)

//获取作者总笔记数,访问量和被收藏量
func getAuthorFlow(userId string) (int, int) {
	note := DBModel.Article{AuthorId: userId}
	view_sum := DBCenter.DbgetSumWithModel(note, "view_num", "authorId")
	col_sum := DBCenter.DbgetSumWithModel(note, "collection", "authorId")
	return view_sum, col_sum
}

//我的笔记-分类
func categoryMyNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("未登录呢")
		w.Write(jsonStr)
		return
	}
	cate := DBModel.MycateNote{}
	cnt := DBCenter.DbgetCountWithSql("select count(cnt) from((select count(*) cnt from article where authorId=" + userId + " group by categoryId)a)")
	var store = make([]interface{}, cnt)
	var sql = "select sum(view_num) View_sum,sum(collection) Collect_sum, count(*) cnt,categoryName from article,category where authorId =" + userId + " and categoryId=category.id GROUP BY categoryId"
	DBCenter.DbgetModelWithSql(&cate, store, sql)
	rb := GenericPackJson(store)
	w.Write(rb)
}

//笔主详细信息页服务[user,note[],category[]]
func GetAuthorNote(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	var userId = ""
	authorId := r.URL.Query().Get("authorId")
	if authorId == "" {
		userId = safeHandler.GetCurrentUserId(w, r)
		if userId == "" {
			w.Write([]byte("未登录呢"))
			return
		}
		authorId = userId
	}
	var notelist []interface{}
	var catelist []interface{}
	user := DBModel.User{Id: authorId}
	note := DBModel.Article{AuthorId: authorId}
	//得到当前用户信息
	DBCenter.DbgetWithOneModel(&user)
	//得到用户笔记浏览量和收藏量
	v_s, c_s := getAuthorFlow(authorId)
	n_s := DBCenter.DbgetCount(&note)
	user.Note_sum = Utils.IntToString(n_s)
	user.View_sum = Utils.IntToString(v_s)
	user.Col_sum = Utils.IntToString(c_s)
	//以下字段对其他用户不可见
	user.Coin = ""
	user.Flow = ""
	user.Gratuity = ""
	user.Money = ""
	//得到当前用户所有笔记
	DBCenter.DbgetAllModel(&note, &notelist)
	//得到所有笔记的对应类型
	for _, e := range notelist {
		if e == nil {
			continue
		}
		cate := DBModel.Category{}
		cate.Id = e.(DBModel.Article).CategoryId
		DBCenter.DbgetWithOneModel(&cate)
		catelist = append(catelist, cate)
	}
	//[ user  notelist<>catelist ]
	var final = []interface{}{user, notelist, catelist}
	w.Write(GenericPackJson(final))
}
func getUserMoney(user *DBModel.User) {
	//得到用户笔记浏览量和收藏量
	v_s, c_s := getAuthorFlow(user.Id)
	user.View_sum = Utils.IntToString(v_s)
	user.Col_sum = Utils.IntToString(c_s)
	user.Flow = Utils.IntToString(v_s + c_s)
	user.Coin = Utils.AddNumString(user.Coin, user.Flow)
	//固定价值两元的100金币不可取出，防止重复销户注册套现
	var money = Utils.StringToFloat(Utils.AddNumString(user.Coin, user.Flow))/50 - 2.0
	if money < 0 {
		money = 0.00
	}
	user.Money = Utils.Float64ToString(money)
}
func getMyUser(w http.ResponseWriter, r *http.Request) {

	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		fmt.Println("未登录")
		w.Write(GenericPackJson("未登录呢"))
		return
	}
	user := DBModel.User{Id: userId}
	DBCenter.DbgetWithOneModel(&user)
	if user.Id == "" {
		w.Write(GenericPackJson("用户不存在!"))
		return
	}
	getUserMoney(&user)
	w.Write(GenericPackJson(user))
}

//打赏
func GiveGratuity(w http.ResponseWriter, r *http.Request) {
	r.ParseForm()
	userId := safeHandler.GetCurrentUserId(w, r)
	authorId := r.Form.Get("authorId")
	if userId == "" {
		PackMsgAndSend("400", "请先登录!", w)
		return
	}
	if authorId == "" {
		PackMsgAndSend("400", "要打赏的作者不存在!", w)
		return
	}
	coin := r.Form.Get("coin")
	coinNum := Utils.StringToInt(coin)
	if coinNum < 1 || coinNum > 10 {
		PackMsgAndSend("400", "失败，请打赏规定范围内的金币数量!", w)
		return
	}
	//查询当前用户信息
	user := DBModel.User{Id: userId}
	DBCenter.DbgetWithOneModel(&user)
	//金币是否充足
	surplus := int(Utils.StringToFloat(Utils.RdcNumString(user.Coin, coin)))
	if surplus < 0 {
		PackMsgAndSend("400", "失败，金币不足!", w)
		return
	}
	//查询被打赏用户信息
	luckUser := DBModel.User{Id: authorId}
	DBCenter.DbgetWithOneModel(&luckUser)
	luckUser.Gratuity = Utils.AddNumString(luckUser.Gratuity, coin)
	luckUser.Coin = Utils.AddNumString(luckUser.Coin, coin)
	luckUptUser := DBModel.User{Id: luckUser.Id, Coin: luckUser.Coin, Gratuity: luckUser.Gratuity}
	if luckUser.Id != "" {
		if !DBCenter.UpdateTable(luckUptUser) {
			PackMsgAndSend("400", "打赏失败，网络或系统故障!", w)
		} else {
			//打赏成功,减少自身金币
			user = DBModel.User{Id: userId, Coin: Utils.IntToString(surplus)}
			DBCenter.UpdateTable(user)
			PackMsgAndSend("100", "打赏成功!", w)
		}
	} else {
		PackMsgAndSend("400", "失败，该用户不存在!", w)
	}
}
