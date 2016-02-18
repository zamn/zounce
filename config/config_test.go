package config

import (
	"errors"
	"os"
	"reflect"
	"testing"

	"gopkg.in/validator.v2"

	"github.com/zamN/zounce/config/confutils"
	"github.com/zamN/zounce/logging"
	"github.com/zamN/zounce/net"
	"github.com/zamN/zounce/net/perform"
	"github.com/zamN/zounce/user"
	"github.com/zamN/zounce/user/cert"
)

var (
	DataDir        = "data/"
	TemplateFile   = DataDir + "config.toml"
	EmptyFile      = DataDir + "empty.toml"
	PartialFile    = DataDir + "partial.toml"
	BadNetworkFile = DataDir + "badnetwork.toml"
	BaseUser       = "zamn"
	GSNet          = "GameSurge"

	// Config Info all data files will be based off of
	// If you want to add a test, make sure the data (even if partially filled)
	// in it matches this struct
	BaseConfig = &Config{
		Title:  "Zounce Configuration",
		Port:   7777,
		CAPath: "certs/ca.crt",
		Users: map[string]user.User{
			BaseUser: user.User{
				Nick:     "zamn",
				AltNick:  "zamn92",
				Username: "zamn",
				Realname: "Adam",
				Logging: logging.LogInfo{
					Adapter:  "SQLite3",
					Database: "zounce",
				},
				Certs: map[string]cert.Cert{
					"desktop": cert.Cert{
						Path: "certs/zamn.crt",
					},
				},
				Networks: map[string]net.Network{
					GSNet: net.Network{
						Name: "The GameSurge Network",
						Servers: []string{
							"irc.gamesurge.net:6666",
						},
						Password: "",
						PerformInfo: perform.Perform{
							Channels: []string{
								"#zamN",
							},
							Commands: []string{
								"/msg AuthServ@Services.Gamesurge.net auth user pass",
							},
						},
					},
				},
			},
		},
	}
)

func TestMain(m *testing.M) {
	// os.Exit() does not respect defer statements
	ret := m.Run()
	os.Exit(ret)
}

func equals(exConfErr, confErr *confutils.ConfigError) bool {
	if exConfErr.Type == confErr.Type {
		if exConfErr.Id == confErr.Id {
			if len(exConfErr.Errors) == len(confErr.Errors) {
				return containsAll(exConfErr.Errors, confErr.Errors)
			}
		}
	}
	return false
}

// Probably will need to be improved in the future
func containsAll(exInput, err []error) bool {
	// This is probably horribly slow
	var expected []error
	expected = append(expected, exInput...)

	for _, e := range err {
		confErr, ok := e.(*confutils.ConfigError)
		for i, ex := range expected {
			if !ok {
				if e.Error() == ex.Error() {
					expected = append(expected[:i], expected[i+1:]...)
					break
				}
			} else {
				exConfErr, ok := ex.(*confutils.ConfigError)
				if ok {
					if equals(exConfErr, confErr) {
						expected = append(expected[:i], expected[i+1:]...)
						break
					}
				}
			}
		}
	}

	return len(expected) == 0
}

func TestValidConfig(t *testing.T) {
	c, err := LoadConfig(TemplateFile)

	if err != nil {
		t.Fatalf("Error in valid config file! Messages: %s\n", err)
	}

	if !reflect.DeepEqual(c, BaseConfig) {
		t.Fatalf("Parsed config file and expected config file are not the same!")
	}
}

func TestNetworkErrors(t *testing.T) {
	_, err := LoadConfig(BadNetworkFile)

	if err == nil {
		t.Fatalf("Error(s) not found in bad networks config.")
	}

	servMinErr, _ := confutils.GetErrExpln(confutils.NetworkType, "Servers", validator.ErrMin)
	nameEmptyErr, _ := confutils.GetErrExpln(confutils.NetworkType, "Name", validator.ErrZeroValue)
	caPathEmptyErr, _ := confutils.GetErrExpln(confutils.BaseType, "CAPath", validator.ErrZeroValue)

	expUserErrors := &confutils.ConfigError{
		Type: confutils.UserType,
		Id:   "Users",
		Errors: []error{
			&confutils.ConfigError{
				Type: confutils.UserType,
				Id:   "zamn",
				Errors: []error{
					&confutils.ConfigError{
						Type: confutils.NetworkType,
						Id:   "Networks",
						Errors: []error{
							&confutils.ConfigError{
								Type: confutils.NetworkType,
								Id:   "GameSurge",
								Errors: []error{
									&confutils.ConfigError{
										Type: confutils.NetworkType,
										Id:   "Name",
										Errors: []error{
											errors.New(nameEmptyErr),
										},
									},
									&confutils.ConfigError{
										Type: confutils.NetworkType,
										Id:   "Servers",
										Errors: []error{
											errors.New(servMinErr),
										},
									},
								},
							},
						},
					},
				},
			},
		},
	}
	expected := []error{
		&confutils.ConfigError{
			Type: confutils.BaseType,
			Id:   "CAPath",
			Errors: []error{
				errors.New(caPathEmptyErr),
			},
		},
		expUserErrors,
	}

	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s.\nExpected: %d\nGot: %d\n", PartialFile, len(expected), len(err))
	}

	if !containsAll(expected, err) {
		t.Fatalf("Returned errors not equal to expected errors.\nExpected: %v\nGot: %v\n", expected, err)
	}
}

