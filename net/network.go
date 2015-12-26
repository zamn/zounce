package net

import (
	"reflect"

	"github.com/zamN/zounce/config/confutils"
	"github.com/zamN/zounce/net/perform"
)

func ValidateNetworks(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	if st.Kind() == reflect.Map {
		netErrors := &confutils.ConfigError{
			Type: confutils.NetworkType,
			Id:   "Networks",
		}

		return confutils.ValidateMap(netErrors, st)
	}

	return nil
}

type Network struct {
	Name        string   `validate:"nonzero"`
	Servers     []string `validate:"min=1"`
	Password    string
	PerformInfo perform.Perform `toml:"perform"`
}
