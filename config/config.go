package config

import (
	"fmt"
	"log"
	"reflect"

	"gopkg.in/validator.v2"

	"github.com/BurntSushi/toml"
	"github.com/zamN/zounce/config/confutils"
	"github.com/zamN/zounce/net"
	"github.com/zamN/zounce/user"
	"github.com/zamN/zounce/user/cert"
)

type Config struct {
	Title  string               `validate:"nonzero"`
	Port   int                  `validate:"nonzero"`
	CAPath string               `toml:"ca_path" validate:"nonzero"`
	Users  map[string]user.User `validate:"nonzero,validusers"`
}

type ConfigError struct {
	Field   string
	Message string
}

func (ce ConfigError) Error() string {
	return fmt.Sprintf("%s: %s", ce.Field, ce.Message)
}

func LoadConfig(configFile string) (*Config, []error) {
	var c Config
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		log.Fatalf("Cannot load config file! Error: %s\n", err)
	}

	var errs []error

	validator.SetValidationFunc("validusers", user.ValidateUsers)
	validator.SetValidationFunc("validnetworks", net.ValidateNetworks)
	validator.SetValidationFunc("validcerts", cert.ValidateCerts)
	errMap := validator.Validate(c)

	if errMap != nil {
		errMap = errMap.(validator.ErrorMap)
		fmt.Printf("%#v\n", errMap)
		for k, v := range errMap.(validator.ErrorMap) {
			for _, err := range v {
				switch reflect.TypeOf(err).String() {
				// For dealing with sub-errors within config segments
				case "*confutils.MultiError":
					temp := reflect.ValueOf(err).Interface().(*confutils.MultiError)
					for i := 0; i < len(temp.Errors); i++ {
						errs = append(errs, temp.Errors[i].FormatErrors()...)
					}
					break
				default:
					kErr, ok := confutils.GetErrExpln(k, err)
					if ok {
						errs = append(errs, &ConfigError{k, kErr})
					} else {
						errs = append(errs, &ConfigError{k, err.Error()})
					}
					fmt.Println(reflect.TypeOf(k))
					fmt.Println(k)
					break
				}
			}
		}
	}

	// TODO: Config validation, default values, etc

	return &c, errs
}
