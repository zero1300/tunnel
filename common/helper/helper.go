package helper

import (
	"log"
	"net"
	"time"
)

func NewListener(addr string) *net.TCPListener {
	TCPAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		panic(err.Error())
	}
	listener, err := net.ListenTCP("tcp", TCPAddr)
	return listener
}

func NewConn(addr string) (*net.TCPConn, error) {
	TCPAddr, err := net.ResolveTCPAddr("tcp", addr)
	if err != nil {
		return nil, err
	}
	conn, err := net.DialTCP("tcp", nil, TCPAddr)
	return conn, err
}

func KeepAlive(conn *net.TCPConn, keepAliveStr string) {
	for {
		_, err := conn.Write([]byte(keepAliveStr))
		if err != nil {
			log.Printf("KeepAlive Error %s", err)
			return
		}
		time.Sleep(60 * time.Second)
	}
}
