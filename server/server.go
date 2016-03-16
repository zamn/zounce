package server

import (
	"fmt"
	"io/ioutil"
	"net"

	"github.com/spacemonkeygo/openssl"
)

func Start() error {
	fmt.Println("Starting Server!")
	servCtx, err := openssl.NewCtxFromFiles("keys/zamn.net.cert", "keys/zamn.net.key")

	certFile, err := ioutil.ReadFile("keys/zamn.net.chain")

	if err != nil {
		panic(err)
	}

	cert, err := openssl.LoadCertificateFromPEM(certFile)
	servCtx.AddChainCertificate(cert)

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
			msg := make([]byte, 2048)
			if err != nil {
				panic(err)
			}

			n, err := c.Read(msg)
			if err != nil {
				panic(err)
			}
			fmt.Println(n, msg)
		}(conn)
	}

	return nil
}
