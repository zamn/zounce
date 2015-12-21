package net

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/validator.v2"

	"github.com/zamN/zounce/config/confutils"
	"github.com/zamN/zounce/net/perform"
)

type NetworkError struct {
	Network string
	Message string
}

func (ne NetworkError) Error() string {
	return fmt.Sprintf("[networks.%s] -> %s", ne.Network, ne.Message)
}

func (ne NetworkError) FormatErrors() []error {
	return []error{errors.New(fmt.Sprintf("[networks.%s] -> %s", ne.Network, ne.Message))}
}

func ValidateNetworks(v interface{}, param string) error {

	st := reflect.ValueOf(v)
	mError := &confutils.MultiError{}
	if st.Kind() == reflect.Map {
		keys := st.MapKeys()
		for _, server := range keys {
			errMap := validator.Validate(st.MapIndex(server).Interface()).(validator.ErrorMap)
			if errMap != nil {
				for k, v := range errMap {
					for _, err := range v {
						errorMsg, ok := confutils.GetErrExpln(k, err)
						if ok {
							mError.Add(NetworkError{server.String(), errorMsg})
						} else {
							mError.Add(NetworkError{server.String(), fmt.Sprintf("Unknown error: %s", err)})
						}
					}
				}
			}
		}
	}

	if !mError.Empty() {
		return mError
	}
	return nil
}

type Network struct {
	Name        string   `validate:"nonzero"`
	Servers     []string `validate:"min=1"`
	Password    string
	PerformInfo perform.Perform `toml:"perform"`
}
