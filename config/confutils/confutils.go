package confutils

// Config Utility functions
import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/validator.v2"
)

type ErrorType int

// Only fields in the User struct which have custom validations
const (
	BaseType ErrorType = iota
	UserType
	NetworkType
	CertType
)

func (et ErrorType) String() string {
	switch et {
	case BaseType:
		return "BaseType"
		break
	case UserType:
		return "UserType"
		break
	case NetworkType:
		return "NetworkType"
		break
	case CertType:
		return "CertType"
		break
	}

	return "Unknown ErrorType"
}

type ConfigError struct {
	Type   ErrorType
	Id     string
	Errors []error
}

func (ce *ConfigError) Error() string {
	var tag string
	switch ce.Type {
	case BaseType:
		if len(ce.Id) > 0 {
			tag = fmt.Sprintf("[%s]", ce.Id)
		} else {
			tag = "[ERROR]"
		}
		break
	case UserType:
		tag = fmt.Sprintf("[users.%s]", ce.Id)
		break
	case NetworkType:
		tag = fmt.Sprintf("[networks.%s]", ce.Id)
		break
	case CertType:
		tag = fmt.Sprintf("[certs.%s]", ce.Id)
		break
	default:
		tag = fmt.Sprintf("[unknown.%s]", ce.Id)
		break
	}

	if len(ce.Errors) > 0 {
		if len(ce.Errors) == 1 {
			return fmt.Sprintf("%s %s", tag, ce.Errors[0])
		}
		return fmt.Sprintf("%s %s", tag, "There are multiple errors for this block.")
	}

	return fmt.Sprintf("%s %s", tag, "There are no errors for this block.")
}

// A ConfigError is only a container if it contains
// One or more ConfigError
func (ce *ConfigError) IsContainer() bool {
	for _, e := range ce.Errors {
		if _, ok := e.(*ConfigError); ok {
			return true
		}
	}
	return false
}

func (ce *ConfigError) IsEmpty() bool {
	return len(ce.Errors) == 0
}

type TextConvert map[error]string

var errExpl = map[ErrorType]map[string]TextConvert{
	BaseType: map[string]TextConvert{
		"Title": TextConvert{
			validator.ErrZeroValue: "Title not supplied, using default title 'Zounce Config'",
		},
		"Port": TextConvert{
			validator.ErrZeroValue: "Port not supplied, using default port",
			validator.ErrMin:       "Port must be greater than 0, using default port",
		},
		"CAPath": TextConvert{
			validator.ErrZeroValue: "You must specify the CA for your user certificates to validate against.",
		},
	},
	UserType: map[string]TextConvert{
		"Users": TextConvert{
			validator.ErrZeroValue: "You must specify at least one user in order to use to zounce.",
		},
		"Logging.Adapter": TextConvert{
			validator.ErrZeroValue: "An adapter is required. Valid Options: SQLite3, Flatfile",
		},
		"Logging.Database": TextConvert{
			validator.ErrZeroValue: "You must specify the name of the logging database.",
		},
		"Nick": TextConvert{
			validator.ErrZeroValue: "You must specify a nickname in order to connect to an IRC server.",
			validator.ErrMax:       "Nickname can only be 9 characters long.",
		},
		"AltNick": TextConvert{
			validator.ErrZeroValue: "You must specify a alternate nickname in order to connect to an IRC server.",
			validator.ErrMax:       "Altenate nickname can only be 9 characters long.",
		},
	},
	NetworkType: map[string]TextConvert{
		"Networks": TextConvert{
			validator.ErrZeroValue: "You must specify at least one network in order to use to zounce.",
		},
		"Servers": TextConvert{
			validator.ErrMin: "You must specify at least one server in order to use this network with zounce.",
		},
		"Name": TextConvert{
			validator.ErrZeroValue: "You must specify a name for this network!",
		},
	},
	CertType: map[string]TextConvert{
		"Certs": TextConvert{
			validator.ErrZeroValue: "You must specify at least one certificate in order to authenticate to zounce.",
		},
	},
}

func GetErrExpln(eType ErrorType, field string, err error) (string, bool) {
	expln, ok := errExpl[eType][field][err]
	return expln, ok
}

func ValidateMap(container *ConfigError, segMap reflect.Value) error {
	segKeys := segMap.MapKeys()

	// Since doing validation ourselves, make sure we *have* something
	if len(segKeys) == 0 {
		errStr, _ := GetErrExpln(container.Type, container.Id, validator.ErrZeroValue)
		return &ConfigError{
			Type: container.Type,
			Id:   container.Id,
			Errors: []error{
				errors.New(errStr),
			},
		}
	}

	for _, sk := range segKeys {
		segError := &ConfigError{
			Type: container.Type,
			Id:   sk.String(),
		}

		errMap := validator.Validate(segMap.MapIndex(sk).Interface())
		if errMap != nil {
			errMap := errMap.(validator.ErrorMap)

			for field, errArr := range errMap {
				for _, err := range errArr {
					errStr, ok := GetErrExpln(segError.Type, field, err)
					if ok {
						ce := &ConfigError{
							Type: container.Type,
							Id:   field,
							Errors: []error{
								errors.New(errStr),
							},
						}
						segError.Errors = append(segError.Errors, ce)
					} else {
						segError.Errors = append(segError.Errors, err)
					}
				}
			}
			container.Errors = append(container.Errors, segError)
		}
	}

	if !container.IsEmpty() {
		return container
	}

	return nil
}
