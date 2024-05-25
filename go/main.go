package main

import (
	"chat_server_golang/config"
	"chat_server_golang/network"
	"flag"
	"log"
)

func init() {
	log.Println("먼저 시작 됩니다.")
}

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set ")

func main() {
	c := config.NewConfig(*pathFlag)

	log.Println("나중에 시작 됩니다.")
	n := network.NewServer()
	n.StartServer()
}
