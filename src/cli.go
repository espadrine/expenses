package expense

import (
	"fmt"
	"log"
	"os"
	"slices"
	"strings"
)

var toplevelCommand Command

func init() {
	toplevelCommand = Command{
		Names: []string{"expense"},
		doc: "analyze expenses.\n\n" +
			"The commands rely on the following environment variables:\n" +
			"- $USER to determine the current username to associate with created operation;\n" +
			"- $EXPENSE_DB to define the location of the sqlite database;\n" +
			"  the default is in ~/.config/expense/db.sqlite.",
		Execute: printHelp,
		id:      1,
		subcommands: []Command{
			{
				Names:   []string{"help", "--help", "-h"},
				doc:     "print this usage information.",
				Execute: printHelp,
				id:      2,
			},
			{
				Names:   []string{"user"},
				doc:     "view or modify users.",
				Execute: listUsers,
				id:      3,
				subcommands: []Command{
					{
						Names:   []string{"help", "--help", "-h"},
						doc:     "print this usage information.",
						Execute: printHelp,
						id:      4,
					},
					{
						Names:   []string{"list"},
						doc:     "list the known usernames.",
						Execute: listUsers,
						id:      5,
					},
					{
						Names: []string{"create"},
						doc: "add a new user. " +
							"It takes a single parameter, its username, and returns its user ID.",
						Execute: createUser,
						id:      6,
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
	id          int // These IDs only stay the same within an execution.
	subcommands []Command
}

func (c *Command) String() string {
	return fmt.Sprintf("Command{id: %d, name: %s}", c.id, c.Names[0])
}

type Params struct {
	Command      Command
	CommandChain []Command
	Args         []string
}

func (p *Params) String() string {
	var chainStrings []string
	for _, command := range p.CommandChain {
		chainStrings = append(chainStrings, command.String())
	}
	chain := "[" + strings.Join(chainStrings, ", ") + "]"
	return fmt.Sprintf("Params{%v, chain: %s}", &p.Command, chain)
}

func ParseFlags(args []string) Params {
	var params Params
	params.CommandChain, params.Args = parseCommandChain(args, toplevelCommand)
	params.Command = params.CommandChain[len(params.CommandChain)-1]
	return params
}

func parseCommandChain(args []string, command Command) ([]Command, []string) {
	if len(args) == 0 {
		return []Command{command}, []string{} // Default command (typically the help).
	}
	commandName := args[0]
	for _, c := range command.subcommands {
		if slices.Contains(c.Names, commandName) {
			if len(args) > 1 && isSubcommandName(args[1], c) {
				cmdChain, cmdArgs := parseCommandChain(args[1:], c)
				return append([]Command{command}, cmdChain...), cmdArgs
			} else {
				return []Command{command, c}, args[1:]
			}
		}
	}
	return []Command{command}, args[1:] // Default command (typically the help).
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

	commandDoc := lastCommand.doc

	var commandsHelp string
	if hasCommands {
		commandsHelp = "Commands:\n"
		for _, command := range lastCommand.subcommands {
			commandsHelp += "- " + command.Names[0] + ": " + command.doc + "\n"
		}
	}

	return usageLine + "\n\nUse it to " + commandDoc + "\n\n" + commandsHelp
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

func createUser(params *Params, store *Store) {
	user, err := store.createUser(params.Args[0])
	if err != nil {
		log.Fatalf("createUsers: %s\n", err)
	}
	fmt.Println(user.id)
}
