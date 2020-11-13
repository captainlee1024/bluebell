package controller

// 专门用来放接口文档用到的model
// 因为我们的接口文档返回的数据格式一致，但是具体的data类型不一致

type _ResponsePostList struct {
	Code ResCode     `json:"code"`           // 业务状态码
	Msg  interface{} `json:"msg"`            // 提示信息
	Data interface{} `json:"data,omitempty"` // 数据
}
