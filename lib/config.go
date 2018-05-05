package lib

import (
	"github.com/BurntSushi/toml"
)

type Config struct {
	Backlog    backlog
	Github     github
	Confluence confluence
	PageSize   string `toml:"page_size"`
}

type backlog struct {
	Url    string `toml:"url"`
	ApiKey string `toml:"api_key"`
}

type github struct {
	Url                string `toml:"url"`
	ApiKey             string `toml:"api_key"`
	OptionalConditions string `toml:"optional_conditions"`
}

type confluence struct {
	Url                string `toml:"url"`
	OptionalConditions string `toml:"optional_conditions"`
	Id                 string `toml:"id"`
	Pass               string `toml:"pass"`
}

func Read() Config {
	var conf Config
	if _, err := toml.DecodeFile("./config.toml", &conf); err != nil {
		panic(err)
	}
	return conf
}
