package main

import (
	"fmt"
	"math/rand"
	"time"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
)

/*
	用户个性化推荐算法
*/
//分类model[特征id,权重(from->to)]
type CateAnay struct {
	CateId string
	from   int
	to     int
}
type Recommend struct {
	CateAnays []interface{} //分类概率分析表 CateAnays[CateAnay{from->to代表权重}]
	max_len   int
	edu_char  string //教育经历首字母 屏蔽不同年级的笔记
	lim_wit   int    //分类权重上限 <= 分析数据的1/2
}

func (recd *Recommend) initEnviroment(userId string) {
	recd.CateAnays = make([]interface{}, 0)
	rand.Seed(time.Now().Unix())
	user := DBModel.User{Id: userId}
	DBCenter.DbgetWithOneModel(&user)
	//推荐笔记面向的教育层级， 除了不相同的年级，所有分类参与推荐计算
	vune := []rune(user.Education)
	if len(vune) > 1 {
		recd.edu_char = string(vune[0:1])
	}
}

/*
	根据cateId出现的次数计算从0开始的from和to,
	第二个分类的from在第一个to的基础上加1[from -> to  = 权重]
	one[from:0, to:2]->次数3
	two[from:3, to:5]->次数3
	...
*/
func (recd *Recommend) CalcCate() {
	//初始化分类模型
	ca := CateAnay{}
	var incre = 0
	if len(recd.CateAnays) > 1 {
		incre = recd.CateAnays[0].(CateAnay).to
		for i := 1; i < len(recd.CateAnays); i++ {
			if recd.CateAnays[i] != nil {
				ca = recd.CateAnays[i].(CateAnay)
				ca.from = incre + 1
				ca.to = ca.from + ca.to - 1
				recd.CateAnays[i] = ca
				incre = ca.to
			}
		}
	}
}

//添加分类到CateAnays中，按cateId结组并计算每个分类模型cateId出现的次数[to]
func (recd *Recommend) addCate(cateId string) {
	//初始化分类模型
	ca := CateAnay{}
	var find = false
	for i, obj := range recd.CateAnays {
		if obj != nil {
			ca = obj.(CateAnay)
			if cateId == ca.CateId {
				find = true
				//限制权重上限
				if ca.to < recd.lim_wit {
					ca.to++
					recd.CateAnays[i] = ca
				}
				break
			}
		}
	}
	if !find {
		ca.CateId = cateId
		ca.to = 1
		recd.CateAnays = append(recd.CateAnays, ca)
	}
}

//传递一组数据日志，筛选每条数据的分类特征
func (recd *Recommend) anaylysisLogs(viLogs []interface{}, userId string) {
	//初始化环境
	recd.initEnviroment(userId)
	recd.lim_wit = len(viLogs) / 2 //限制权重上限
	for _, obj := range viLogs {
		if obj != nil {
			var cid = obj.(DBModel.Visit_log).CategoryId
			//分离不同的categoryId
			recd.addCate(cid)
		}
	}

	cate := DBModel.Category{}
	var cates []interface{}
	DBCenter.DbgetAllModel(&cate, &cates)
	for _, obj := range cates {
		if obj != nil {
			//当前分类特征是否允许参加权重计算
			if !recd.allowCalcWeight(obj.(DBModel.Category).CategoryContain) {
				continue
			}
			var cid = obj.(DBModel.Category).Id
			//分离不同的categoryId
			recd.addCate(cid)
		}
	}
	//计算每个categoryId的权重
	recd.CalcCate()
	//max_len: from -> to De数量,概率分布范围[0,max(to)]
	ca := recd.CateAnays[len(recd.CateAnays)-1].(CateAnay)
	recd.max_len = ca.to + 1
	//
}

//得到一个推荐类型
func (recd *Recommend) getOneCateId() string {
	var rn = rand.Intn(recd.max_len)
	for i := 0; i < len(recd.CateAnays); i++ {
		obj := recd.CateAnays[i].(CateAnay)
		if rn >= obj.from && rn <= obj.to {
			return obj.CateId
		}
	}
	return ""
}

//不同年级的分类不参加权重计算
func (recd *Recommend) allowCalcWeight(edu string) bool {
	vune := []rune(edu)
	if len(vune) < 1 {
		return true
	}
	var ec = string(vune[0:1])
	if ec != recd.edu_char && (ec == "小" || ec == "初" || ec == "高" || ec == "大") {
		return false
	}
	return true
}

/*
	=============================推荐笔记委托结构体
*/
type ReNoteAri struct {
	has_noteid map[string]string
	cd         Recommend //用户行为特征分析结构体
	isInited   bool      //是否被初始化
}

func (rn *ReNoteAri) initEnviroment(viLogs []interface{}, userId string) {
	rn.has_noteid = make(map[string]string)
	rn.cd = Recommend{}
	rn.isInited = true
	//分析当前用户行为特征结果
	rn.cd.anaylysisLogs(viLogs, userId)
}
func (rn *ReNoteAri) avail_noteId(noteid string) bool {
	if noteid == "" {
		return false
	}
	if rn.has_noteid[noteid] == "" {
		rn.has_noteid[noteid] = "1"
		return true
	} else {
		return false
	}
}

//获取一个推荐笔记
func (rn *ReNoteAri) give_recommend_note(curNum int) interface{} {
	//获取一个推荐的特征id
	cid := rn.cd.getOneCateId()
	if cid != "" {
		//根据特征id 查询笔记
		article := DBModel.Article{}
		article.CategoryId = cid
		//test
		DBCenter.DbgetWithModel(&article, nil, 0, "createTime desc")
		if article.ArticleId == "" {
			fmt.Println("笔穷^-^")
			return nil
		}
		//笔记去重复
		if rn.avail_noteId(article.ArticleId) {
			return article
		} else {
			//如果重复了，偏移量+1 继续寻找，找到返回，否则结束循环，返回空对象
			for true {
				article = DBModel.Article{}
				article.CategoryId = cid
				curNum++
				fmt.Println("冲突增加偏移量:", curNum)
				DBCenter.DbgetWithModel(&article, nil, curNum, "createTime desc")
				if article.ArticleId != "" && rn.avail_noteId(article.ArticleId) {
					return article
				}
				if article.ArticleId == "" {
					break
				}

			}
		}
	}
	return nil
}
