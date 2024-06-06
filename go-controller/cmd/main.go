package main

import (
	"controller_server_golang/cmd/app"
	"controller_server_golang/config"
	"flag"
)

var pathFlag = flag.String("config", "../config.toml", "config set")

func main() {
	flag.Parse()
	c := config.NewConfig(*pathFlag)

	// todo app 객체를 사용하여 서버 시작
	a := app.NewApp(c)
	a.Start()
}
