package types

import "strings"

// 소켓 통신을 할 때 사용할 버퍼 사이즈 설정
const (
	SocketBufferSize  = 1024 // 큰 사이즈 통신이 잦다면 소켓 버퍼 사이즈 크기 늘려주어야 함
	MessageBufferSize = 256  // 이미지, 동영상과 같은 큰 버퍼 사이즈 데이터를 전송해야 하는 경우 크기 늘려주어야 함
)

type LoginReq struct {
	Name string `json:"name" binding:"required"`
}

type header struct {
	Result int    `json:"result"`
	Data   string `json:"data"`
}

func newHeader(result int, data ...string) *header {
	return &header{
		Result: result,
		Data:   strings.Join(data, ","),
	}
}

type response struct {
	*header
	Result interface{} `json:"result"`
}

func NewRes(result int, res interface{}, data ...string) *response {
	return &response{
		header: newHeader(result, data...),
		Result: res,
	}
}
