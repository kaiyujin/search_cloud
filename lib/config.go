package lib

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Backlog backlog
}

type backlog struct {
	Url    string `toml:"url"`
	ApiKey string `toml:"api_key"`
}

func Read() Config {
	var conf Config
	if _, err := toml.DecodeFile("./lib/Config.toml", &conf); err != nil {
		panic(err)
	}
	return conf
}
