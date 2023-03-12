package main

import (
	"bufio"
	"fmt"
	"github.com/spf13/viper"
	"io"
	"log"
	"tunnel/common/constants"
	"tunnel/common/helper"
)

func init() {
	viper.SetConfigName("client")
	viper.SetConfigType("yaml")
	viper.AddConfigPath(".")
	err := viper.ReadInConfig()
	if err != nil {
		panic("加载配置文件失败: " + err.Error())
	}
}

func main() {
	controlServerAddr := viper.GetString("client.controlServerAddr")
	fmt.Println(controlServerAddr)
	conn, err := helper.NewConn(controlServerAddr)
	if err != nil {
		panic("连接服务器失败....")
	}
	fmt.Println("连接成功: ", conn.RemoteAddr().String())
	for {
		data, err := bufio.NewReader(conn).ReadString('\n')
		if err != nil {
			log.Println(err.Error())
			return
		}
		fmt.Println(data)
		if string(data) == constants.NewConnection {
			go messageForward()
		}
	}
}

// messageForward 消息转发
func messageForward() {
	// 连接服务端的隧道
	tunnelServerAddr := viper.GetString("client.tunnelServerAddr")
	tunnelConn, err := helper.NewConn(tunnelServerAddr)
	fmt.Println(tunnelServerAddr)
	if err != nil {
		panic(err.Error())
	}
	localServerAddr := viper.GetString("client.localServerAddr")
	localConn, err := helper.NewConn(localServerAddr)
	if err != nil {
		panic(err.Error())
	}
	fmt.Println(localServerAddr)
	go io.Copy(localConn, tunnelConn)
	go io.Copy(tunnelConn, localConn)
}
