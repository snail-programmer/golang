package main

import (
	"fmt"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../Utils"
)

type NotesId struct {
	note1 []string
	note2 []string
	note3 []string
}
type CateInfo struct {
	cate1   []string
	cate2   []string
	cates   []string //交集
	notesId []string
}
type PersonInfo struct {
	uid     []string
	notesId []string
}

//封装类对象
type ReSrh struct {
	nsid       NotesId
	cateinfo   CateInfo
	personinfo PersonInfo
	sortType   string
}

func _append(sr []string, de []string) []string {
	find := false
	for i := 0; i < len(de); i++ {
		find = false
		for j := 0; j < len(sr); j++ {
			if de[i] == sr[j] {
				find = true
				continue
			}
		}
		if !find {
			sr = append(sr, de[i])
		}

	}
	return sr
}

//分类条件过滤
func (re *ReSrh) detachCategory(sr []string) {
	for i := 0; i < len(sr); i++ {
		if sr[i] == "" {
			continue
		}
		key := sr[i]
		//search all categoryName
		cate := DBModel.Category{Id: "identify", CategoryName: key}
		store := DBCenter.DbgetIdentify(&cate, 1)
		if len(store) > 0 {
			re.cateinfo.cate1 = append(re.cateinfo.cate1, store...)
			sr[i] = ""
		} else {
			//search categoryContain
			cate := DBModel.Category{Id: "identify", CategoryContain: key}
			store := DBCenter.DbgetIdentify(&cate, 1)
			if len(store) > 0 {
				re.cateinfo.cate2 = append(re.cateinfo.cate2, store...)
				sr[i] = ""
			}
		}

	}
	//两个非空列表取交集
	if len(re.cateinfo.cate1) == 0 {
		re.cateinfo.cates = append(re.cateinfo.cates, re.cateinfo.cate2...)
	} else if len(re.cateinfo.cate2) == 0 {
		re.cateinfo.cates = append(re.cateinfo.cates, re.cateinfo.cate1...)
	} else {
		for ca1 := 0; ca1 < len(re.cateinfo.cate1); ca1++ {
			for ca2 := 0; ca2 < len(re.cateinfo.cate2); ca2++ {
				if re.cateinfo.cate1[ca1] == re.cateinfo.cate2[ca2] {
					re.cateinfo.cates = append(re.cateinfo.cates, re.cateinfo.cate1[ca1])
				}
			}
		}
	}
	//根据cateId 查找noteid
	for _, v := range re.cateinfo.cates {
		note := DBModel.Article{ArticleId: "identify", CategoryId: v}
		//根据分类得到笔记id
		store := DBCenter.DbgetIdentify(&note, 0)
		if len(store) > 0 {
			re.nsid.note1 = append(re.nsid.note1, store...)
		}
	}
}

//附加信息 过滤
func (re *ReSrh) detachRemark(sr []string) {

	//获取 remark 弄忒多信息
	for i := 0; i < len(sr); i++ {
		if sr[i] == "" {
			continue
		}
		note := DBModel.Article{ArticleId: "identify", Remark: sr[i]}
		store := DBCenter.DbgetIdentify(&note, 1)
		if len(store) > 0 {
			re.nsid.note2 = _append(re.nsid.note2, store)
			sr[i] = ""
		}
	}
}

//昵称过滤
func (re *ReSrh) detachUser(sr []string) {
	for i := 0; i < len(sr); i++ {
		if sr[i] == "" {
			continue
		}
		sur := DBModel.User{Id: "identify", NickName: sr[i]}
		store := DBCenter.DbgetIdentify(&sur, 2)
		if len(store) > 0 {
			re.personinfo.uid = _append(re.personinfo.uid, store)
			sr[i] = ""
		}
	}
	//根据userId 查找noteId
	for i := 0; i < len(re.personinfo.uid); i++ {
		note := DBModel.Article{ArticleId: "identify", AuthorId: re.personinfo.uid[i]}
		store := DBCenter.DbgetIdentify(&note, 0)
		if len(store) > 0 {
			re.nsid.note3 = _append(re.nsid.note3, store)
		}
	}
}

func (re *ReSrh) analysisKeyWord(keyword string) []string {
	//关键字处理对象
	deal := dealKeyWord{}
	//关键字纠错
	queryword := deal.adjustError(keyword)
	//关键字每2个字符为一组分解数组
	arr := deal.cvtArrWithSepNum(queryword, 2)
	re.detachCategory(arr)
	re.detachRemark(arr)
	//查询用户
	keylen := len(queryword)
	//动态拆词,重查
	for i := 0; i < keylen; i++ {
		//关键字每keylen-i个字符为一组分解数组
		retryarr := deal.cvtArrWithSepNum(queryword, keylen-i)
		fmt.Println(retryarr)
		//查询用户
		re.detachUser(retryarr)
		if len(re.nsid.note3) > 0 && i >= (keylen-1)/2 {
			break
		}
	}
	fmt.Println("nsid.note1", re.nsid.note1)

	fmt.Println("nsid.note2", re.nsid.note2)

	fmt.Println("nsid.note3", re.nsid.note3)
	//取 nsid 交集
	//var capture []string
	capture := intersect(re.nsid.note1, re.nsid.note2, re.nsid.note3)
	if len(capture) == 0 {
		return capture
	}
	//对capture 排序
	retry_nid := "("
	for _, v := range capture {
		retry_nid += v + ","
	}
	retry_nid = retry_nid[0:len(retry_nid)-1] + ")"
	sql := "select ArticleId from article where ArticleId in" + retry_nid + " order by " + re.sortType
	fmt.Println(sql)
	retry_sort := make([]interface{}, len(capture))
	DBCenter.DbgetWithSql(sql, len(capture), 0, &DBModel.Article{}, retry_sort)
	capture = make([]string, 0)
	for _, obj := range retry_sort {
		capture = append(capture, obj.(DBModel.Article).ArticleId)
	}
	return capture
}
func resouce_search(keyword string, curNum string, sortType string, store []interface{}) {
	//初始化存储分析容器
	re := ReSrh{nsid: NotesId{}, cateinfo: CateInfo{}, personinfo: PersonInfo{}, sortType: sortType}
	//对关键字数组分析result[noteid]
	noteids := re.analysisKeyWord(keyword)
	fmt.Println("noteids:", noteids)

	crnum := Utils.StringToInt(curNum)
	for i := crnum; i < crnum+6; i++ {
		if i >= len(noteids) {
			break
		}
		note := DBModel.Article{ArticleId: noteids[i]}
		DBCenter.DbgetWithOneModel(&note)
		store[i-crnum] = note
	}
}

//取所有数组交集
func intersect(arrs ...[]string) []string {
	var store = map[string]int{}
	var ret []string
	var topcnt = 0
	for k, arr := range arrs {
		if len(arr) == 0 {
			continue
		}
		if k == 0 {
			for _, v := range arr {
				store[v] = 1
				topcnt = 1
			}
		} else {
			for _, v := range arr {
				store[v] += 1
				if store[v] > topcnt {
					topcnt = store[v]
				}
			}
		}
	}
	for k, v := range store {
		if v == topcnt {
			ret = append(ret, k)
		}
	}
	return ret
}
