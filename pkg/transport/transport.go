package transport

import (
	"bytes"
	"crypto/tls"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"sync"

	"github.com/quic-go/quic-go"
	"github.com/quic-go/quic-go/http3"
)

// Retrieve Binary from Remote using http protocol
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

// Retrieve Binary from Remote using http3 protocol (Quic)
// Put ca.pem in current directory
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
