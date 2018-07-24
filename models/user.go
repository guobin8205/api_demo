/**
 * Created by IDEA.
 * User: godsoul
 */
package models

import (
	"github.com/guobin8205/api_demo/utils/helper/custerror"
	"strconv"
	"github.com/garyburd/redigo/redis"
	"log"
	"time"
	"github.com/guobin8205/api_demo/models/db/rds"
	"github.com/guobin8205/api_demo/utils/rediskey"
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/guobin8205/api_demo/models/db/kafka"
	"net/rpc"
	"github.com/tidwall/gjson"
)

var official map[int]string

func init() {
	official = map[int]string{
		0: "士兵",
		1: "十夫长",
		2: "百夫长",
		3: "千夫长",
		4: "校尉",
		5: "先锋将军",
		6: "中军将军",
		7: "领军将军",
		8: "骠骑将军",
		9: "大将军",
	}
}

type User struct {
	Uid      int64  `gorm:"primary_key" redis:"uid"`
	Account  string `redis:"account"`
	Nickname string `gorm:"-" redis:"nickname"`
	Openid   int    `redis:"openid"`
	IsReg    int    `redis:"is_reg"`
	IP       string `gorm:"-" redis:"-"`
}

type Push struct {
	Openid   int
	Nickname string
	Content  string
	Roomid   int
	IP       string
}

