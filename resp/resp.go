package resp

import "time"

type RespCode int32

const (
	NONE RespCode = iota
)

type Resp struct {
	Timestamp int64       `json:"timestamp"`
	Msg       string      `json:"msg"`
	Code      RespCode    `json:"code"`
	Data      interface{} `json:"data"`
}

func Ok(data ...interface{}) *Resp {
	return &Resp{
		Timestamp: time.Now().Unix(),
		Msg:       "OK",
		Code:      NONE,
		Data:      data,
	}
}
