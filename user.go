package main

// User 客户端用户操作数据
type User struct {
	Sip       string `json:"sip"`
	Region    string `json:"region"`
	Namespace string `json:"namespace"`
}
