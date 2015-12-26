package config

import (
	"errors"
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
	Title  string               `validate:"validbase=Title"`
	Port   int                  `validate:"validbase=Port"`
	CAPath string               `toml:"ca_path" validate:"validbase=CAPath"`
	Users  map[string]user.User `validate:"validusers"`
}

func ValidateConfigBase(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	if st.Interface() == reflect.Zero(st.Type()).Interface() {
		ce := &confutils.ConfigError{
			Type: confutils.BaseType,
			Id:   param,
		}
		errStr, _ := confutils.GetErrExpln(ce.Type, ce.Id, validator.ErrZeroValue)
		ce.Errors = []error{errors.New(errStr)}
		return ce
	}

	if st.Kind() == reflect.Int {
		if reflect.ValueOf(st.Interface()).Int() < 0 {
			ce := &confutils.ConfigError{
				Type: confutils.BaseType,
				Id:   param,
			}
			errStr, _ := confutils.GetErrExpln(ce.Type, ce.Id, validator.ErrMin)
			ce.Errors = []error{errors.New(errStr)}
			return ce
		}
	}

	return nil
}

func LoadConfig(configFile string) (*Config, []error) {
	var c Config
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		log.Fatalf("Cannot load config file! Error: %s\n", err)
	}

	var errs []error

	validator.SetValidationFunc("validnetworks", net.ValidateNetworks)
	validator.SetValidationFunc("validcerts", cert.ValidateCerts)
	validator.SetValidationFunc("validusers", user.ValidateUsers)
	validator.SetValidationFunc("validbase", ValidateConfigBase)
	errMap := validator.Validate(c)

	if errMap != nil {
		errMap := errMap.(validator.ErrorMap)

		// TODO: Some way to display warnings lol.
		warnings := map[string](func()){
			"Title": (func() {
				c.Title = "Zounce Configuration"
			}),
			"Port": (func() {
				c.Port = 7777
			}),
		}

		for _, errArr := range errMap {
			for _, err := range errArr {
				confErr, ok := err.(*confutils.ConfigError)
				if ok {
					defFunc, ok := warnings[confErr.Id]
					if ok {
						defFunc()
					} else {
						errs = append(errs, err)
					}
				}
			}
		}
	}

	// TODO: Handle defaults

	return &c, errs
}
