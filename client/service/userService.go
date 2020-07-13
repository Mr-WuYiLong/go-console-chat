package service

import (
	"bufio"
	"chat_system/client/exception"
	"chat_system/model"
	"chat_system/utils"
	"encoding/json"
	"fmt"
	"io"
	"net"
	"os"
	"strings"

	"github.com/gomodule/redigo/redis"
)

// UserService 用户服务结构体
type UserService struct {
	onlineUser []string
}

// func clientWrite(conn net.Conn) {
// 	reader := bufio.NewReader(os.Stdin) // 标准输入        x
// 	for {
// 		data := reader.ReadString('\n')
// 		_, err := conn.Write([]byte(data))
// 		if err != nil {
// 			fmt.Printf("*****客户端发送数据出错****err->%v", err)
// 		}
// 	}
// }

// Register 用户注册
func (userService *UserService) Register(username string, password string) (err error) {
	r := utils.GetRedis()
	defer r.Close()
	userVal, err := redis.Strings(r.Do("HVALS", "user"))
	if err != nil {
		fmt.Printf("从redis获取所有的值出错err->%v", err)
	}
	for _, v := range userVal {
		user := new(model.User)
		json.Unmarshal([]byte(v), user)
		if username == user.Username {
			return exception.NewException("用户已被注册")
		}
	}

	// 进行注册
	user := model.NewUser(username, password)
	userSlice, err := json.Marshal(*user)
	if err != nil {
		fmt.Printf("json序列化出错err->%v", err)
	}
	_, err = r.Do("HSET", "user", user.ID, string(userSlice))
	if err != nil {
		fmt.Printf("用户redis-存储失败err->%v", err)
	}
	return nil
}

// Login 登录
func (userService *UserService) Login(username string, password string) (err error) {
	r := utils.GetRedis()
	defer r.Close()
	userVal, err := redis.Strings(r.Do("HVALS", "user"))
	if err != nil {
		fmt.Printf("从redis获取所有的值出错err->%v", err)
	}
	var loop = false
	for _, v := range userVal {
		user := new(model.User)
		json.Unmarshal([]byte(v), user)
		if username == user.Username {
			if password == user.Password {
				loop = true
			}
		}
	}
	if !loop {
		return exception.NewException("用户名或密码错误")
	}

	//查找所有在线用户
	onlineUsers, _ := redis.Strings(utils.GetRedis().Do("SMEMBERS", "onlineUser"))
	for _, v := range onlineUsers {
		if username == v {
			return exception.NewException("你的账号已登录")
		}
	}
	// 连接服务端
	conn, err := net.Dial("tcp", "127.0.0.1:8080")
	if err != nil {
		fmt.Printf("客户端无法连接服务端err->%v", err)
		return err
	}

	fmt.Printf("*****登录成功****\n")

	user := model.NewUser(username, password)
	// fmt.Printf("序列化后的数据->%v\n", user)
	// 序列化
	userSlice, err := json.Marshal(*user)
	if err != nil {
		fmt.Printf("json序列化出错err->%v", err)
	}
	// 序列化发送的消息体
	message := model.NewMessage("login", string(userSlice), len(userSlice), username)
	msgSlice, err := json.Marshal(*message)
	if err != nil {
		fmt.Printf("json序列化出错err->%v", err)
	}
	_, err = conn.Write(msgSlice)
	if err != nil {
		fmt.Printf("客户端写入出错err->%v", err)
	}
	// 从服务端读取数据
	go func() {

		for {
			num := make([]byte, 1024)
			n, err := conn.Read(num)
			if err == io.EOF {
				fmt.Printf("服务端->%v已退出err->%v\n", conn.RemoteAddr(), err)
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

			if msg.Type == "online" {
				fmt.Println("提示：默认群聊")
				fmt.Println("******在线列表*******")
				userService := new(UserService)
				userService.onlineUser = strings.Split(msg.Data, ",")
				// fmt.Printf("服务端传过来的数据->%v\n", userService.onlineUser)

				// userService := new(UserService)
				if userService.onlineUser != nil {
					for i, v := range userService.onlineUser {
						fmt.Printf("%v.%v\n", i+1, v)
					}
				}
			} else {
				fmt.Printf("%v:%v\n", msg.Sender, msg.Data)
			}

		}
	}()

	for {

		reader := bufio.NewReader(os.Stdin) // 标准输入
		data, _ := reader.ReadString('\n')
		data = strings.Trim(data, "\n")
		msg := model.NewMessage("msg", data, len(data), username)
		msgSlice, _ := json.Marshal(*msg)
		_, err := conn.Write([]byte(msgSlice))
		if err != nil {
			fmt.Printf("*****客户端发送数据出错****err->%v", err)
		}

	}

	return nil
}
