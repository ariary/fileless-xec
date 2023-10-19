package main

import (
	"flag"
	"fmt"
	"net/http"
	"path"
	"runtime"
	"strings"
	"sync"

	_ "net/http/pprof"

	"github.com/lucas-clemente/quic-go/http3"
	"github.com/quic-go/quic-go"
)

type binds []string

func (b binds) String() string {
	return strings.Join(b, ",")
}

func (b *binds) Set(v string) error {
	*b = strings.Split(v, ",")
	return nil
}

// Size is needed by the /demo/upload handler to determine the size of the uploaded file
type Size interface {
	Size() int64
}

func setupHandler(www string) http.Handler {
	mux := http.NewServeMux()

	mux.Handle("/", http.FileServer(http.Dir(www)))
	return mux
}

func main() {
	port := flag.String("p", "6121", "http3 server port")

	bs := binds{}
	www := "."

	flag.Parse()
	fmt.Printf("Listening on 0.0.0.0:%s...\n", *port)
	bs = binds{"0.0.0.0:" + *port}

	handler := setupHandler(www)
	quicConf := &quic.Config{}

	var wg sync.WaitGroup
	wg.Add(len(bs))
	for _, b := range bs {
		bCap := b
		go func() {
			var err error

			server := http3.Server{
				Server:     &http.Server{Handler: handler, Addr: bCap},
				QuicConfig: quicConf,
			}

			//Get cert
			_, filename, _, ok := runtime.Caller(0)
			if !ok {
				panic("Failed to get current frame")
			}
			certPath := path.Dir(filename)
			cert := path.Join(certPath, "cert.pem")
			privateKey := path.Join(certPath, "priv.key")

			err = server.ListenAndServeTLS(cert, privateKey)

			if err != nil {
				fmt.Println(err)
			}
			wg.Done()
		}()
	}
	wg.Wait()
}

//sudo tcpdump -i lo udp port 6121 -vv -X
