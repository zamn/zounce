package config

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/BurntSushi/toml"
	"github.com/go-validator/validator"
)

type Config struct {
	Title   string
	Port    int
	Logging LogInfo
	Users   map[string]User `validate:"hasusers"`
}

type LogInfo struct {
	Adapter  string `toml:"adapter" validate:"nonzero"`
	Database string `toml:"database" validate:"nonzero"`
}

type User struct {
	Nick     string `validate:"nonzero,max=9"`
	AltNick  string `validate:"nonzero,max=9"`
	Username string
	Realname string
	AuthInfo Auth            `toml:"auth"`
	Certs    map[string]Cert `validate:"min=1"`
	Networks map[string]Network
}

type Cert struct {
	Path string `toml:"cert_path validate:"nonzero`
}

type Network struct {
	Servers     []string `validate:"min=1"`
	Password    string
	PerformInfo Perform `toml:"perform"`
}

type Perform struct {
	Channels []string
	Commands []string
}

type Auth struct {
	CAPath string `toml:"ca_path" validate:"nonzero"`
}

var errorExpl = map[string]map[error]string{
	"Logging.Adapter":  map[error]string{validator.ErrZeroValue: "An adapter is required. Valid Options: SQLite3"},
	"Logging.Database": map[error]string{validator.ErrZeroValue: "You must specify the name of the logging database."},
	"Nick":             map[error]string{validator.ErrZeroValue: "ERROR [users.%s]: You must specify a nickname in order to connect to an IRC server.", validator.ErrMax: "ERROR [users.%s]: Nickname can only be 9 characters long."},
	"AltNick":          map[error]string{validator.ErrZeroValue: "ERROR [users.%s]: You must specify a alternate nickname in order to connect to an IRC server.", validator.ErrMax: "ERROR [users.%s]: Altenate nickname can only be 9 characters long."},
	"Certs":            map[error]string{validator.ErrZeroValue: "ERROR [users.%s]: You must specify at least one certificate in order to authenticate to zounce.", validator.ErrMin: "ERROR [users.%s]: You must have at least 1 certificate on your user in order to authenticate."},
	"Users":            map[error]string{validator.ErrZeroValue: "You must specify at least one user in order to use to zounce."},
	"AuthInfo.CAPath":  map[error]string{validator.ErrZeroValue: "You must specify the CA for your certificate to verify."},
}

func validateUsers(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	/*
		defValMap := map[string]string{
			"Logging.Adapter":  "SQLite3",
			"Logging.Database": "zounce",
			"Nick":             "zounceuser",
			"AltNick":          "zounceuser-alt",
		}
	*/

	if st.Kind() == reflect.Map {
		keys := st.MapKeys()
		for _, k := range keys {
			isValid, errMap := validator.Validate(st.MapIndex(k).Interface())
			if !isValid {
				for k, v := range errMap {
					for _, err := range v {
						errorMsg := errorExpl[k][err]
						if strings.Contains(errorMsg, "%s") {
							errorMsg = fmt.Sprintf(errorExpl[k][err], k)
						}
						fmt.Println(errorMsg)
					}
				}
			}
		}
	}
	return nil
}

func LoadConfig(configFile string) (*Config, error) {
	var c Config
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		log.Fatalf("Cannot load config file! Error: %s\n", err)
	}

	validator.SetValidationFunc("hasusers", validateUsers)
	isValid, errMap := validator.Validate(c)

	if !isValid {
		for k, v := range errMap {
			for _, err := range v {
				errorMsg := errorExpl[k][err]
				if strings.Contains(errorMsg, "%s") {
					errorMsg = fmt.Sprintf(errorExpl[k][err], k)
				}
				fmt.Println(errorMsg)
			}
		}
	}

	// TODO: Config validation, default values, etc

	return &c, nil
}
