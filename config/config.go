package config

import (
	"fmt"
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
	adapter  string
	database string
}

type User struct {
	Nick     string
	AltNick  string
	AuthInfo Auth `toml:"auth"`
	Certs    map[string]Cert
	Networks map[string]Network
}

type Cert struct {
	cert_path string
}

type Network struct {
	Servers     []string
	PerformInfo Perform `toml:"perform"`
}

type Perform struct {
	channels []string
	commands []string
}

type Auth struct {
	ca_path string
}

func LoadConfig(configFile string) (Config, error) {
	var c Config
	if _, err := toml.DecodeFile(configFile, &c); err != nil {
		log.Fatalf("Cannot load config file! Error: %s\n", err)
	}
	fmt.Printf("%#v\n", c)
	return c, nil
}
