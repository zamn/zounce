package user

import (
	"errors"
	"fmt"
	"reflect"

	"gopkg.in/validator.v2"

	"github.com/zamN/zounce/config/confutils"
	"github.com/zamN/zounce/logging"
	"github.com/zamN/zounce/net"
	"github.com/zamN/zounce/user/cert"
)

type UserError struct {
	User   string
	Errors []error
}

func (ue *UserError) Add(err error) {
	ue.Errors = append(ue.Errors, err)
}

func (ue *UserError) Empty() bool {
	return len(ue.Errors) == 0
}

func (ue UserError) Error() string {
	errOut := ""
	errs := ue.FormatErrors()
	for _, e := range errs {
		errOut += e.Error() + "\n"
	}
	return errOut
}

func (ue UserError) FormatErrors() []error {
	var final []error
	for _, e := range ue.Errors {
		switch reflect.TypeOf(e).String() {
		case "*confutils.MultiError":
			temp := reflect.ValueOf(e).Interface().(*confutils.MultiError)
			for _, err := range temp.Errors {
				final = append(final, errors.New(fmt.Sprintf("[users.%s]%s", ue.User, err)))
			}
			break
		case "*net.NetworkError":
			final = append(final, errors.New(fmt.Sprintf("[users.%s]%s", ue.User, e)))
			break
		default:
			final = append(final, errors.New(fmt.Sprintf("[users.%s] -> %s", ue.User, e)))
			break
		}
	}
	return final
}

func ValidateUsers(v interface{}, param string) error {
	st := reflect.ValueOf(v)

	multiUserErr := &confutils.MultiError{}

	if st.Kind() == reflect.Map {
		keys := st.MapKeys()
		for _, user := range keys {
			userError := &UserError{User: user.String()}

			errMap := validator.Validate(st.MapIndex(user).Interface()).(validator.ErrorMap)
			if errMap != nil {
				for k, v := range errMap {
					for _, err := range v {
						// TODO: Change structure to if errorExpl[k] has key err, then add errorMsg
						errorMsg, ok := confutils.GetErrExpln(k, err)

						// If this is a top level error
						if ok {
							userError.Add(errors.New(errorMsg))
						} else {
							userError.Add(err)
						}
					}
				}
				multiUserErr.Add(userError)
				fmt.Println(userError.Errors)
				fmt.Println("MUE", multiUserErr)
			}
		}
	}

	if !multiUserErr.Empty() {
		return multiUserErr
	}
	return nil
}

type User struct {
	Nick     string `validate:"nonzero,max=9"`
	AltNick  string `validate:"nonzero,max=9"`
	Username string
	Realname string
	Logging  logging.LogInfo
	Certs    map[string]cert.Cert   `validate:"nonzero,validcerts"`
	Networks map[string]net.Network `validate:"validnetworks"`
}
