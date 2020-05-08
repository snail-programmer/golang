package main

import (
	"fmt"
	"net/http"
	"strings"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../Utils"
	"../safeHandler"
)

//获取所有分类categoryName:[categoryContain...]
func CategoryList(w http.ResponseWriter, r *http.Request) {
	cate_rs := make(map[string][]string)
	cate := DBModel.Category{CategoryName: "distinct"}
	catename := DBCenter.DbgetIdentify(&cate, 0)
	for _, v := range catename {
		cate = DBModel.Category{CategoryName: v, CategoryContain: "identify"}
		catecontain := DBCenter.DbgetIdentify(&cate, 0)
		cate_rs[v] = catecontain
	}
	w.Write(GenericPackJson(cate_rs))
}

/*
	查询一段数据，默认6，cur为以有的数据量，再次查询他应该从cur偏移的索引位置开始查询
	array[
		map[author,note]
		...
		]
*/
func queryPartNote(w http.ResponseWriter, r *http.Request) {
	curNum := r.URL.Query().Get("curNum")
	sortType := r.URL.Query().Get("sortType")
	keyword := r.URL.Query().Get("keyword")
	algo := false
	if strings.Index(keyword, "'") > -1 || strings.Index(keyword, "\\") > -1 {
		PackMsgAndSend("200", "请不要输入非法字符!", w)
		return
	}
	switch sortType {
	case "hotNote":
		sortType = "view_num desc"
		break
	case "recommend":
		keyword = ""
		algo = true
	default:
		sortType = "createTime desc"
		break
	}
	note := DBModel.Article{}
	store := make([]interface{}, 0)
	realStoreLen := 6
	//传统方式查询笔记
	if algo == false {
		store = make([]interface{}, realStoreLen)
		//关键字搜索
		if keyword != "" {
			resouce_search(keyword, curNum, sortType, store)
		} else {
			DBCenter.DbgetWithModel(&note, store, Utils.StringToInt(curNum), sortType)
		}
	} else {
		//推荐算法查询笔记-需登录
		userId := safeHandler.GetCurrentUserId(w, r)
		if userId == "" {
			//jsonStr := GenericPackJson("未登录呢")
			PackMsgAndSend("410", "使用该功能请先登录(ง •_•)ง ", w)
			return
		}
		//个性推荐查询笔记
		if curNum == "0" {
			get_recommend_init(userId) //初始化
		}
		for j := 0; j < realStoreLen; j++ {
			var rnote interface{} = nil
			for ri := 0; ri < 10; ri++ {
				rnote = get_recommend_note(userId, curNum)
				//笔记存在并且忽略自身的笔记
				if rnote != nil && rnote.(DBModel.Article).AuthorId != userId {
					break
				} else {
					rnote = nil
				}
			}
			if rnote != nil {
				note = rnote.(DBModel.Article)
				store = append(store, note)
				cate := DBModel.Category{Id: note.CategoryId}
				DBCenter.DbgetWithOneModel(&cate)
				fmt.Println("推荐:", note.ArticleId, cate)
			}
		}
	}
	//处理查询到的真正长度数据
	for i := 0; i < len(store); i++ {
		if store[i] == nil {
			// 计算出真实长度
			realStoreLen = i
			break
		}
		//每项map->model(Article)
		tmpNote := store[i].(DBModel.Article)
		//查询每个文章所属作者
		author := DBModel.User{Id: tmpNote.AuthorId}
		DBCenter.DbgetWithOneModel(&author)
		//查询每个文章所属分类
		category := DBModel.Category{Id: tmpNote.CategoryId}
		DBCenter.DbgetWithOneModel(&category)
		//传递笔记类型
		tmpNote.CategoryName = category.CategoryName
		tmpNote.CategoryContain = category.CategoryContain
		GeneArry := map[string]interface{}{"author": author, "note": tmpNote}
		//fmt.Println("tmpUser:", tmpUser)
		//更新到store
		store[i] = GeneArry
	}
	if realStoreLen == 0 && keyword != "" {
		PackMsgAndSend("400", "没有你想要的东西 ︿(￣︶￣) ", w)
		return
	}
	var realStore = make([]interface{}, realStoreLen)
	copy(realStore, store)
	//打包通用json格式回传
	w.Write(GenericPackJson(realStore))
}

//存储可复用的分析对象
var user_ari_reuse = make(map[string]ReNoteAri)

func get_recommend_init(userId string) {
	user_ari_reuse[userId] = ReNoteAri{}
}
func get_recommend_note(userId string, curNum string) interface{} {
	//查询访问日志记录,获取当前用户产生的<=10条记录
	viLog := DBModel.Visit_log{VisitId: userId}
	viLogs := make([]interface{}, 0)
	DBCenter.DbgetAllModel(&viLog, &viLogs)
	rn := user_ari_reuse[userId]

	//分析用户id为7 的个性特征
	if !rn.isInited {
		rn.initEnviroment(viLogs, userId)
		user_ari_reuse[userId] = rn
		fmt.Println("userId:" + userId + "===rn==init===")
	}

	//获取一条推荐笔记Utils.StringToInt(curNum)
	intr_note := rn.give_recommend_note(Utils.StringToInt(curNum))
	if intr_note == nil {
		return nil
	}
	note := intr_note.(DBModel.Article)
	return note
}
