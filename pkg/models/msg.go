package models

//响应
type GeneralRes struct {
	Code int64 `json:"code"`
	Msg string `json:"msg"`
}

const(
	ControlConn = iota
	WorkConn
)


//客户端请求
type ClientCtlReq struct {
	Type int64 `json:"type"`
	ProxyName string `json:"proxy_name"`
	Passwd string `json:"passwd"`
}

type ClientCtlRes struct {
	GeneralRes
}

type ServerCtlReq struct {
	Type int64 `json:"type"`
}