func TestEmptyFileErrors(t *testing.T) {
	_, err := LoadConfig(EmptyFile)

	if len(err) == 0 {
		t.Fatalf("Validated an empty configuration file!\n")
	}

	usersEmptyErr, _ := confutils.GetErrExpln(confutils.UserType, "Users", validator.ErrZeroValue)
	caPathEmptyErr, _ := confutils.GetErrExpln(confutils.BaseType, "CAPath", validator.ErrZeroValue)

	expected := []error{
		&confutils.ConfigError{
			Type: confutils.UserType,
			Id:   "Users",
			Errors: []error{
				errors.New(usersEmptyErr),
			},
		},
		&confutils.ConfigError{
			Type: confutils.BaseType,
			Id:   "CAPath",
			Errors: []error{
				errors.New(caPathEmptyErr),
			},
		},
	}

	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s.\nExpected: %d\nGot: %d\n", PartialFile, len(expected), len(err))
	}

	if !containsAll(expected, err) {
		t.Fatalf("Returned errors not equal to expected errors.\nExpected: %v\nGot: %v\n", expected, err)
	}
}

func TestPartialFileErrors(t *testing.T) {
	_, err := LoadConfig(PartialFile)

	if len(err) == 0 {
		t.Fatalf("Validated a config with errors.")
	}

	adapterEmptyErr, _ := confutils.GetErrExpln(confutils.UserType, "Logging.Adapter", validator.ErrZeroValue)
	dbEmptyErr, _ := confutils.GetErrExpln(confutils.UserType, "Logging.Database", validator.ErrZeroValue)
	nickEmptyErr, _ := confutils.GetErrExpln(confutils.UserType, "Nick", validator.ErrZeroValue)
	altNickEmptyErr, _ := confutils.GetErrExpln(confutils.UserType, "AltNick", validator.ErrZeroValue)
	certsEmptyErr, _ := confutils.GetErrExpln(confutils.CertType, "Certs", validator.ErrZeroValue)
	caPathEmptyErr, _ := confutils.GetErrExpln(confutils.BaseType, "CAPath", validator.ErrZeroValue)
	netBlockEmptyErr, _ := confutils.GetErrExpln(confutils.NetworkType, "Networks", validator.ErrZeroValue)

	// TODO: Create generators(?) for these monsters.
	expUserErrors := &confutils.ConfigError{
		Type: confutils.UserType,
		Id:   "Users",
		Errors: []error{
			&confutils.ConfigError{
				Type: confutils.UserType,
				Id:   "zamn",
				Errors: []error{
					&confutils.ConfigError{
						Type: confutils.NetworkType,
						Id:   "Networks",
						Errors: []error{
							errors.New(netBlockEmptyErr),
						},
					},
					&confutils.ConfigError{
						Type: confutils.CertType,
						Id:   "Certs",
						Errors: []error{
							errors.New(certsEmptyErr),
						},
					},
					&confutils.ConfigError{
						Type: confutils.UserType,
						Id:   "Logging.Adapter",
						Errors: []error{
							errors.New(adapterEmptyErr),
						},
					},
					&confutils.ConfigError{
						Type: confutils.UserType,
						Id:   "Logging.Database",
						Errors: []error{
							errors.New(dbEmptyErr),
						},
					},
					&confutils.ConfigError{
						Type: confutils.UserType,
						Id:   "Nick",
						Errors: []error{
							errors.New(nickEmptyErr),
						},
					},
					&confutils.ConfigError{
						Type: confutils.UserType,
						Id:   "AltNick",
						Errors: []error{
							errors.New(altNickEmptyErr),
						},
					},
				},
			},
		},
	}

	expected := []error{
		&confutils.ConfigError{
			Type: confutils.BaseType,
			Id:   "CAPath",
			Errors: []error{
				errors.New(caPathEmptyErr),
			},
		},
		expUserErrors,
	}

	// Although this is checking at the baseline level, good enough
	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s.\nExpected: %d\nGot: %d\n", PartialFile, len(expected), len(err))
	}

	if !containsAll(expected, err) {
		t.Fatalf("Returned errors not equal to expected errors.\nExpected: %v\nGot: %v\n", expected, err)
	}
}
