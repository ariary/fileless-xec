package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
)

func main() {
	//Read shell
	command := ""
	script, err := os.Open("myscript.sh")
	if err != nil {
		log.Fatal(err)
	}
	defer script.Close()
	scanner := bufio.NewScanner(script)

	for scanner.Scan() {
		//todo: Replace " by ' if possible"
		command += scanner.Text() + ";"
	}

	if err := scanner.Err(); err != nil {
		log.Fatal(err)
	}

	fmt.Println(command)

	//Write go program
	program := `
package main

import (
	"fmt"
	"os/exec"
)

func main() {
	//Shell content hardcoded
	content := "` + command + `"

	//execute binary
	cmd := exec.Command("/bin/sh", "-c", content)
	stdout, err := cmd.Output()

	if err != nil {
		fmt.Println(err.Error())
		return
	}

	// Print the output
	fmt.Println(string(stdout))
}`
	f, err := os.Create("nestedscript.go")

	if err != nil {
		log.Fatal(err)
	}

	defer f.Close()

	_, err = f.WriteString(program)

	if err != nil {
		log.Fatal(err)
	}
}
