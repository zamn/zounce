package config

import (
	"fmt"
	"log"
	"os"
	"os/exec"
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

	fmt.Println(err)
	if len(err) == 0 {
		log.Fatalf("Validated an empty configuration file!\n")
	}
}

func TestPartialFileSafeDefaults(t *testing.T) {
	c, err := LoadConfig(PartialFile)

	if len(err) == 0 {
		log.Fatalf("Failed to load config file! Error: %s\n", err)
	}
	fmt.Println("C", *c)
}
