package apiserver

import "webapp/internal/app/store"

//структура конфига
type Config struct {
	BindAddress string `toml:"bind_addres"`
	Store       *store.Config
	SessionKey  string `toml:"session_key"`
}

//возвращает структуру конфига
func NewConfig() *Config {
	return &Config{
		BindAddress: ":4040",
		Store:       store.NewConfig(),
	}
}
