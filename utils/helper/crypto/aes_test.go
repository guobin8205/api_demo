package crypto

import (
	"testing"
	//"encoding/base64"
	"encoding/json"
	"fmt"
	"net/http"
	"bytes"
	"io/ioutil"
	"log"
	"encoding/base64"
	"github.com/guobin8205/api_demo/utils/helper"
)

type Anchor struct {
	AnchorId   int    `json:"anchor_id"`
}

type Video struct {
	Id           int    `json:"-"`
	VideoId      int    `json:"video_id"`
	Name         string `json:"name"`
	AnchorId     int    `json:"anchor_id"`
	Hot          int    `json:"hot"`
	Time         int    `json:"time"`
	ViewTimes    int    `json:"view_times"`
	CollectTimes int    `json:"collect_times"`
}

func TestAes(t *testing.T){
	//account := "15858239848"
	roomId:= 253497
	//c,_ := json.Marshal(map[string]interface{}{"title":"推送测试001","content":"你关注的主播开播啦","payload":"1,54,1,5"})
	params:=fmt.Sprintf("act=push&tags=room_%d&title=主播上线拉&content=%s",roomId,"主播上线了")
	//params:= fmt.Sprintf("act=tag&type=add&account=%s&areaid=10&tags=%d",account,roomId)
	res := helper.HttpPost("http://10.225.10.111/m_send/tag.php",params)
	//log.Println(params)
	log.Println(res)


	//scomment := 5000*20/(5000+2)
	//s := float64(10000 + scomment + 3000*20)
	//log.Println("ss:",s)
	//now := int(time.Now().Unix())
	//a := math.Floor(float64((now - 1513846943)/3600))+0.25
	//log.Println("hours",a)
	//
	//calc_hot := fmt.Sprintf("%.8f",s/math.Pow(a,1.8))
	//log.Println(strconv.ParseFloat(calc_hot,64))
	//str := []byte(`{"uid":"988743680","limit":100,"anchor_id":1,"xx":"aaaaaaa","bb":3,"operate":0}`)
	//dst := NewAesCrypto().Encrypt(str)
	//fmt.Println(dst)
	//sdst,_ := NewAesCrypto().Decrypt(dst)
	//log.Println(string(sdst))

	//append(dst, []byte())
}

func TestAesCrypto_Encrypt(t *testing.T) {
	body, _ :=  json.Marshal(map[string]interface{}{"uuid":"bc61f307-258a-481c-9efa-2a3dc896b48e",
	"content":"真不錯啊！！","roomid":121748,"account":"ex01047588835qjb.dsp","anchor_id":58,"xx":"你好","bb":3,"limit":50,
	"operate":1,"time":"3.3","type":1})
	//body, _ :=  json.Marshal(map[string]interface{}{"nickname":"漂浮一生","uid":210038261794})
	//body, _ :=  json.Marshal(map[string]interface{}{"video_id":211,"anchor_id":1,"limit":10})
	//fmt.Println(string(body))
	dst := NewAesCrypto().Encrypt(body)
	//str := base64.StdEncoding.EncodeToString(dst)d
	//fmt.Println(str)
	//str = "gELQS4rNBc9zZNCxQgESesvKzzh46GgX06x6aq2oV07qX096Xn94mwIvS4G147s36P+U5PwawdAfInlg12fQi3OwMOEb+vIzA9SNQ51193Y="
	domain := "http://api.live.sanguosha.com"//正式服
	//domain := "http://121.199.6.172:5200"//外网测试
	//domain := "http://10.225.10.238:8080"//内网测试

	//api := "/anchor/test"
	//api := "/anchor/FansTop"
	//api := "/anchor/count"		//直播观看数统计
	//api := "/user/push"		//直播弹幕
	//api := "/video/push"		//视频留言弹幕
	//api := "/user/login"		//用户登录
	//api := "/video/get"		//视频信息
	//api := "/video/list"		//视频列表
	//api := "/video/DoCol"		//视频收藏
	api := "/anchor/list"		//主播列表
	//api := "/anchor/hotTop"		//主播列表
	//api := "/anchor/MyFollow"	//我的关注
	//api := "/anchor/DoFollow"	//主播关注操作
	//api := "/anchor/get" 		//主播信息
	fmt.Println(base64.StdEncoding.EncodeToString(dst))
	url := domain+api
	req, err := http.NewRequest("POST", url, bytes.NewBuffer(dst))
	req.Header.Set("Content-Type", "application/json")
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		panic(err)
	}
	defer resp.Body.Close()
	fmt.Println("response Status:", resp.Status)
	fmt.Println("response Headers:", resp.Header)
	body, _ = ioutil.ReadAll(resp.Body)
	//body, _ = ioutil.ReadAll(resp.Body)
	fmt.Println("response Body:", string(body))
}
