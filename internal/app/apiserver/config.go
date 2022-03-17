package apiserver

import "webapp/internal/app/store"

//структура конфига
type Config struct {
	BindAddress string `toml:"bind_addres"`
	Store       *store.Config
}

//возвращает структуру конфига
func NewConfig() *Config {
	return &Config{
		BindAddress: ":8080",
		Store:       store.NewConfig(),
	}
}
