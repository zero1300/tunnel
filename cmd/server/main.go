package main

import (
	"context"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"net"
	"os/signal"
	"sync"
	"syscall"
	"tunnel/common/constants"
	"tunnel/common/helper"
)

var userConn *net.TCPConn
var controlConn *net.TCPConn
var wg sync.WaitGroup

func init() {
	viper.SetConfigName("server")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("加载配置文件失败: " + err.Error())
	}
}

func main() {
	wg.Add(1)
	go ControlAliveServerListen()
	go userRequestListen()
	go tunnelServerListen()
	go func() {
		notifyContext, stop := signal.NotifyContext(context.Background(), syscall.SIGINT, syscall.SIGTERM)
		defer stop()
		<-notifyContext.Done()
		wg.Done()
	}()
	wg.Wait()

}

func ControlAliveServerListen() {
	controlServerAddr := viper.GetString("server.controlServerAddr")
	listener := helper.NewListener(controlServerAddr)
	log.Println("控制服务正在监听中: " + controlServerAddr)
	var err error
	for {
		controlConn, err = listener.AcceptTCP()
		if err != nil {
			log.Println("KeepAliveServerListen AcceptTCP error: " + err.Error())
			return
		}
		go helper.KeepAlive(controlConn, constants.KeepAliveStr)
	}
}

// 接受用户的请求, 并吧请求转发到tunnel
func userRequestListen() {
	userRequestServerAddr := viper.GetString("server.userRequestServerAddr")
	listener := helper.NewListener(userRequestServerAddr)
	log.Println("用户请求服务正在监听中: " + userRequestServerAddr)
	var err error
	for {
		userConn, err = listener.AcceptTCP()
		if err != nil {
			log.Println("userRequestListen AcceptTCP error: " + err.Error())
			return
		}
		// 发送一个消息, 告诉客户端有一个新的连接
		if controlConn != nil {
			fmt.Println("userRequestServer 接收到到一个用户请求.....")
			_, _ = controlConn.Write([]byte(constants.NewConnection))
		}
	}
}

func tunnelServerListen() {
	tunnelServerAddr := viper.GetString("server.tunnelServerAddr")
	listener := helper.NewListener(tunnelServerAddr)
	log.Println("隧道服务正在监听中: " + tunnelServerAddr)
	for {
		conn, err := listener.AcceptTCP()
		if err != nil {
			log.Println("tunnelServerListen AcceptTCP error: " + err.Error())
			return
		}
		go io.Copy(conn, userConn)
		go io.Copy(userConn, conn)
	}
}
