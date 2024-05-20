package main

import (
	"chat_server_golang/network"
	"log"
)

func init() {
	log.Println("먼저 시작 됩니다.")
}

func main() {
	log.Println("나중에 시작 됩니다.")
	n := network.NewServer()
	n.StartServer()
}
