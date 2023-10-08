package main

import (
	"fmt"
	"io"
	"log"
	"os"
	"os/exec"
	"regexp"
	"slices"
	"strings"
	"testing"
)

func TestMain(t *testing.T) {
	// Smoe tests require sequential execution of commands,
	// to check the effect of previous commands.
	type command struct {
		args      []string
		buildArgs func(ti int, env map[string]string) []string
		stdout    func(stdout string, env map[string]string) bool
		stderr    func(stderr string, env map[string]string) bool
	}
	tests := []struct {
		commands []command
	}{

		// Test the help command.
		{
			commands: []command{
				{
					args: []string{"help"},
					stderr: func(stderr string, env map[string]string) bool {
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
					stderr: func(stderr string, env map[string]string) bool {
						return stderr == "usage: expense user [commands]\n"+
							"\n"+
							"Use it to view or modify users.\n"+
							"\n"+
							"Commands:\n"+
							"- help: print this usage information.\n"+
							"- list: list the known usernames.\n"+
							"- create: add a new user. It takes a single parameter, its username, and returns its user ID.\n"+
							"- name: get the username associated with a user ID.\n"+
							"- id: get the user IDs associated with a username.\n"
					},
				},
			},
		},

		// Test the user command.
		{
			commands: []command{
				{
					buildArgs: func(ti int, env map[string]string) []string {
						os.Setenv("USER", "username")
						os.Setenv("EXPENSE_DB", fmt.Sprintf("./.test%d.sqlite", ti))
						return []string{"user", "list"}
					},
					stdout: func(stdout string, env map[string]string) bool {
						matched, err := regexp.MatchString("[a-z2-7]{26}\tusername", stdout)
						if err != nil {
							log.Fatal(err)
						}
						return matched
					},
				},
				{
					args: []string{"user", "create", "archimedes"},
					stdout: func(stdout string, env map[string]string) bool {
						env["archimedesUserID"] = strings.TrimSpace(stdout)
						return true
					},
				},
				{
					args: []string{"user", "list"},
					stdout: func(stdout string, env map[string]string) bool {
						matched, err := regexp.MatchString("[a-z2-7]{26}\t(username|archimedes)\n[a-z2-7]{26}\t(username|archimedes)", stdout)
						if err != nil {
							log.Fatal(err)
						}
						return matched
					},
				},
				{
					buildArgs: func(ti int, env map[string]string) []string {
						return []string{"user", "name", env["archimedesUserID"]}
					},
					stdout: func(stdout string, env map[string]string) bool {
						return stdout == "archimedes\n"
					},
				},
				{
					args: []string{"user", "id", "archimedes"},
					stdout: func(stdout string, env map[string]string) bool {
						return stdout == env["archimedesUserID"]+"\n"
					},
				},
				{
					args: []string{"user", "create", "archimedes"},
					stdout: func(stdout string, env map[string]string) bool {
						env["archimedesUserID2"] = strings.TrimSpace(stdout)
						return true
					},
				},
				{
					args: []string{"user", "id", "archimedes"},
					stdout: func(stdout string, env map[string]string) bool {
						ids := strings.Split(strings.TrimSpace(stdout), "\n")
						return len(ids) == 2 &&
							slices.Contains(ids, env["archimedesUserID"]) &&
							slices.Contains(ids, env["archimedesUserID2"])
					},
				},
			},
		},
	}

	for ti, test := range tests {
		// A fresh environment for each test,
		// to communicate data between commands.
		env := make(map[string]string)
		for ci, command := range test.commands {
			args := command.args
			if command.buildArgs != nil {
				args = command.buildArgs(ti, env)
			}
			stdout, stderr, err := execute(args)
			if err != nil {
				log.Fatal("Error on test #", ti, ": ", err, "\n", "stdout: ", stdout, "\n", "stderr: ", stderr)
			}

			if command.stdout != nil && !command.stdout(stdout, env) {
				t.Errorf("Test %d: command %d `%s` yielded invalid output\n\n%s",
					ti, ci, "expense "+strings.Join(command.args, " "), stdout)
			}
			if command.stderr != nil && !command.stderr(stderr, env) {
				t.Errorf("Test %d: command %d `%s` yielded invalid error message\n\n%s",
					ti, ci, "expense "+strings.Join(command.args, " "), stderr)
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
