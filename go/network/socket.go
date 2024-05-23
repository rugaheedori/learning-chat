package network

import (
	"chat_server_golang/types"
	"log"
	"net/http"
	"time"

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

func (c *Client) Read() {
	// 클라이언트가 들어오는 메세지를 읽는 함수
	defer c.Socket.Close()
	for {
		var msg *message
		log.Println("Read", msg, "Name", c.Name)
		err := c.Socket.ReadJSON(&msg)
		if err != nil {
			// Close가 됐음에도 동시성 이슈로 인한 ReadJSON 실행 시 socket 끊김 에러 처리
			if !websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway) {
				break
			} else {
				panic(err)
			}
		} else {
			msg.Time = time.Now().Unix()
			msg.Name = c.Name

			c.Room.Forward <- msg
		}
	}
}

func (c *Client) Write() {
	defer c.Socket.Close()
	// 클라이언트가 메세지를 전송하는 함수
	for msg := range c.Send {
		log.Println("Write", msg, "Name", c.Name)
		err := c.Socket.WriteJSON(msg)

		if err != nil {
			panic(err)
		}
	}
}

func (r *Room) RunInit() {
	// Room에 있는 모든 채널값을 받는 역할
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true
		case client := <-r.Leave:
			r.Clients[client] = false
			delete(r.Clients, client)
			close(client.Send) // 채널을 닫아주는 역할
		case msg := <-r.Forward:
			for client := range r.Clients {
				client.Send <- msg
			}
		}
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

	go client.Write()

	client.Read()
}
