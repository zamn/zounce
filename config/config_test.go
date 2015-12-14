package config

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"testing"

	"github.com/go-validator/validator"
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
					"GameSurge": net.Network{
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
								"/msg AuthServ@Services.Gamesurge.net user pass",
							},
						},
					},
				},
			},
		},
	}
)

func TestMain(m *testing.M) {
	retCode := m.Run()

	os.Exit(retCode)
}

// Lets make sure I didn't break my config while developing, heh
// Requires tomlv
func TestValidTomlTemplate(t *testing.T) {
	cmd := exec.Command("tomlv", TemplateFile)

	_, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Config file not valid TOML! %s Error: %s\n", TemplateFile, err)
	}
}

type ByErrorFunc []error

func (bef ByErrorFunc) Less(i, j int) bool {
	return bef[i].Error() < bef[j].Error()
}

func (bef ByErrorFunc) Swap(i, j int) {
	bef[i], bef[j] = bef[j], bef[i]
}

func (bef ByErrorFunc) Len() int {
	return len(bef)
}

func sameErrors(expected, err []error, dataFile string, t *testing.T) bool {
	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s. Expected: %d, Got: %d\n", dataFile, len(expected), len(err))
	}

	sort.Sort(ByErrorFunc(expected))
	sort.Sort(ByErrorFunc(err))

	for i := 0; i < len(expected); i++ {
		if err[i].Error() != expected[i].Error() {
			t.Fatalf("Expected: \"%s\" and got: \"%s\"", expected[i], err[i])
		}
	}

	return true
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

	expUserErrors := UserError{
		User: BaseUser,
		Errors: []error{
			&NetworkError{GSNet, errorExpl["Servers"][validator.ErrMin]},
			&NetworkError{GSNet, errorExpl["Name"][validator.ErrZeroValue]},
		},
	}
	expected := []error{
		&ConfigError{"CAPath", errorExpl["CAPath"][validator.ErrZeroValue]},
	}

	expected = append(expected, expUserErrors.FormatErrors()...)

	sameErrors(expected, err, BadNetworkFile, t)
}

func TestEmptyFileErrors(t *testing.T) {
	_, err := LoadConfig(EmptyFile)

	if len(err) == 0 {
		t.Fatalf("Validated an empty configuration file!\n")
	}

	expected := []error{
		&ConfigError{"Users", errorExpl["Users"][validator.ErrZeroValue]},
		&ConfigError{"CAPath", errorExpl["CAPath"][validator.ErrZeroValue]},
	}

	sameErrors(expected, err, EmptyFile, t)
}

func TestPartialFileErrors(t *testing.T) {
	_, err := LoadConfig(PartialFile)

	if len(err) == 0 {
		t.Fatalf("Validated a config with errors.")
	}

	expUserErrors := UserError{
		User: BaseUser,
		Errors: []error{
			errors.New(errorExpl["Logging.Adapter"][validator.ErrZeroValue]),
			errors.New(errorExpl["Logging.Database"][validator.ErrZeroValue]),
			errors.New(errorExpl["Nick"][validator.ErrZeroValue]),
			errors.New(errorExpl["AltNick"][validator.ErrZeroValue]),
			errors.New(errorExpl["Certs"][validator.ErrZeroValue]),
		},
	}
	expected := []error{
		&ConfigError{"CAPath", errorExpl["CAPath"][validator.ErrZeroValue]},
	}

	expected = append(expected, expUserErrors.FormatErrors()...)

	sameErrors(expected, err, PartialFile, t)
}
