package expense

import (
	"fmt"
	"log"
	"os"
	"slices"
)

var toplevelCommand Command

func init() {
	toplevelCommand = Command{
		Names:   []string{"expense"},
		doc:     "analyze expenses",
		Execute: printHelp,
		subcommands: []Command{
			{
				Names:   []string{"help", "--help", "-h"},
				doc:     "print this usage information",
				Execute: printHelp,
			},
			{
				Names:   []string{"user"},
				doc:     "analyze or modify users",
				Execute: listUsers,
				subcommands: []Command{
					{
						Names:   []string{"help", "--help", "-h"},
						doc:     "print this usage information",
						Execute: printHelp,
					},
					{
						Names:   []string{"list"},
						doc:     "list the known usernames",
						Execute: listUsers,
					},
				},
			},
		},
	}
}

type Command struct {
	Names       []string
	doc         string
	Execute     func(*Params, *Store)
	subcommands []Command
}

type Params struct {
	Command      Command
	CommandChain []Command
}

func ParseFlags() Params {
	var params Params
	params.CommandChain = parseCommandChain(os.Args[1:], toplevelCommand)
	params.Command = params.CommandChain[len(params.CommandChain)-1]
	return params
}

func parseCommandChain(args []string, command Command) []Command {
	if len(args) == 0 {
		return []Command{command} // Default command (typically the help).
	}
	commandName := args[0]
	for _, c := range command.subcommands {
		if slices.Contains(c.Names, commandName) {
			if len(args) > 1 && isSubcommandName(args[1], c) {
				return append([]Command{command}, parseCommandChain(args[1:], c)...)
			} else {
				return []Command{command, c}
			}
		}
	}
	return []Command{command} // Default command (typically the help).
}

func isSubcommandName(name string, command Command) bool {
	for _, subcom := range command.subcommands {
		if slices.Contains(subcom.Names, name) {
			return true
		}
	}
	return false
}

func printHelp(params *Params, store *Store) {
	fmt.Fprint(os.Stderr, helpString(params))
}

func helpString(params *Params) string {
	usageLine := "usage: "
	cclen := len(params.CommandChain)
	// If there are no commands given, add the help command.
	if cclen == 1 {
		params.CommandChain = append(params.CommandChain, params.Command.subcommands[0])
		cclen = len(params.CommandChain)
	}
	for _, command := range params.CommandChain[:cclen-1] {
		usageLine += command.Names[0] + " "
	}
	lastCommand := params.CommandChain[cclen-2]
	hasCommands := len(lastCommand.subcommands) > 0
	if hasCommands {
		usageLine += "[commands]"
	}

	var commandHelp string
	if hasCommands {
		commandHelp = "Commands:\n"
		for _, command := range lastCommand.subcommands {
			commandHelp += "- " + command.Names[0] + ": " + command.doc + "\n"
		}
	}

	return usageLine + "\n\n" + commandHelp
}

func listUsers(params *Params, store *Store) {
	users, err := store.getUsers()
	if err != nil {
		log.Fatalf("listUsers: %s\n", err)
	}
	for _, user := range users {
		fmt.Println(user.id + "\t" + user.name)
	}
}
