package main

import (
	"chat_server_golang/config"
	"chat_server_golang/network"
	"chat_server_golang/repository"
	"chat_server_golang/service"
	"flag"
)

var pathFlag = flag.String("config", "./config.toml", "config set")
var port = flag.String("port", ":1010", "port set")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)

	if rep, err := repository.NewRepository(c); err != nil {
		panic(err)
	} else {
		s := network.NewServer(service.NewService(rep), rep, *port)
		s.StartServer()
	}
}
