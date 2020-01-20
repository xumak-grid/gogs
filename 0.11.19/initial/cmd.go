package main

import (
	"bufio"
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

//NewCMD returns a new Command
func NewCMD(cmd string, args ...string) Command {
	return Command{CmdName: cmd, CmdArgs: args}
}

// Command represents a Command to exec in the terminal
type Command struct {
	CmdName string
	CmdArgs []string
}

//Exec execs the command and shows the Output
func (c *Command) Exec() error {
	fmt.Println(c.CmdName, strings.Join(c.CmdArgs, " "))

	var stderr bytes.Buffer
	cmd := exec.Command(c.CmdName, c.CmdArgs...)

	cmdReader, err := cmd.StdoutPipe()
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		return err
	}
	cmd.Stderr = &stderr

	scanner := bufio.NewScanner(cmdReader)
	go func() {
		for scanner.Scan() {
			fmt.Println(scanner.Text())
		}
	}()

	err = cmd.Start()
	if err != nil {
		fmt.Fprintln(os.Stderr, stderr.String(), err)
		return err
	}

	err = cmd.Wait()
	if err != nil {
		fmt.Fprintln(os.Stderr, stderr.String(), err)
		return err
	}

	return nil
}
