package expense

import "testing"

func TestFlags(t *testing.T) {
	tests := []struct {
		args   []string
		params Params
	}{
		{
			args: []string{},
			params: Params{
				Command: toplevelCommand,
				CommandChain: []Command{
					toplevelCommand,
				},
			},
		},
		{
			args: []string{"help"},
			params: Params{
				Command: toplevelCommand.subcommands[0],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[0],
				},
			},
		},
		{
			args: []string{"--help"},
			params: Params{
				Command: toplevelCommand.subcommands[0],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[0],
				},
			},
		},
		{
			args: []string{"-h"},
			params: Params{
				Command: toplevelCommand.subcommands[0],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[0],
				},
			},
		},
		{
			args: []string{"user"},
			params: Params{
				Command: toplevelCommand.subcommands[1],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[1],
				},
			},
		},
		{
			args: []string{"user", "-h"},
			params: Params{
				Command: toplevelCommand.subcommands[1].subcommands[0],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[1],
					toplevelCommand.subcommands[1].subcommands[0],
				},
			},
		},
		{
			args: []string{"user", "list"},
			params: Params{
				Command: toplevelCommand.subcommands[1].subcommands[1],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[1],
					toplevelCommand.subcommands[1].subcommands[1],
				},
			},
		},
		{
			args: []string{"user", "create", "archimedes"},
			params: Params{
				Command: toplevelCommand.subcommands[1].subcommands[2],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[1],
					toplevelCommand.subcommands[1].subcommands[2],
				},
				Args: []string{"archimedes"},
			},
		},
		{
			args: []string{"user", "name", "6gqsxomsj3hir4msckf4wxtmsuoqa3w5"},
			params: Params{
				Command: toplevelCommand.subcommands[1].subcommands[3],
				CommandChain: []Command{
					toplevelCommand,
					toplevelCommand.subcommands[1],
					toplevelCommand.subcommands[1].subcommands[3],
				},
				Args: []string{"6gqsxomsj3hir4msckf4wxtmsuoqa3w5"},
			},
		},
	}
	for _, test := range tests {
		params := ParseFlags(test.args)
		if !matchesParams(params, test.params) {
			t.Errorf("Arguments %v yielded parameters %v instead of %v", test.args, &params, &test.params)
		}
	}
}

func matchesParams(params Params, expectedParams Params) bool {
	sameCommand := params.Command.id == expectedParams.Command.id
	sameCommandChainLength := len(params.CommandChain) == len(expectedParams.CommandChain)
	chainContainsExpectedCommands := true
	if sameCommandChainLength {
		for i, command := range params.CommandChain {
			chainContainsExpectedCommands = command.id == expectedParams.CommandChain[i].id
			if !chainContainsExpectedCommands {
				break
			}
		}
	}
	sameCommandChain := sameCommandChainLength && chainContainsExpectedCommands

	sameArgs := true
	for i, arg := range params.Args {
		if arg != expectedParams.Args[i] {
			sameArgs = false
			break
		}
	}

	return sameCommand && sameCommandChain && sameArgs
}
