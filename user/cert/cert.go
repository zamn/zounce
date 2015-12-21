package cert

import (
	"errors"
	"fmt"
	"reflect"

	"github.com/zamN/zounce/config/confutils"
	"gopkg.in/validator.v2"
)

type CertError struct {
	Name    string
	Message string
}

func (ce CertError) Error() string {
	return fmt.Sprintf("[certs.%s] -> %s", ce.Name, ce.Message)
}

func (ce CertError) FormatErrors() []error {
	return []error{errors.New(fmt.Sprintf("[certs.%s] -> %s", ce.Name, ce.Message))}
}

func ValidateCerts(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	multiUserErr := &confutils.MultiError{}

	if st.Kind() == reflect.Map {
		keys := st.MapKeys()
		for _, certName := range keys {
			certError := &CertError{Name: certName.String()}

			errMap := validator.Validate(st.MapIndex(certName).Interface()).(validator.ErrorMap)
			if errMap != nil {
				for k, v := range errMap {
					for _, err := range v {
						// TODO: Change structure to if errorExpl[k] has key err, then add errorMsg
						errorMsg, ok := confutils.GetErrExpln(k, err)

						// If this is a top level error
						if ok {
							certError.Message = errorMsg
						} else {
							certError.Message = err.Error()
						}
					}
				}
				multiUserErr.Add(certError)
			}
		}
	}

	if !multiUserErr.Empty() {
		return multiUserErr
	}
	return nil
}

type Cert struct {
	Path string `toml:"cert_path" validate:"nonzero"`
}
