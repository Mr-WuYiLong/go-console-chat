package main

import (
	"chat_system/client/service"
	"fmt"
	"os"
)

func main() {
	var key int
	var username, password string
	userService := new(service.UserService)
	fmt.Println("      欢迎进入聊天系统")
	for {
		fmt.Println("      1.登录")
		fmt.Println("      2.注册")
		fmt.Println("      3.退出系统")
		fmt.Println("      请选择1~3选项")
		fmt.Scanln(&key)
		switch key {
		case 1:
			fmt.Println("输入用户名:")
			fmt.Scanln(&username)
			fmt.Println("输入密码:")
			fmt.Scanln(&password)
			err := userService.Login(username, password)
			if err != nil {
				fmt.Printf("*****登录失败err->%v\n", err)
			}

		case 2:
			fmt.Println("输入用户名:")
			fmt.Scanln(&username)
			fmt.Println("输入密码:")
			fmt.Scanln(&password)
			err := userService.Register(username, password)
			if err != nil {
				fmt.Printf("*****注册失败err->%v\n", err)
				break
			}
			fmt.Println("注册成功")

		case 3:
			os.Exit(0) // 退出系统
		default:
			fmt.Println("输入有误,请选择正确的选项")
		}
	}

}
