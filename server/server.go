package server

import (
	"fmt"
	"io"
	"net"

	"github.com/spacemonkeygo/openssl"
)

func Start() error {
	fmt.Println("Starting Server!")
	servCtx, err := openssl.NewCtxFromFiles("keys/zamn.net.cert", "keys/zamn.net.key")

	if err != nil {
		panic(err)
	}

	l, err := openssl.Listen("tcp", ":7777", servCtx)

	if err != nil {
		panic(err)
	}

	defer l.Close()

	sslListener := openssl.NewListener(l, servCtx)

	for {
		conn, err := sslListener.Accept()
		if err != nil {
			panic(err)
		}
		fmt.Printf("Got: %#V\n", conn)

		fmt.Println("SSL", conn)

		go func(c net.Conn) {
			io.Copy(c, c)
			c.Close()
		}(conn)
	}

	return nil
}
