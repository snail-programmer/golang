package main

import (
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	DBCenter "../DataBaseCenter"
	DBModel "../DataBaseCenter/DataBaseModel"
	"../safeHandler"
)

//application/octet-stream 接收二进制流
//1.保存图片路径; 2.保存文本到数据库
type FileInfo struct {
	Attachment string
	Name       string
	Filename   string
	ExtName    string
}

//笔记内图片上传
func xhedit_uploadImg(w http.ResponseWriter, r *http.Request) {
	cd := r.Header.Get("Content-Disposition")
	fi := FileInfo{}
	LoadModelWithByte(&fi, []byte(cd))
	//去引号
	if strings.Index(fi.Filename, "\"") > -1 {
		fi.Filename = strings.ReplaceAll(fi.Filename, "\"", "")
	}
	//取格式名
	rx := strings.LastIndex(fi.Filename, ".")
	if rx != -1 {
		fi.ExtName = fi.Filename[rx:len(fi.Filename)]
	}
	fmt.Println(fi)
	data := r.Body
	nano := time.Now().UnixNano()
	filename := strconv.FormatInt(nano, 10)
	filename += fi.ExtName
	filepath := "../view/image/upload/" + filename
	fmt.Println("filepath:", filepath)
	file, err := os.Create(filepath)
	if err != nil {
		fmt.Println(err)
	}
	_, error := io.Copy(file, data)
	if error != nil {
		fmt.Println(error)
	}

	defer file.Close()

	//{err:"错误信息提示",msg:"图片路径"}
	var msg = "失败"
	if err == nil {
		msg = "http://localhost:9000/image/upload/" + filename
	}
	var res = map[string]string{"err": "", "msg": msg}
	jsonStr, _ := json.Marshal(&res)
	w.Write(jsonStr)
}

//用户头像上传
func upload(w http.ResponseWriter, r *http.Request) {
	userId := safeHandler.GetCurrentUserId(w, r)
	if userId == "" {
		jsonStr := GenericPackJson("未登录呢")
		w.Write(jsonStr)
		return
	}
	var safexg = []string{"jpg", "jpeg", "png", "bmp"}
	r.ParseMultipartForm(1024 * 100)
	file, handle, err := r.FormFile("uploadimg")
	if err != nil {
		fmt.Println(err)
	}
	var filename = handle.Filename
	//有些浏览器会发送包含的路径，截取最后的文件名
	var lastix = strings.LastIndex(filename, "\\")
	if lastix > -1 {
		filename = filename[lastix+1 : len(filename)]
	}
	var check = filename[len(filename)-4 : len(filename)]
	check = strings.ToLower(check)
	var allow = false
	for _, e := range safexg {
		if strings.Contains(check, e) {
			allow = true
			break
		}
	}
	if !allow {
		GenericPackJson("error")
		fmt.Println("what")
		return
	}
	filename = userId + "_" + filename
	defer file.Close()
	imgurl := "image/userimg/" + filename
	f, err := os.OpenFile("../view/"+imgurl, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		fmt.Println(imgurl)
		fmt.Println("upload-108:", err)
		return
	}
	defer f.Close()
	io.Copy(f, file)

	//更新头像至数据库,删除之前的头像文件
	user := DBModel.User{Id: userId}
	tmp := user
	DBCenter.DbgetWithOneModel(&tmp)
	user.HeadImg = imgurl
	if DBCenter.UpdateTable(user) {
		fmt.Println("update headimg success")
		//删除之前的文件
		err := os.Remove("../view/" + tmp.HeadImg)
		if err != nil {
			fmt.Println("删除原始头像失败")
		}
		w.Write([]byte("success"))
	} else {
		fmt.Println("update headimg failed")
		w.Write([]byte("failed"))
	}

}
