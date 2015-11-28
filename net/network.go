package net

import "github.com/zamN/zounce/net/perform"

type Network struct {
	Name        string   `validate:"nonzero"`
	Servers     []string `validate:"min=1"`
	Password    string
	PerformInfo perform.Perform `toml:"perform"`
}
