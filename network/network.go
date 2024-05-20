package network

import (
	"log"

	"github.com/gin-gonic/gin"
)

type Network struct {
	engin *gin.Engine
}

// 프레임워크를 사용할 수 있는 객체값 리턴함수 생성
func NewServer() *Network {
	n := &Network{engin: gin.New()}

	return n
}

// 서버 시작 함수
func (n *Network) StartServer() error {
	log.Println("Starting Server")

	return n.engin.Run(":8080")
}
