/**
 * Created by IDEA.
 * User: godsoul
 */
package rds

import (
	"errors"
	"log"

	"github.com/garyburd/redigo/redis"
	"reflect"
	"strings"
)

func (r *Rds) Set(k string, v interface{}) error {
	_, err := r.Con.Do("SET", k, v)
	return err
}

//注意这里的参数顺序
func (r *Rds) SetEx(k string, exp int, v interface{}) error {
	_, err := r.Con.Do("SETEX", k, exp, v)
	return err
}

func (r *Rds) Get(k string) (string, error) {
	return redis.String(r.Con.Do("GET", k))
}

func (r *Rds) HGet(key, field string) (interface{}, error) {
	return r.Con.Do("HGET", key, field)
}

func (r *Rds) HSet(key, field string, value interface{}) error {
	_, err := redis.Int(r.Con.Do("HSET", key, field, value))
	return err
}

func (r *Rds) HMSet(key string, args interface{}) error {
	_, err := r.Con.Do("HMSET", redis.Args{}.Add(key).AddFlat(args)...)
	return err
}


func (r *Rds) HINCRBY(key, field string, value int) (int, error) {
	return redis.Int(r.Con.Do("HINCRBY", key, field, value))
}

func (r *Rds) Del(k string) error {
	_, err := r.Con.Do("DEL", k)
	return err
}

func (r *Rds) Expire(k string, d int) error {
	_, err := r.Con.Do("EXPIRE", k, d)
	return err
}

func (r *Rds) Incr(key string) (int, error) {
	return redis.Int(r.Con.Do("INCR", key))
}

func (r *Rds) Exists(k string) bool {
	e, err := redis.Bool(r.Con.Do("EXISTS", k))
	if err != nil {
		return false
	}
	return e
}

func (r *Rds) HGetAll(key string) map[string]string {
	re, err := redis.StringMap(r.Con.Do("HGETALL", key))
	if err != nil && err != redis.ErrNil {
		log.Println(err)
	}
	return re
}

func (r *Rds) Map2Struct(key string, obj interface{}) error {
	src, err := redis.Values(r.Con.Do("HGETALL", key))
	if err != nil {
		return err
	}
	if len(src) == 0 {
		return errors.New("key not exists")
	}
	return redis.ScanStruct(src, obj)
}

//d 阻塞时间,单位秒
func (r *Rds) BLPOP(key string, d int) []string {
	re, err := redis.Strings(r.Con.Do("BLPOP", key, d))
	if err != nil {
		log.Println(err)
	}
	return re
}

func (r *Rds) LPOP(key string) (interface{}, error) {
	return r.Con.Do("LPOP", key)
}

func (r *Rds) LRange(key string,begin ,end int)([]interface{}, error) {
	return  redis.Values(r.Con.Do("lrange",key,begin,end))
}

func (r *Rds) SisMember(key string, m int) (bool, error) {
	return redis.Bool(r.Con.Do("SISMEMBER", key, m))
}

func (r *Rds) SMembers_Ints(key string) ([]int, error) {
	return redis.Ints(r.Con.Do("SMEMBERS", key))
}

func (r *Rds) SMembers_Strings(key string) ([]string, error) {
	return redis.Strings(r.Con.Do("SMEMBERS", key))
}

func (r *Rds) SAdd(key, value string) error {
	_, err := r.Con.Do("SADD", key, value)
	return err
}

func (r *Rds) SREM(key, value string) error {
	_, err := r.Con.Do("SREM", key, value)
	return err
}

func (r *Rds) SCard(key string) (int, error) {
	return redis.Int(r.Con.Do("SCARD", key))
}

func (r *Rds) LPush(key, value string) error {
	_, err := r.Con.Do("LPUSH", key, value)
	return err
}

func (r *Rds) RPush(key, value string) error {
	_, err := r.Con.Do("RPUSH", key)
	return err
}

func (r *Rds) Scard(key string) int {
	i, _ := redis.Int(r.Con.Do("SCARD", key))
	return i
}

func (r *Rds) ZAdd(key string, score interface{}, member string) error {
	_, err := r.Con.Do("ZADD", key, score, member)
	return err
}

func (r *Rds) ZCard(key string) (int, error) {
	return redis.Int(r.Con.Do("ZCARD", key))
}

func (r *Rds) ZREVRANK(key string, s int) (int, error) {
	return redis.Int(r.Con.Do("ZREVRANK", key, s))
}

func(r *Rds) ZRANGEBYSCORE(key,min,max string,offset,count int) ([]interface{}, error){
	return redis.Values(r.Con.Do("ZRANGEBYSCORE", key, "("+min, max,"limit",offset,count))
}

func(r *Rds) ZREM(key,value string) error{
	_,err := r.Con.Do("ZREM", key, value)
	return err
}

func (r *Rds) HMGet(key string,  args... string) ([]interface{}, error) {
	data, err := redis.Values(r.Con.Do("hmget", redis.Args{}.Add(key).AddFlat(args)...))
	return data,err
}

func (r *Rds) GetHashList(keys map[string][]string)(map[string]map[string]interface{},error){
	lkeys := len(keys)
	if lkeys < 1{
		return nil,errors.New("key is empty")
	}
	res := make(map[string]map[string]interface{})
	var err error
	var allKeys [][]interface{}
	for k, v := range keys {
		var key []interface{}
		if len(v) > 0 { //取部分
			key = append(key, k)
			for _, str := range v {
				key = append(key, str)
			}
			err = r.Con.Send("HMGET", key...)
			if err != nil {
				key = make([]interface{}, 0) //有错误将此值赋空值
			}

		}
		allKeys = append(allKeys, key)
	}
	if err = r.Con.Flush(); err != nil {
		return nil, err
	}
	var reply interface{}
	for i := 0; i < lkeys; i++ {
		if len(allKeys[i]) > 0 {
			if reply, err = r.Con.Receive(); err == nil {
				redisKey := allKeys[i][:1][0].(string)
				hashKeys := allKeys[i][1:]
				res[redisKey] = r.buildKeyVals(hashKeys, reply)
			}
		}
	}
	return res, nil
}

func(r *Rds) buildKeyVals(keys []interface{}, vals interface{}) map[string]interface{} {
	rv := reflect.Indirect(reflect.ValueOf(vals))
	rt := strings.Replace(rv.Type().String(), " ", "", -1)
	var res map[string]interface{} = make(map[string]interface{})
	if rt == "[]interface{}" {
		valsArr := vals.([]interface{})
		if len(keys) == len(valsArr) {
			for k, v := range keys {
				mKey, _ := redis.String(v, nil)
				if mKey != "" {
					res[mKey] = valsArr[k]
				}
			}

		}
	}
	return res
}
