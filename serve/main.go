package main

import (
	"chat_system/model"
	"chat_system/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"strings"

	"github.com/gomodule/redigo/redis"
)

var (
	clientConns = make([]model.ClientConn, 0)
)

// 从客户端读数据
func receive(conn net.Conn) {
	defer conn.Close()
	for {
		num := make([]byte, 1024)
		n, err := conn.Read(num)
		if err == io.EOF {
			fmt.Printf("客户端->%v已退出err->%v\n", conn.RemoteAddr(), err)
			for i, v := range clientConns {
				if conn == v.CliConn {
					clientConns = append(clientConns[:i], clientConns[i+1:]...)
					// 删除redis对应的在线名单
					utils.GetRedis().Do("SREM", "onlineUser", v.Username)

				}
				_, err := v.CliConn.Write(showOnlineUser())
				if err != nil {
					fmt.Printf("****服务器端推送数据失败err->%v", err)
				}

			}
			break
		}
		msg := new(model.Message)
		err = json.Unmarshal(num[:n], msg)
		if err != nil {
			fmt.Printf("服务端反序列化失败err->%v", err)
		}

		if len(msg.Data) != msg.Length {
			fmt.Printf("客户端发包不全")
		}

		if msg.Type == "login" {
			user := new(model.User)
			err = json.Unmarshal([]byte(msg.Data), user)
			if err != nil {
				fmt.Printf("服务端反序列化失败err->%v", err)
			}
			clientConn := model.ClientConn{CliConn: conn, Username: user.Username}
			clientConns = append(clientConns, clientConn)
			utils.GetRedis().Do("SADD", "onlineUser", user.Username)
			fmt.Printf("%v->客户端登录了->用户名:%v\n", conn.RemoteAddr(), user.Username)
			// 登录成功就推送在线用户列表
			for _, v := range clientConns {
				_, err := v.CliConn.Write(showOnlineUser())
				if err != nil {
					fmt.Printf("****服务器端推送数据失败err->%v", err)
				}
			}
		} else {

			for _, v := range clientConns {

				if v.CliConn != conn {

					msg := new(model.Message)
					err = json.Unmarshal(num[:n], msg)
					if err != nil {
						fmt.Printf("服务端反序列化失败err->%v", err)
					}
					msg.Sender = msg.Sender
					msgByte, _ := json.Marshal(msg)
					_, err := v.CliConn.Write(msgByte)
					if err != nil {
						fmt.Printf("****服务器端推送数据失败err->%v", err)
					}
				}

			}

			// fmt.Printf("%v客户端传过来的数据->%v\n", conn.RemoteAddr(), string(num[:n]))
		}

	}

}

// 推送数据到客户端
func write() {
	var content string
	for {
		fmt.Scanln(&content)
		msg := new(model.Message)
		msg.Type = "msg"
		msg.Data = content
		msg.Length = len(content)
		msgByte, _ := json.Marshal(msg)
		for _, v := range clientConns {
			_, err := v.CliConn.Write(msgByte)
			if err != nil {
				fmt.Printf("****服务器端推送数据失败err->%v", err)
			}
		}
	}

}

// 获取在线人数
func showOnlineUser() []byte {
	userList, _ := redis.Strings(utils.GetRedis().Do("SMEMBERS", "onlineUser"))
	// userByte, err := json.Marshal(userList)
	// if err != nil {
	// 	fmt.Printf("数据序列化出错err->%v", err)
	// }
	msg := new(model.Message)
	msg.Type = "online"
	msg.Data = strings.Join(userList, ",")
	msg.Length = len(strings.Join(userList, ","))
	msgByte, _ := json.Marshal(msg)
	return msgByte
}

// 主函数入口
func main() {
	ls, err := net.Listen("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("服务端监听失败->%v", err)
		return
	}

	go write()
	for {
		fmt.Println("等待客户端的连接.....")
		conn, err := ls.Accept()
		if err != nil {
			fmt.Printf("客户端连接不上->%v", err)
		}

		go receive(conn)
		fmt.Printf("客户端%v正在连接服务端\n", clientConns)
	}

}
