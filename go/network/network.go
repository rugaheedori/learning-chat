package network

import (
	"chat_server_golang/repository"
	"chat_server_golang/service"
	"log"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engin *gin.Engine

	service    *service.Service
	repository *repository.Repository

	port string
	ip   string
}

// 프레임워크를 사용할 수 있는 객체값 리턴함수 생성
func NewServer(service *service.Service, repository *repository.Repository, port string) *Server {
	s := &Server{engin: gin.New(), service: service, repository: repository, port: port}

	// s.engin.Use: app.use와 같은 모든 라우터에 대해 특정 범용처리 하는 부분
	s.engin.Use(gin.Logger())
	// gin.Recovery: 오류로 인해 서버가 죽은 경우 자동으로 서버를 다시 올려주는 역할
	s.engin.Use(gin.Recovery())
	s.engin.Use(cors.New(cors.Config{
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "PATCH"},
		AllowHeaders:     []string{"*"},
		ExposeHeaders:    []string{"*"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			return true
		},
	}))

	registerServer(s.engin)

	return s
}

// 서버 시작 함수
func (s *Server) StartServer() error {
	log.Println("Starting Server")

	return s.engin.Run(s.port)
}
