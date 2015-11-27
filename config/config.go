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
	Users   map[string]User `validate:"nonzero,hasusers"`
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
	AuthInfo Auth               `toml:"auth"`
	Certs    map[string]Cert    `validate:"min=1"`
	Networks map[string]Network `validate:"validnetworks"`
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

type ConfigError struct {
	Field   string
	Message string
}

func (ce ConfigError) Error() string {
	return fmt.Sprintf("%s: %s", ce.Field, ce.Message)
}

type MultiError struct {
	Errors []error
}

func (me MultiError) Error() string {
	var errStr string
	for _, e := range me.Errors {
		errStr += ", " + e.Error()
	}
	errStr = strings.TrimLeft(errStr, ",")
	return fmt.Sprintf("%s", errStr)
}

type UserError struct {
	User    string
	Message string
}

func (ue UserError) Error() string {
	return fmt.Sprintf("[users.%s] -> %s", ue.User, ue.Message)
}

type UsersError struct {
	Field   string
	Message string
}

var errorExpl = map[string]map[error]string{
	"Logging.Adapter":  map[error]string{validator.ErrZeroValue: "An adapter is required. Valid Options: SQLite3"},
	"Logging.Database": map[error]string{validator.ErrZeroValue: "You must specify the name of the logging database."},
	"Nick":             map[error]string{validator.ErrZeroValue: "You must specify a nickname in order to connect to an IRC server.", validator.ErrMax: "Nickname can only be 9 characters long."},
	"AltNick":          map[error]string{validator.ErrZeroValue: "You must specify a alternate nickname in order to connect to an IRC server.", validator.ErrMax: "Altenate nickname can only be 9 characters long."},
	"Certs":            map[error]string{validator.ErrZeroValue: "You must specify at least one certificate in order to authenticate to zounce.", validator.ErrMin: "You must have at least 1 certificate on your user in order to authenticate."},
	"Users":            map[error]string{validator.ErrZeroValue: "You must specify at least one user in order to use to zounce."},
	"AuthInfo.CAPath":  map[error]string{validator.ErrZeroValue: "You must specify the CA for your certificate to verify."},
}

func validateNetworks(v interface{}, param string) error {
	return nil
}

func validateUsers(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	var mError MultiError
	if st.Kind() == reflect.Map {
		keys := st.MapKeys()
		for _, user := range keys {
			isValid, errMap := validator.Validate(st.MapIndex(user).Interface())
			if !isValid {
				for k, v := range errMap {
					for _, err := range v {
						errorMsg := errorExpl[k][err]
						if len(errorMsg) > 0 {
							mError.Errors = append(mError.Errors, &UserError{user.String(), errorMsg})
						}
					}
				}
			}
		}
	}
	return mError
}

func LoadConfig(configFile string) (*Config, []error) {
	var c Config
	_, err := toml.DecodeFile(configFile, &c)
	if err != nil {
		log.Fatalf("Cannot load config file! Error: %s\n", err)
	}

	var errs []error

	validator.SetValidationFunc("hasusers", validateUsers)
	validator.SetValidationFunc("validnetworks", validateNetworks)
	isValid, errMap := validator.Validate(c)

	if !isValid {
		for k, v := range errMap {
			for _, err := range v {
				switch reflect.TypeOf(err).String() {
				case "config.MultiError":
					errors := reflect.ValueOf(err).FieldByName("Errors")
					for i := 0; i < errors.Len(); i++ {
						ue := errors.Index(i).Interface().(*UserError)
						errs = append(errs, ue)
					}
					break
				case "*errors.errorString":
					kErr := errorExpl[k][err]
					if len(kErr) > 0 {
						errs = append(errs, &ConfigError{k, kErr})
					} else {
						errs = append(errs, &ConfigError{k, err.Error()})
					}
				default:
					fmt.Println("log this?")
					break
				}
			}
		}
	}

	// TODO: Config validation, default values, etc

	return &c, errs
}
