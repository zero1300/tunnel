package test

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"net"
	"sync"
	"testing"
	"tunnel/common/helper"
)

const ControlServerAddr = "0.0.0.0:8080"
const RequestServerAddr = "0.0.0.0:8081"
const KeepAliveStr = "KeepAlive\n"

var clientConn *net.TCPConn
var wg sync.WaitGroup

func TestUserServer(t *testing.T) {
	wg.Add(1)
	go ControlServer()
	go RequestServer()
	wg.Wait()
}

func ControlServer() {
	listener := helper.NewListener(ControlServerAddr)
	log.Println("ControlServer Started")
	var err error
	for {
		clientConn, err = listener.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
		}
		go helper.KeepAlive(clientConn, KeepAliveStr)
	}
}

func RequestServer() {
	listener := helper.NewListener(RequestServerAddr)
	log.Println("RequestServer Started")
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			fmt.Println(err.Error())
			return
		}
		go func() {
			_, err := io.Copy(clientConn, conn)
			if err != nil {

			}
		}()
		go func() {
			_, err := io.Copy(conn, clientConn)
			if err != nil {

			}
		}()
	}
}

// 客户端
func TestUserClient(t *testing.T) {
	conn, err := helper.NewConn(ControlServerAddr)
	if err != nil {
		log.Println("连接失败....")
	}
	for {
		s, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			return
		}
		log.Println("获取到数据:" + s)
	}
}

// 用户端
func TestUserReq(t *testing.T) {
	conn, err := helper.NewConn(RequestServerAddr)
	if err != nil {
		log.Println("连接失败....")
	}
	_, err = conn.Write([]byte("Hello there\n"))
	if err != nil {
		log.Println(err.Error())
	}
}
