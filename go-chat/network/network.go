package network

import (
	"chat_server_golang/service"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"syscall"

	"github.com/confluentinc/confluent-kafka-go/kafka"
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type Server struct {
	engine *gin.Engine

	service *service.Service

	port string
	ip   string
}

// 프레임워크를 사용할 수 있는 객체값 리턴함수 생성
func NewServer(service *service.Service, port string) *Server {
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

	registerServer(s)

	return s
}

func (s *Server) setServerInfo() {
	// IP를 가져오고, IP를 기반으로 MySQL serverInfo 테이블 변경
	if addrs, err := net.InterfaceAddrs(); err != nil {
		panic(err.Error())
	} else {
		var ip net.IP

		for _, addr := range addrs {
			// 타입 상속 확인
			if ipnet, ok := addr.(*net.IPNet); ok {
				if !ipnet.IP.IsLoopback() && ipnet.IP.To4() != nil {
					ip = ipnet.IP
					break
				}
			}
		}

		if ip == nil {
			panic("no ip address found")
		} else {
			if err = s.service.ServerSet(ip.String()+s.port, true); err != nil {
				panic(err)
			} else {
				s.ip = ip.String()
			}
		}
	}
}

// 서버 시작 함수
func (s *Server) StartServer() error {
	s.setServerInfo()

	// 일종의 이벤트를 받을 수 있는 변수를 선언
	channel := make(chan os.Signal, 1)
	// 서버가 죽었을 때 감지하여 채널에 메세지를 전송함
	signal.Notify(channel, syscall.SIGINT)

	go func() {
		<-channel // 서버가 죽었다는 의미

		if err := s.service.ServerSet(s.ip+s.port, false); err != nil {
			// todo 실패 케이스에 대해 추가 처리 필요 ex) retry option
			log.Println("Failed to Set Server Into When Close", "err", err)
		}

		// Kafka에 이벤트 전송

		type ServerInfoEvent struct {
			IP     string
			Status bool
		}

		e := &ServerInfoEvent{
			IP:     s.ip + s.port,
			Status: false,
		}
		ch := make(chan kafka.Event)

		// 값을 전송할 때는 배열 바이트 값으로 전송해야 함
		if v, err := json.Marshal(e); err != nil {
			log.Println("Failed To Marshal")
		} else if result, err := s.service.PublishEvent("chat", v, ch); err != nil {
			// TODO Send Event To Kafka
			log.Println("Failed To Send Event to Kafka", "err", err)
		} else {
			log.Println("Success To Send Event", result)
		}
	}()

	log.Println("Starting Server")
	return s.engine.Run(s.port)
}
