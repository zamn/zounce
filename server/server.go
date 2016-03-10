package server

import (
	"fmt"

	"github.com/spacemonkeygo/openssl"
)

func Start() error {
	fmt.Println("Starting Server!")
	serv_ctx, err := openssl.NewCtxFromFiles("keys/zamn.net.cert", "keys/zamn.net.key")

	if err != nil {
		panic(err)
	}

	l, err := openssl.Listen("tcp", ":7777", serv_ctx)

	if err != nil {
		panic(err)
	}

	conn, err := l.Accept()
	fmt.Printf("Got: %#V\n", conn)

	fmt.Printf("Addr", l.Addr())

	return nil
}