func (User) TableName() string {
	return "user"
}
//直播弹幕
func (u *User) UserPush(roomid int, nickname, content, ip string) (interface{}, *custerror.Error) {
	if roomid == 0 || u.Uid == 0 || content == "" {
		return "参数错误", custerror.New(1000)
	}
	_redis := redisPool.Get()
	defer _redis.Close()
	openid, redis_err := redis.Int(_redis.HGet("user:"+u.Account, "openid"))
	officiallevel, _ := redis.Int(_redis.HGet("gameuser:"+u.Account, "officiallevel"))
	var officialname string
	if val, ok := official[officiallevel]; !ok {
		officialname = "士兵"
	} else {
		officialname = val
	}
	content = content + "(" + officialname + ")"
	if redis_err != nil || openid == 0 {
		return "用户未登录", custerror.New(4000)
	}
	//client,err := InitRPC()
	//if err != nil {
	//	return "Can't connect rpc server", custerror.New(9999)
	//}
	var reply int
	push := &Push{
		Openid:   openid,
		Nickname: nickname,
		Content:  content,
		Roomid:   roomid,
		IP:       ip,
	}
	if RPC == nil {
		InitRPC()
	}
	call_err := RPC.Call("RPC.Push", push, &reply)
	if call_err != nil {
		if call_err == rpc.ErrShutdown {
			InitRPC()
		}
		return call_err.Error(), custerror.New(9999)
	}

	if err := kafka.MSGpushKafka(map[string]interface{}{"app_id": 1, "area_id": 10, "area_name": "sj_live_show", "server_id": 30000,
		"event_info": map[string]interface{}{"id": 1001, "timestamp": time.Now().Unix() + 8*3600, "account": u.Account, "nk": nickname, "msg_type": 3000,
			"msg": content, "dest_nick": roomid}}); err != nil {
		log.Println(err)
	}
	room_info_key := fmt.Sprintf("room:%d", roomid)
	anchor_id, err := _redis.Get(room_info_key)
	if err == nil {
		_,week := time.Now().ISOWeek()
		week_times_key := fmt.Sprintf("anchor_count_%d:%s",week,anchor_id)
		_redis.HINCRBY(week_times_key, "push_times", 1)
	}
	log.Println("push content:", map[string]interface{}{"account": u.Account, "nk": nickname, "msg": content, "ip": ip})
	return reply, nil
}
//异步更新用户登录
func UserLoginList() {
	_redis := redisPool.Get()
	defer _redis.Close()
	for {
		end := _redis.BLPOP("user_login_list", 0)
		if len(end) > 1 {
			if end[1] != "" {
				uid := int(gjson.Get(end[1], "uid").Int())
				openid := int(gjson.Get(end[1], "openid").Int())
				if uid != 0{
					if err := Db.Model(&User{}).Where("uid = ?", uid).
						Updates(map[string]interface{}{"is_reg": 1, "openid": openid}).Error; err != nil {
						log.Println("登录更新失败，uid:", uid,",openid:",openid)
					}else{
						log.Println("登录更新成功，uid:", uid,",openid:",openid)
					}
				}
			}
		}
	}
}
//用户登录
func (u *User) Login(uuid string) (interface{}, *custerror.Error) {
	if uuid == "" || u.Account == "" {
		return "参数错误", custerror.New(1000)
	}
	//client,err := InitRPC()
	//if err != nil {
	//	return "Can't connect rpc server", custerror.New(9999)
	//}
	_redis := redisPool.Get()
	defer _redis.Close()
	if !CheckUser(u.Account, uuid) {
		return "用户验证不通过", custerror.New(4002)
	}
	nickname, _ := redis.String(_redis.HGet("gameuser:"+u.Account, "nick_name"))
	user := GetUserByAccount(u.Account)
	//第一次登录的账号注册
	if user == nil {
		user, err = newUser(u.Account)
		if err != nil {
			return "用户注册失败", custerror.New(4003)
		}
	}
	user.Nickname = nickname
	user.createCache(&_redis)
	if user.IsReg == 0 {
		var reply int
		//call_err := rpcx.Call(conf.RPCServer,"RPC.UserLogin", user, &reply)
		if RPC == nil {
			InitRPC()
		}
		user.IP = u.IP
		call_err := RPC.Call("RPC.UserLogin", user, &reply)
		if call_err != nil {
			if call_err == rpc.ErrShutdown {
				InitRPC()
			}
			return call_err.Error(), custerror.New(9999)
		}
		return map[string]int{"openid": 0}, nil
	} else {
		return map[string]int{"openid": user.Openid}, nil
	}
}
//创建redis缓存
func (u *User) createCache(r *rds.Rds) error {
	userKey, expire := rediskey.UserKey(u.Account)
	err := r.HMSet(userKey, u)
	if err == nil {
		uid_key := fmt.Sprintf("uid:%d", u.Uid)
		err = r.Set(uid_key, u.Account)
		r.Expire(userKey, expire)
	}
	return err
}
//创建新用户
func newUser(account string) (*User, error) {
	var user = &User{
		Account: account,
	}
	if err := Db.Create(user).Error; err != nil {
		return nil, err
	}
	return user, nil
}
//根据账号获取uid
func GetUid(account string) (int64, *custerror.Error) {
	_redis := redisPool.Get()
	defer _redis.Close()
	user_key := "user:" + account
	uid, err := redis.Int64(_redis.HGet(user_key, "uid"))
	if err != nil || uid == 0 {
		var u = &User{}
		err := Db.Where("account = ?", account).
			First(u).Error
		if err != nil {
			return 0, custerror.New(4004)
		} else {
			return u.Uid, nil
		}
	} else {
		return uid, nil
	}
}
//根据账号获取用户
func GetUserByAccount(account string) *User {
	var u = &User{}
	_redis := redisPool.Get()
	defer _redis.Close()
	user_key := "user:" + account
	if _redis.Exists(user_key) {
		user := _redis.HGetAll(user_key)
		u.Uid, _ = strconv.ParseInt(user["uid"], 10, 64)
		u.Account = user["account"]
		u.Nickname = user["nickname"]
		u.Openid, _ = strconv.Atoi(user["openid"])
		u.IsReg, _ = strconv.Atoi(user["is_reg"])
	} else {
		err := Db.Where("account = ?", account).
			First(u).Error
		if err != nil {
			log.Println(err)
			return nil
		}
	}
	return u
}
//检查用户登录状态
func CheckUser(account, uuid string) bool {
	_redis := redisPool.Get()
	defer _redis.Close()
	redis_uuid, err := redis.String(_redis.HGet("gameuser:"+account, "uuid"))
	if err == gorm.ErrRecordNotFound || err != nil || redis_uuid == "" || redis_uuid != uuid {
		log.Println("redis error:", err)
		return false
	} else {
		return true
	}
}

//func GetUserByUid(uid int64) *User{
//	var u = &User{}
//	_redis := redisPool.Get()
//	defer _redis.Close()
//	user_key := fmt.Sprintf("user:%d",uid)
//	if _redis.Exists(user_key){
//		user := _redis.HGetAll(user_key)
//		u.Uid,_ = strconv.ParseInt(user["uid"], 10, 64)
//		u.Account = user["account"]
//		u.Nickname = user["nickname"]
//		u.Openid,_ = strconv.Atoi(user["openid"])
//	}else{
//		err := Db.Where("uid = ?", uid).
//			First(u).Error
//		if err != nil {
//			log.Println(err)
//			return nil
//		}
//	}
//	return u
//}

//测试hmget
func GetRedis() {
	_redis := redisPool.Get()
	defer _redis.Close()
	user2, err := _redis.Con.Do("HMGET", "user:2", "uid", "openid")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("user2", user2)
	//s1 := []string{"user:2","uid","openid"}
	user, err := _redis.HMGet("user:2", "uid", "openid")
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println("user", user)
}
