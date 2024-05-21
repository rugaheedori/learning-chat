package network

import (
	"chat_server_golang/types"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// HTTP Connection을 Websocket Connection으로 upgrade 해줌 & CheckOrigin: 연결 시 오리진을 체크
var upgrader = &websocket.Upgrader{ReadBufferSize: types.SocketBufferSize, WriteBufferSize: types.MessageBufferSize, CheckOrigin: func(r *http.Request) bool { return true }}

// 채팅방
type Room struct {
	Forward chan *message // 수신되는 메세지를 보관하는 값

	Join  chan *Client // Socket이 연결되는 경우에 동작
	Leave chan *Client // Socket이 끊어지는 경우에 동작

	Clients map[*Client]bool // 현재 방에 있는 Client의 정보를 저장
}

type message struct {
	Name    string
	Message string
	Time    int64
}

type Client struct {
	Send   chan *message
	Room   *Room
	Name   string
	Socket *websocket.Conn
}

// Room객체 필드들을 초기화해주는 함수
func NewRoom() *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
	}
}

// gin 사용 시 API 연결 가능하게 해줌
func (r *Room) SocketServe(c *gin.Context) {
	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		panic(err)
	}

	userCookie, err := c.Request.Cookie("auth")
	if err != nil {
		panic(err)
	}

	client := &Client{
		Socket: socket,
		Room:   r,
		Name:   userCookie.Value,
		Send:   make(chan *message, types.MessageBufferSize),
	}

	r.Join <- client

	// 함수가 종료될 때 실행 됨
	defer func() {
		r.Leave <- client
	}()
}
