package main

import (
	"flag"
	"fmt"
	"webapp/internal/app/apiserver"

	"github.com/BurntSushi/toml"
	"github.com/gorilla/sessions"
)

var configPath string

func init() {
	//CLI флаг для передачи пути на config файл
	flag.StringVar(&configPath, "config-path", "configs/apiserver.toml", "PATH TO CONFIG API SERVER")
}

func main() {
	config := apiserver.NewConfig()
	_, err := toml.DecodeFile(configPath, config)
	if err != nil {
		return
	}

	sessionStore := sessions.NewCookieStore([]byte(config.SessionKey))
	server := apiserver.NewServer(config, sessionStore)

	if err := server.Start(); err != nil {
		fmt.Println("Error: start server!")
		return
	}
}
