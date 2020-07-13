package model

import (
	"strconv"
	"time"
)

// User 用户
type User struct {
	ID       string `json:"id"`
	Username string `json:"username"`
	Password string `json:"password"`
}

// NewUser 实例化
func NewUser(username string, password string) *User {

	// 获得当前的秒
	nowTime := time.Now().Unix()
	return &User{
		ID:       strconv.FormatInt(nowTime, 16),
		Username: username,
		Password: password,
	}
}
