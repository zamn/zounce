package config

import (
	"errors"
	"log"
	"os"
	"os/exec"
	"sort"
	"testing"
)

var TemplateFile string
var EmptyFile string
var PartialFile string

func setup() {
	TemplateFile = "config.toml"
	EmptyFile = "empty.toml"
	PartialFile = "partial.toml"
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

func TestEmptyFileSafeDefaults(t *testing.T) {
	_, err := LoadConfig(EmptyFile)

	if len(err) == 0 {
		log.Fatalf("Validated an empty configuration file!\n")
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

func TestPartialFileSafeDefaults(t *testing.T) {
	_, err := LoadConfig(PartialFile)

	if len(err) == 0 {
		log.Fatalf("Validated a config with errors.")
	}

	expected := []error{
		errors.New("Logging.Adapter: An adapter is required. Valid Options: SQLite3"),
		errors.New("Logging.Database: You must specify the name of the logging database."),
		errors.New("[users.zamn] -> You must specify a nickname in order to connect to an IRC server."),
		errors.New("[users.zamn] -> You must specify a alternate nickname in order to connect to an IRC server."), errors.New("[users.zamn] -> You must specify the CA for your certificate to verify."), errors.New("[users.zamn] -> You must have at least 1 certificate on your user in order to authenticate.")}

	sort.Sort(ByErrorFunc(expected))
	sort.Sort(ByErrorFunc(err))

	if len(err) != 6 {
		log.Fatalf("Invalid number of errors returned for %s", PartialFile)
	}

	for i := 0; i < 6; i++ {
		if err[i].Error() != expected[i].Error() {
			log.Fatalf("Expected: \"%s\" and got: \"%s\"", expected[i], err[i])
		}
	}

}
