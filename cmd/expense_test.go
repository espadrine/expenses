package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	// Smoe tests require sequential execution of commands,
	// to check the effect of previous commands.
	type command struct {
		before func(ti int)
		args   []string
		stdout func(stdout string) bool
		stderr func(stderr string) bool
	}
	tests := []struct {
		commands []command
	}{

		// Test the help command.
		{
			commands: []command{
				{
					args: []string{"help"},
					stderr: func(stderr string) bool {
						return stderr == "usage: expense [commands]\n"+
							"\n"+
							"Use it to analyze expenses.\n"+
							"\n"+
							"The commands rely on the following environment variables:\n"+
							"- $USER to determine the current username to associate with created operation;\n"+
							"- $EXPENSE_DB to define the location of the sqlite database;\n"+
							"  the default is in ~/.config/expense/db.sqlite.\n"+
							"\n"+
							"Commands:\n"+
							"- help: print this usage information.\n"+
							"- user: view or modify users.\n"
					},
				},
			},
		},
		{
			commands: []command{
				{
					args: []string{"user", "help"},
					stderr: func(stderr string) bool {
						return stderr == "usage: expense user [commands]\n"+
							"\n"+
							"Use it to view or modify users.\n"+
							"\n"+
							"Commands:\n"+
							"- help: print this usage information.\n"+
							"- list: list the known usernames.\n"+
							"- create: add a new user. It takes a single parameter, its username, and returns its user ID.\n"
					},
				},
			},
		},

		// Test the user command.
		{
			commands: []command{
				{
					before: func(ti int) {
						os.Setenv("USER", "username")
						os.Setenv("EXPENSE_DB", fmt.Sprintf("./.test%d.sqlite", ti))
					},
					args: []string{"user", "list"},
					stdout: func(stdout string) bool {
						matched, err := regexp.MatchString("[a-z2-7]{26}\tusername", stdout)
						if err != nil {
							log.Fatal(err)
						}
						return matched
					},
				},
				{
					args: []string{"user", "create", "archimedes"},
				},
				{
					args: []string{"user", "list"},
					stdout: func(stdout string) bool {
						matched, err := regexp.MatchString("[a-z2-7]{26}\tusername\n[a-z2-7]{26}\tarchimedes", stdout)
						if err != nil {
							log.Fatal(err)
						}
						return matched
					},
				},
			},
		},
	}

	for ti, test := range tests {
		for _, command := range test.commands {
			if command.before != nil {
				command.before(ti)
			}
			stdout, stderr, err := execute(command.args)
			if err != nil {
				log.Fatal("Error on test #", ti, ": ", err, "\n", "stdout: ", stdout, "\n", "stderr: ", stderr)
			}

			if command.stdout != nil && !command.stdout(stdout) {
				t.Errorf("Test %d: command `%s` yielded invalid output\n\n%s",
					ti, "expense "+strings.Join(command.args, " "), stdout)
			}
			if command.stderr != nil && !command.stderr(stderr) {
				t.Errorf("Test %d: command `%s` yielded invalid error message\n\n%s",
					ti, "expense "+strings.Join(command.args, " "), stderr)
			}
		}
	}
}

func execute(args []string) (stdout string, stderr string, err error) {
	cmd := exec.Command("../expense", args...)
	stdoutReader, err := cmd.StdoutPipe()
	if err != nil {
		return
	}
	stderrReader, err := cmd.StderrPipe()
	if err != nil {
		return
	}

	if err = cmd.Start(); err != nil {
		return
	}

	stdoutBytes, err := io.ReadAll(stdoutReader)
	if err != nil {
		return
	}
	stdout = string(stdoutBytes)
	stderrBytes, err := io.ReadAll(stderrReader)
	if err != nil {
		return
	}
	stderr = string(stderrBytes)

	if err = cmd.Wait(); err != nil {
		return
	}
	return
}
