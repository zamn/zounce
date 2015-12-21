package config

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
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

func contains(expected, err []error) bool {
	var found bool
	for _, ex := range expected {
		found = false
		for _, e := range err {
			if ex.Error() == e.Error() {
				found = true
				// figure out how to throw !found in that loop...
				break
			}
		}
	}

	return found
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

	servMinErr, _ := confutils.GetErrExpln("Servers", validator.ErrMin)
	nameEmptyErr, _ := confutils.GetErrExpln("Name", validator.ErrZeroValue)
	caPathEmptyErr, _ := confutils.GetErrExpln("CAPath", validator.ErrZeroValue)

	expUserErrors := user.UserError{
		User: BaseUser,
		Errors: []error{
			&net.NetworkError{GSNet, servMinErr},
			&net.NetworkError{GSNet, nameEmptyErr},
		},
	}
	expected := []error{
		&ConfigError{"CAPath", caPathEmptyErr},
	}

	expected = append(expected, expUserErrors.FormatErrors()...)

	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s. Expected: %d, Got: %d\n", PartialFile, len(expected), len(err))
	}

	if !contains(expected, err) {
		t.Fatalf("Returned errors not equal to expected errors. Expected: %v, Got: %v\n", expected, err)
	}
}

func TestEmptyFileErrors(t *testing.T) {
	_, err := LoadConfig(EmptyFile)

	if len(err) == 0 {
		t.Fatalf("Validated an empty configuration file!\n")
	}

	usersEmptyErr, _ := confutils.GetErrExpln("Users", validator.ErrZeroValue)
	caPathEmptyErr, _ := confutils.GetErrExpln("CAPath", validator.ErrZeroValue)
	titleEmptyErr, _ := confutils.GetErrExpln("Title", validator.ErrZeroValue)
	portEmptyErr, _ := confutils.GetErrExpln("Port", validator.ErrZeroValue)
	expected := []error{
		&ConfigError{"Users", usersEmptyErr},
		&ConfigError{"CAPath", caPathEmptyErr},
		&ConfigError{"Title", titleEmptyErr},
		&ConfigError{"Port", portEmptyErr},
	}

	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s. Expected: %d, Got: %d\n", PartialFile, len(expected), len(err))
	}

	if !contains(expected, err) {
		t.Fatalf("Returned errors not equal to expected errors. Expected: %v, Got: %v\n", expected, err)
	}
}

func TestPartialFileErrors(t *testing.T) {
	_, err := LoadConfig(PartialFile)

	if len(err) == 0 {
		t.Fatalf("Validated a config with errors.")
	}
	fmt.Println("ERR", err)
	adapterEmptyErr, _ := confutils.GetErrExpln("Logging.Adapter", validator.ErrZeroValue)
	dbEmptyErr, _ := confutils.GetErrExpln("Logging.Database", validator.ErrZeroValue)
	nickEmptyErr, _ := confutils.GetErrExpln("Nick", validator.ErrZeroValue)
	altNickEmptyErr, _ := confutils.GetErrExpln("AltNick", validator.ErrZeroValue)
	certsEmptyErr, _ := confutils.GetErrExpln("Certs", validator.ErrZeroValue)
	caPathEmptyErr, _ := confutils.GetErrExpln("Certs", validator.ErrZeroValue)

	expUserErrors := user.UserError{
		User: BaseUser,
		Errors: []error{
			errors.New(adapterEmptyErr),
			errors.New(dbEmptyErr),
			errors.New(nickEmptyErr),
			errors.New(altNickEmptyErr),
			errors.New(certsEmptyErr),
		},
	}
	expected := []error{
		&ConfigError{"CAPath", caPathEmptyErr},
	}

	expected = append(expected, expUserErrors.FormatErrors()...)

	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s. Expected: %d, Got: %d\n", PartialFile, len(expected), len(err))
	}

	if !contains(expected, err) {
		t.Fatalf("Returned errors not equal to expected errors. Expected: %v, Got: %v\n", expected, err)
	}
}
