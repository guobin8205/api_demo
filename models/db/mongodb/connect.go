package mongodb

import (
	"net"
	"github.com/go-mgo/mgo"
)

func NewMongoDB(host, port string) *mgo.Session {
	session, err := mgo.Dial(net.JoinHostPort(host, port))
	if err != nil {
		panic(err)
	}
	return session
}
