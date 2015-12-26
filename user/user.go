package user

import (
	"reflect"

	"github.com/zamN/zounce/config/confutils"
	"github.com/zamN/zounce/logging"
	"github.com/zamN/zounce/net"
	"github.com/zamN/zounce/user/cert"
)

func ValidateUsers(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	if st.Kind() == reflect.Map {
		userErrors := &confutils.ConfigError{
			Type: confutils.UserType,
			Id:   "Users",
		}

		return confutils.ValidateMap(userErrors, st)
	}

	return nil
}

type User struct {
	Nick     string `validate:"nonzero,max=9"`
	AltNick  string `validate:"nonzero,max=9"`
	Username string
	Realname string
	Logging  logging.LogInfo
	Certs    map[string]cert.Cert   `validate:"validcerts"`
	Networks map[string]net.Network `validate:"validnetworks"`
}
