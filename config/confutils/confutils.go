package confutils

// Config Utility functions
import (
	"fmt"
	"strings"

	"gopkg.in/validator.v2"
)

var errorExpl = map[string]map[error]string{
	// TODO: Don't hardcode adapter 'valid options', also this is pretty ugly
	"Logging.Adapter":  map[error]string{validator.ErrZeroValue: "An adapter is required. Valid Options: SQLite3, Flatfile"},
	"Logging.Database": map[error]string{validator.ErrZeroValue: "You must specify the name of the logging database."},
	"Nick":             map[error]string{validator.ErrZeroValue: "You must specify a nickname in order to connect to an IRC server.", validator.ErrMax: "Nickname can only be 9 characters long."},
	"AltNick":          map[error]string{validator.ErrZeroValue: "You must specify a alternate nickname in order to connect to an IRC server.", validator.ErrMax: "Altenate nickname can only be 9 characters long."},
	"Certs":            map[error]string{validator.ErrZeroValue: "You must specify at least one certificate in order to authenticate to zounce."},
	"Users":            map[error]string{validator.ErrZeroValue: "You must specify at least one user in order to use to zounce."},
	"Servers":          map[error]string{validator.ErrMin: "You must specify at least one server in order to use this network with zounce."},
	"Name":             map[error]string{validator.ErrZeroValue: "You must specify a name for this network!"},
	"CAPath":           map[error]string{validator.ErrZeroValue: "You must specify the CA for your user certificates to validate against."},
	"Title":            map[error]string{validator.ErrZeroValue: "Title not supplied, using default title 'Zounce Config'"},
	"Port":             map[error]string{validator.ErrZeroValue: "Port not supplied, using default port 1337"},
}

func GetErrExpln(field string, err error) (string, bool) {
	expln, ok := errorExpl[field][err]
	return expln, ok
}

type ProperError interface {
	Error() string
	FormatErrors() []error
}

type MultiError struct {
	Errors []ProperError
}

func (me *MultiError) Add(err ProperError) {
	me.Errors = append(me.Errors, err)
}

func (me *MultiError) Empty() bool {
	return len(me.Errors) == 0
}

func (me MultiError) Error() string {
	var errStr string
	for _, e := range me.Errors {
		errStr += ", " + e.Error()
	}
	errStr = strings.TrimLeft(errStr, ",")
	return fmt.Sprintf("%s", errStr)
}
