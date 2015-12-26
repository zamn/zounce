package cert

import (
	"reflect"

	"github.com/zamN/zounce/config/confutils"
)

func ValidateCerts(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	if st.Kind() == reflect.Map {
		certErrors := &confutils.ConfigError{
			Type: confutils.CertType,
			Id:   "Certs",
		}

		return confutils.ValidateMap(certErrors, st)
	}

	return nil
}

type Cert struct {
	Path string `toml:"cert_path" validate:"nonzero"`
}
