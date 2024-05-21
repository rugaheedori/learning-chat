package network

import "github.com/gorilla/websocket"

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
