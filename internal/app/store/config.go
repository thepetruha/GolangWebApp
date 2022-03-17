package store

//структура конфиг для подключения к бд
type Config struct {
	DatabaseURL string `toml:"database_url"`
}

//возращает конфиг
func NewConfig() *Config {
	return &Config{}
}
