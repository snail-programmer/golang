package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"
)

type dealKeyWord struct {
	// keymap map[string][]string
}

var keymap = make(map[string][]string)

func (kw *dealKeyWord) initWordLibs() {
	if len(keymap) > 0 {
		fmt.Println("retry")
		return
	}
	file, err := os.Open("d:/bitch.txt")
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
