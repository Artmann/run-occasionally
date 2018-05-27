package run_occasionally

import (
	"log"
	"os/exec"
	"os"
	"bytes"
)

type Command struct {
	childProcess *os.Process

}

func (cr *Command) Execute(command string, args string) *Command {
	var outputBuffer, errorBuffer bytes.Buffer

	cmd := exec.Command(command, args)
	cmd.Stdout = &outputBuffer
	cmd.Stderr = &errorBuffer

	cr.childProcess = cmd.Process

	go func() {

		err := cmd.Start()

		if err != nil {
			log.Println(err)
			return
		}

		cmd.Wait()

		log.Println("out:", outputBuffer.String(), "err:", errorBuffer.String())

		log.Println("Command Done")

	}()

	return cr
}

func (cr *Command) Abort() {
	if cr.childProcess != nil {
		cr.childProcess.Signal(os.Interrupt)
		log.Println("Aborted Command")
	}
}