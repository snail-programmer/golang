package main

import (
	"bytes"
	"context"
	"crypto/tls"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math/rand"
	"net/http"
	"time"
	_ "unsafe"
)

func runloop_send(client http.Client, request *http.Request) map[string]interface{} {
	response, err_c := client.Do(request.WithContext(context.TODO()))
	var rk map[string]interface{}

	if err_c != nil {
		fmt.Println("client err:", err_c)
		return rk
	}
	resByte, err_b := ioutil.ReadAll(response.Body)
	if err_b != nil {
		fmt.Println("err readAll:", err_b)
		return rk
	}
	err := json.Unmarshal(resByte, &rk)
	if err != nil {
		fmt.Println("to map error:", err)
		return rk
	}
	return rk
}
func init_param() {
	rand.Seed(time.Now().Unix())
	url := "https://39.105.252.205/api/mission/readChapter"
	post := "sign=f7288da01859640aab116075a6e4b786&time=1587972066430&bookId=4207956&chapterNum=7"
	postByte := []byte(post)
	postBuf := bytes.NewBuffer(postByte)
	request, err := http.NewRequest("POST", url, postBuf)
	if err != nil {
		fmt.Println("错误", err)
		return
	}
	request.Header.Set("COOKIE", "sessionid=d2d2a153bc7a49a09e55502162147ab4;vId=f711b1e7e8a143bdbf0d82fee2240027")
	request.Header.Set("X-Client", "sv=9;pm=YAL-AL50;ss=1080*2232;imei=866420047629145;imsi=460030915392511;mac=02:00:00:00:00:00;dID=00077993d88e75b5;version=5.1.35.26.53255;username=f711b1e7e8a143bdbf0d82fee2240027;signVersion=2;webVersion=new;oaid=null;pkv=1;ddid=DulNp2BgJWlkMdeQsvLMIFcIW0DKexrg+F1rJQMEB7EtS6sJ2xN79UfjXfrOv9cnmYKV3EsFSyhnlhgN6Xg+7BaQ;")
	request.Header.Set("Host", "api.ibreader.com")
	request.Header.Set("Accept-Encoding", "gzip")
	request.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	request.Header.Set("Content-Length", "85")

	tr := &http.Transport{TLSClientConfig: &tls.Config{InsecureSkipVerify: true}}
	client := http.Client{Transport: tr}
	var shown = false
	for true {
		var slp = rand.Intn(600)
		time.Sleep(time.Duration(slp) * time.Millisecond)
		rk := runloop_send(client, request)
		rs := rk["data"].(map[string]interface{})
		if rs["result"] != 100.0 {
			fmt.Println("failed=========:", rk["code"], rk["msg"])
			break
		} else {
			if !shown {
				fmt.Println("success========:", rk)
				shown = true
			}
		}
	}
}
func main() {
	init_param()
}
