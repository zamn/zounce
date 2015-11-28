package auth

type Auth struct {
	CAPath string `toml:"ca_path" validate:"nonzero"`
}
