package config

import (
	"errors"
	"log"

	"github.com/BurntSushi/toml"
)

type Config struct {
	Title   string
	Port    int
	Logging LogInfo
	Users   map[string]User
}

type LogInfo struct {
	Adapter  string `toml:"adapter"`
	Database string `toml:"database"`
}

type User struct {
	Nick     string
	AltNick  string
	AuthInfo Auth `toml:"auth"`
	Certs    map[string]Cert
	Networks map[string]Network
}

type Cert struct {
	Path string `toml:"cert_path"`
}

type Network struct {
	Servers     []string
	PerformInfo Perform `toml:"perform"`
}

type Perform struct {
	Channels []string
	Commands []string
}

type Auth struct {
	CAPath string `toml:"ca_path"`
}

func LoadConfig(configFile string) (Config, error) {
	var c Config
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		log.Fatalf("Cannot load config file! Error: %s\n", err)
	}

	// TODO: Config validation, default values, etc

	if c.Title == "Zounce Configuration" {
		return c, nil
	} else {
		return Config{}, errors.New("Bad stuff")
	}
}
