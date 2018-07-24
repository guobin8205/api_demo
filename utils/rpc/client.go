package rpc

import (
	"net/rpc"
	"fmt"
	"time"
)

const(
	pingDuration = 1 * time.Second
)

type NoArg struct {
}

type NoReply struct {
}



func InitRPCClient(addr string) (*rpc.Client,error){
	client, err := rpc.DialHTTP("tcp", addr)
	if err != nil {
		return nil,err
	}
	return client,nil
}

func Call(srv string, rpcname string, args interface{}, reply interface{}) error {
	c, errx := rpc.Dial("tcp", srv)
	if errx != nil {
		return fmt.Errorf("ConnectError: %s", errx.Error())
	}
	defer c.Close()
	return c.Call(rpcname, args, reply)
}