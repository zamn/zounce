package config

import (
	"errors"
	"os"
	"os/exec"
	"reflect"
	"sort"
	"testing"
)

var TemplateFile string
var EmptyFile string
var PartialFile string
var BadNetworkFile string

func setup() {
	TemplateFile = "config.toml"
	EmptyFile = "empty.toml"
	PartialFile = "partial.toml"
	BadNetworkFile = "badnetwork.toml"
}

func TestMain(m *testing.M) {
	setup()

	retCode := m.Run()

	os.Exit(retCode)
}

// Lets make sure I didn't break my config while developing, heh
// Requires tomlv
func TestValidTomlTemplate(t *testing.T) {
	cmd := exec.Command("tomlv", TemplateFile)

	out, err := cmd.CombinedOutput()
	if err != nil {
		t.Fatalf("Config file not valid TOML! %s Error: %s\n", TemplateFile, out)
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

func sameErrors(expected, err []error, t *testing.T) bool {
	if len(err) != len(expected) {
		t.Fatalf("Invalid number of errors returned for %s. Expected: %s, Got: %s", PartialFile, len(expected), len(err))
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

	expectedConfig := &Config{
		Title: "Zounce Configuration",
		Port:  7777,
		Logging: LogInfo{
			Adapter:  "SQLite3",
			Database: "zounce",
		},
		Users: map[string]User{
			"zamn": User{
				Nick:     "zamn",
				AltNick:  "zamn92",
				Username: "zamn",
				Realname: "Adam",
				AuthInfo: Auth{
					CAPath: "certs/ca.crt",
				},
				Certs: map[string]Cert{
					"desktop": Cert{
						Path: "certs/zamn.crt",
					},
				},
				Networks: map[string]Network{
					"GameSurge": Network{
						Name: "The GameSurge Network",
						Servers: []string{
							"irc.gamesurge.net:6666",
						},
						Password: "",
						PerformInfo: Perform{
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

	if !reflect.DeepEqual(c, expectedConfig) {
		t.Fatalf("Parsed config file and expected config file are not the same!")
	}
}

func TestNetworkErrors(t *testing.T) {
	_, err := LoadConfig(BadNetworkFile)

	if err == nil {
		t.Fatalf("Error not found in bad networks config.")
	}

	expected := []error{
		errors.New("[users.zamn] -> [networks.GameSurge] -> You must specify at least one server in order to use this network with zounce."),
	}

	sameErrors(expected, err, t)
}

func TestEmptyFileErrors(t *testing.T) {
	_, err := LoadConfig(EmptyFile)

	if len(err) == 0 {
		t.Fatalf("Validated an empty configuration file!\n")
	}

	expected := []error{
		errors.New("Users: You must specify at least one user in order to use to zounce."),
		errors.New("Logging.Database: You must specify the name of the logging database."),
		errors.New("Logging.Adapter: An adapter is required. Valid Options: SQLite3"),
	}

	sameErrors(expected, err, t)
}

func TestPartialFileErrors(t *testing.T) {
	_, err := LoadConfig(PartialFile)

	if len(err) == 0 {
		t.Fatalf("Validated a config with errors.")
	}

	expected := []error{
		errors.New("Logging.Adapter: An adapter is required. Valid Options: SQLite3"),
		errors.New("Logging.Database: You must specify the name of the logging database."),
		errors.New("[users.zamn] -> You must specify a nickname in order to connect to an IRC server."),
		errors.New("[users.zamn] -> You must specify a alternate nickname in order to connect to an IRC server."), errors.New("[users.zamn] -> You must specify the CA for your certificate to verify."), errors.New("[users.zamn] -> You must have at least 1 certificate on your user in order to authenticate.")}

	sameErrors(expected, err, t)
}
