package network

import (
	"chat_server_golang/service"
	"chat_server_golang/types"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

// 채팅방
type Room struct {
	Forward chan *message // 수신되는 메세지를 보관하는 값
	// 들어오는 메세지를 다른 모든 클라이언트로 보내는 데 사용한다.

	Join  chan *Client // Socket이 연결되는 경우에 동작
	Leave chan *Client // Socket이 끊어지는 경우에 동작

	Clients map[*Client]bool // 현재 방에 있는 Client의 정보를 저장

	service *service.Service
}

type Client struct {
	Socket *websocket.Conn // client의 웹 소켓
	Send   chan *message   // 전송되는 채널
	Room   *Room
	Name   string
}

type message struct {
	Name    string
	Message string
	When    time.Time
}

func (c *Client) Read() {
	// 클라이언트가 ReadMessage 메소드를 통해서 소켓에서 읽을 수 있고, 받은 메세지를 Room 타입에게 계속해서 전송을 한다.
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
			msg.When = time.Now()
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

// Room객체 필드들을 초기화해주는 함수
func NewRoom(service *service.Service) *Room {
	return &Room{
		Forward: make(chan *message),
		Join:    make(chan *Client),
		Leave:   make(chan *Client),
		Clients: make(map[*Client]bool),
		service: service,
	}
}

func (r *Room) Run() {
	// Room에 있는 모든 채널값을 받는 역할
	for {
		select {
		case client := <-r.Join:
			r.Clients[client] = true // client가 새로 들어 올 때
		case client := <-r.Leave:
			r.Clients[client] = false
			delete(r.Clients, client) // 나갈 때에는 map값에서 client를 제거
			close(client.Send)        // 이후 client의 socket을 닫음
		case msg := <-r.Forward: // 특정 메세지가 방에 들어오면
			for client := range r.Clients {
				client.Send <- msg // 모든 client에게 전달
			}
		}
	}
}

const (
	SocketBufferSize  = 1024 // 큰 사이즈 통신이 잦다면 소켓 버퍼 사이즈 크기 늘려주어야 함
	MessageBufferSize = 256  // 이미지, 동영상과 같은 큰 버퍼 사이즈 데이터를 전송해야 하는 경우 크기 늘려주어야 함
)

// HTTP Connection을 Websocket Connection으로 upgrade 해줌
// 기본적으로 HTTP에 웹 소켓을 사용하려면, 이와 같이 업그레이드 해주어야 함 -> 재사용 가능하므로 하나만 만들어도 됨
var upgrader = &websocket.Upgrader{ReadBufferSize: types.SocketBufferSize, WriteBufferSize: types.MessageBufferSize}

// gin 사용 시 API 연결 가능하게 해줌
func (r *Room) ServeHTTP(c *gin.Context) {
	// 이후 요청이 들어오게 된다면 Upgrader를 통해서 소켓을 가져 온다.

	upgrader.CheckOrigin = func(r *http.Request) bool {
		return true
	}

	socket, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Fatal("--- serveHTTP", err)
		return
	}

	authCookie, err := c.Request.Cookie("auth")
	if err != nil {
		log.Fatal("auth cookie is failed", err)
		return
	}

	// 문제가 없다면 client를 생성하여 방에 입장했다고 채널에 전송한다.
	client := &Client{
		Socket: socket,
		Send:   make(chan *message, MessageBufferSize),
		Room:   r,
		Name:   authCookie.Value,
	}

	r.Join <- client

	// defer를 통해서 client가 끝날 때를 대비하여 퇴장하는 작업을 연기함
	defer func() {
		r.Leave <- client
	}()

	// 이 후 고루틴을 통해서 Write 실행
	go client.Write()

	// 이 후 메인 루틴에서 read를 실해ㅐㅇ함으로써 해당 요청을 닫는것을 차단 -> 채널을 활용한 연결 활성화
	client.Read()
}
