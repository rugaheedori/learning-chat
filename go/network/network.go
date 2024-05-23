package network

import (
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Network struct {
	engin *gin.Engine
}

// 프레임워크를 사용할 수 있는 객체값 리턴함수 생성
func NewServer() *Network {
	n := &Network{engin: gin.New()}

	// n.engin.Use: app.use와 같은 모든 라우터에 대해 특정 범용처리 하는 부분
	n.engin.Use(gin.Logger())
	// gin.Recovery: 오류로 인해 서버가 죽은 경우 자동으로 서버를 다시 올려주는 역할
	n.engin.Use(gin.Recovery())
	n.engin.Use(cors.New(cors.Config{
		AllowWebSockets:  true,
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT"},
		AllowHeaders:     []string{"*"},
		AllowCredentials: true,
	}))

	r := NewRoom()
	// go rutine? 고루틴 : 백그라운드에서 동하라는 명령어
	go r.RunInit()

	n.engin.GET("/room", r.SocketServe)

	return n
}

// 서버 시작 함수
func (n *Network) StartServer() error {
	log.Println("Starting Server")

	return n.engin.Run(":8080")
}
