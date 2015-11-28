package logging

type LogInfo struct {
	Adapter  string `toml:"adapter" validate:"nonzero"`
	Database string `toml:"database" validate:"nonzero"`
}
