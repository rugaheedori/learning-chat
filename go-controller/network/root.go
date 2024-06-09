package network

import (
	"controller_server_golang/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine

	service *service.Service

	port string
}

// 프레임워크를 사용할 수 있는 객체값 리턴함수 생성
func NewNetwork(service *service.Service, port string) *Server {
	s := &Server{engine: gin.New(), service: service, port: port}

	// s.engine.Use: app.use와 같은 모든 라우터에 대해 특정 범용처리 하는 부분
	s.engine.Use(gin.Logger())
	// gin.Recovery: 오류로 인해 서버가 죽은 경우 자동으로 서버를 다시 올려주는 역할
	s.engine.Use(gin.Recovery())
	s.engine.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))

	registerTowerAPI(s)

	return s
}

func (s *Server) Start() error {
	log.Printf("Start Tx Server")

	return s.engine.Run(s.port)
}
