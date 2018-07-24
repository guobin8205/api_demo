/**
 * Created by IDEA.
 * User: godsoul
 */
package rds

import (
	"time"
	"net"
	"github.com/garyburd/redigo/redis"
)

type RdsPool struct {
	Host      string
	Port      string
	Password  string
	Maxidle   int
	Maxactive int
	pool      *redis.Pool
}

//对常用命令的二次封装
type Rds struct {
	Con redis.Conn
}

func NewRdsPool(host, port, password string, maxidle, maxactive int) *RdsPool {
	pool := &redis.Pool{
		MaxIdle:     maxidle,
		MaxActive:   maxactive,
		IdleTimeout: 180 * time.Second,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial(
				"tcp",
				net.JoinHostPort(host, port),
				redis.DialPassword(password),
			)
			if err != nil {
				panic(err)
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	return &RdsPool{
		Host:      host,
		Port:      port,
		Password:  password,
		Maxidle:   maxidle,
		Maxactive: maxactive,
		pool:      pool,
	}
}

func (pool *RdsPool) Get() Rds {
	c := pool.pool.Get()
	return Rds{Con: c}
}

func (r *Rds) Close() {
	r.Con.Close()
}
