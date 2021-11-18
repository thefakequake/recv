package recv

import (
	"fmt"
	"strings"

	"github.com/bwmarrin/discordgo"
)

// represents a command that users can use
type Command struct {
	Name        string
	Description string
	Args        []CommandArg
	Aliases     []string
	NoJoin      bool
	Checks      []CommandCheck
	Callback    func(*CommandCtx)
}

// stores basic context
type Ctx struct {
	Command *Command
	Session *discordgo.Session
	Message *discordgo.MessageCreate
}

// stores contextual information for a command such as arguments, session and message
type CommandCtx struct {
	*Ctx
	Args []interface{}
}

// stores contextual information for a converter such as the argument, input, session and message
type ConverterCtx struct {
	*Ctx
	Arg   *CommandArg
	Input string
}

// argument for a command
type CommandArg struct {
	Name        string
	Optional    bool
	Description string
	Converter   *Converter
}

// a callback that decides whether a command can be run or not
type CommandCheck struct {
	Name     string
	Callback func(*Ctx) bool
}

// a callback that "converts" a string input into something else
type Converter struct {
	Name     string
	Callback func(*ConverterCtx) (interface{}, bool)
}

// helper function that gets the usage for a command minus the prefix
func (c Command) Usage() string {
	var args []string

	for _, arg := range c.Args {
		formattedArg := arg.Name
		if arg.Converter != nil {
			formattedArg = fmt.Sprintf("%s: %s", formattedArg, arg.Converter.Name)
		}
		if arg.Optional {
			formattedArg = fmt.Sprintf("[%s]", formattedArg)
		} else {
			formattedArg = fmt.Sprintf("<%s>", formattedArg)
		}
		args = append(args, formattedArg)
	}

	return strings.Join(append([]string{strings.ToLower(c.Name)}, args...), " ")
}
