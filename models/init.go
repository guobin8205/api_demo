/**
 * Created by IDEA.
 * User: godsoul
 */
package models

import (
	"github.com/jinzhu/gorm"
	"github.com/guobin8205/api_demo/models/db/rds"
	"github.com/guobin8205/api_demo/utils/config"
	"github.com/guobin8205/api_demo/models/db/mysql"
	"log"
	"net/rpc"
	crpc "github.com/guobin8205/api_demo/utils/rpc"
)

var (
	redisPool *rds.RdsPool
	Db        *gorm.DB
	RPC       *rpc.Client
	err       error
)

type NoArg struct {
}

type NoReply struct {
}

func init() {
	//初始化kafka
	//kafka_addr := strings.Split(conf.String("kafka.addrs"), ",")
	//err := kafka.InitKafka(kafka_addr)
	//if err != nil {
	//	panic(err)
	//}
	//初始化redis连接池
	redisPool = rds.NewRdsPool(
		conf.String("redis.host"),
		conf.String("redis.port"),
		conf.String("redis.password"),
		conf.Int("redis.maxidle"),
		conf.Int("redis.maxactive"),
	)
	_redis := redisPool.Get()
	defer _redis.Close()
	//初始化mysql连接池
	Db = mysql.NewMysqlCon(
		conf.String("mysql.user"),
		conf.String("mysql.pass"),
		conf.String("mysql.host"),
		conf.String("mysql.port"),
		conf.String("mysql.db"),
		conf.String("mysql.charset"),
		conf.Int("mysql.maxidle"),
		conf.Int("mysql.maxactive"),
	)
	//登录异步处理goroutine
	go UserLoginList() //登录异步处理
	//初始化RPC
	//err = InitRPC()
	//if err != nil {
	//	panic("rpc connect error!")
	//}
}

func InitRPC() error {
	RPC, err = crpc.InitRPCClient(conf.RPCServer)
	if err != nil {
		log.Println("RPC client dialing err:", err)
		return err
	}
	return nil
}
