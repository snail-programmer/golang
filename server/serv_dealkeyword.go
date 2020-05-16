package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type dealKeyWord struct {
}

var keymap = make(map[string][]string)

func (kw *dealKeyWord) initWordLibs() {
	if len(keymap) > 0 {
		fmt.Println("retry")
		return
	}
	file, err := os.Open("../view/config/vague_match.ini")
	if err != nil {
		fmt.Println("read file fail", err)
		return
	}
	defer file.Close()
	bytes, err := ioutil.ReadAll(file)
	if err != nil {
		fmt.Println("read to fd fail", err)
		return
	}
	//arr[record1,...]
	arr := strings.Split(string(bytes), "\r\n")
	for _, av := range arr {
		av = strings.ReplaceAll(av, " ", "")
		//kv[0]:kv[1]
		kv := strings.Split(av, ":")
		if len(kv) < 2 {
			continue
		}
		v := strings.Split(kv[1], ",")
		keymap[kv[0]] = v
	}
}

//纠错
func (kw *dealKeyWord) adjustError(word string) string {
	kw.initWordLibs()
	word = strings.ReplaceAll(word, " ", "")
	for k, arr := range keymap {
		for _, v := range arr {
			if strings.Contains(word, v) {
				word = strings.ReplaceAll(word, v, k)
			}
		}
	}
	return word
}

//str分解为arr
func (kw *dealKeyWord) cvtArrWithSepNum(queryword string, sn int) []string {
	utfword := []rune(queryword)
	arr := make([]string, 0)
	for i := 0; i < len(utfword)-1; i++ {
		var tmp string
		for n := i; n < i+sn; n++ {
			//剩下的字符串不足以组成sn个字符,退出
			if i+sn > len(utfword) {
				break
			} else {
				tmp += string(utfword[n])
			}
		}
		if len(tmp) > 0 {
			arr = append(arr, tmp)
		}
	}
	return arr
}

//动态丢词完成通知回调
type notify func(arr []string) int

//动态丢词
func (kw *dealKeyWord) dynamicDiscardWord(queryword string, callback notify) {
	//查询用户
	keylen := len(queryword)
	//动态拆词,重查
	for i := 0; i < keylen; i++ {
		//关键字每keylen-i个字符为一组分解数组，为1返回
		if keylen-i <= 1 {
			return
		}
		retryarr := kw.cvtArrWithSepNum(queryword, keylen-i)
		//调用回调,返回搜索成功的数量
		callback(retryarr)
		if keylen-i <= 2 {
			break
		}
	}
}
