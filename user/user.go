package user

import (
	"github.com/zamN/zounce/logging"
	"github.com/zamN/zounce/net"
	"github.com/zamN/zounce/user/auth"
	"github.com/zamN/zounce/user/cert"
)

type User struct {
	Nick     string `validate:"nonzero,max=9"`
	AltNick  string `validate:"nonzero,max=9"`
	Username string
	Realname string
	Logging  logging.LogInfo
	AuthInfo auth.Auth              `toml:"auth"`
	Certs    map[string]cert.Cert   `validate:"min=1"`
	Networks map[string]net.Network `validate:"validnetworks"`
}
