package cert

type Cert struct {
	Path string `toml:"cert_path" validate:"nonzero"`
}
