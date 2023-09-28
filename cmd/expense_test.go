package main

import (
	"io"
	"log"
	"os/exec"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	// Smoe tests require sequential execution of commands,
	// to check the effect of previous commands.
	type command struct {
		args   []string
		stdout string
		stderr string
	}
	tests := []struct {
		commands []command
	}{
		{
			commands: []command{
				{
					args: []string{"help"},
					stderr: "usage: expense [commands]\n" +
						"\n" +
						"Commands:\n" +
						"- help: print this usage information\n" +
						"- user: analyze or modify users\n",
				},
			},
		},
		{
			commands: []command{
				{
					args: []string{"user", "help"},
					stderr: "usage: expense user [commands]\n" +
						"\n" +
						"Commands:\n" +
						"- help: print this usage information\n" +
						"- list: list the known usernames\n",
				},
			},
		},
	}
	for ti, test := range tests {
		for _, command := range test.commands {
			stdout, stderr := execute(command.args)
			if stdout != command.stdout {
				t.Errorf("Test %d: command `%s` yielded output\n\n%s\ninstead of\n\n%s",
					ti, "expense "+strings.Join(command.args, " "), stdout, command.stdout)
			}
			if stderr != command.stderr {
				t.Errorf("Test %d: command `%s` yielded error\n\n%s\ninstead of\n\n%s",
					ti, "expense "+strings.Join(command.args, " "), stderr, command.stderr)
			}
		}
	}
}

func execute(args []string) (stdout string, stderr string) {
	cmd := exec.Command("../expense", args...)
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		log.Fatal(err)
	}
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		log.Fatal(err)
	}

	if err := cmd.Start(); err != nil {
		log.Fatal(err)
	}

	stdoutBytes, err := io.ReadAll(stdoutReader)
	if err != nil {
		log.Fatal(err)
	}
	stdout = string(stdoutBytes)
	stderrBytes, err := io.ReadAll(stderrReader)
	if err != nil {
		log.Fatal(err)
	}
	stderr = string(stderrBytes)

	if err := cmd.Wait(); err != nil {
		log.Fatal(err)
	}
	return
}
