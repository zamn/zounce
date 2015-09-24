package config

import "fmt"

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

func LoadConfig(configPath string) (Config, error) {
	fmt.Println(configPath)
	return Config{}, nil
}
