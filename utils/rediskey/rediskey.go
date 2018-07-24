package rediskey

import "fmt"

const (
	EXPIRE_USER       = 31536000 //用户信息失效时间
)
//用户个人信息
func UserKey(account string) (string, int) {
	return "user:"+account, EXPIRE_USER
}
func AnchorKey(anchor_id int)(string,int){
	return fmt.Sprintf("anchor:%d",anchor_id),EXPIRE_USER
}

func VideoKey(video_id int)(string,int){
	return fmt.Sprintf("video:%d",video_id),EXPIRE_USER
}
