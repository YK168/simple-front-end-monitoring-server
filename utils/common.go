package utils

// 基础序列化器
type Response struct {
	Status int    `json:"status"`
	Data   any    `json:"data"`
	Msg    string `json:"msg"`
	Error  string `json:"error"`
}

// 带token的返回值
type TokenData struct {
	User  any    `json:"user"`
	Token string `json:"token"`
}
