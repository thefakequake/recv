package recv

import (
	"strings"

	"github.com/bwmarrin/discordgo"
)

// stores command routes
type CommandRouter struct {
	// maps command names to command objects
	commands map[string]Command
	// maps alias names to command names
	aliases  map[string]string
}

// creates a new CommandRouter
func NewCommandRouter() CommandRouter {
	return CommandRouter{
		commands: map[string]Command{},
		aliases:  map[string]string{},
	}
}

// adds a command to the command router's map
func (r CommandRouter) AddCommand(comm Command) {
	commID := strings.ToLower(comm.Name)

	for _, a := range append(comm.Aliases, commID) {
		r.aliases[strings.ToLower(a)] = commID
	}

	r.commands[commID] = comm
}

// fetches a command from the command router's map via name
func (r CommandRouter) GetCommand(name string) (Command, bool) {
	commID, ok := r.aliases[strings.ToLower(name)]
	if !ok {
		return Command{}, false
	}

	return r.commands[commID], true
}

// the result of calling ProcessCommands, contains the handler for the function and command that was parsed
type ProcessResult struct {
	// command that was processed from the message, nil if no command was found
	Command  *Command
	// callback function that runs the command callback
	Callback func()
}

// uses the prefix of the bot, session and message object in order to determine the handler and context of the command
func (r CommandRouter) ProcessCommands(prefix string, session *discordgo.Session, message *discordgo.MessageCreate) (*ProcessResult, error) {
	result := ProcessResult{}
	if !strings.HasPrefix(message.Content, prefix) {
		return &result, nil
	}
	splitParts := strings.Split(strings.TrimPrefix(message.Content, prefix), " ")
	if len(splitParts) == 0 {
		return &result, nil
	}

	commName := splitParts[0]
	commArgs := splitParts[1:]

	comm, ok := r.GetCommand(commName)
	if !ok {
		return &result, nil
	}

	result.Command = &comm
	ctx := Ctx{
		Command: &comm,
		Session: session,
		Message: message,
	}

	for _, c := range comm.Checks {
		pass := c.Callback(&ctx)
		if !pass {
			return &result, CheckError{
				Check: &c,
			}
		}
	}

	processedArgs := []interface{}{}

	if !comm.NoJoin && len(commArgs) > len(comm.Args) {
		commArgs = append(commArgs[:len(comm.Args)-1], strings.Join(commArgs[len(comm.Args)-1:], " "))
	}

	for i, arg := range comm.Args {
		if i+1 > len(commArgs) {
			if !arg.Optional {
				return &result, MissingRequiredArgumentError{
					Arg:         &arg,
					ArgPosition: i + 1,
				}
			}
			break
		}

		var convert interface{} = commArgs[i]
		if arg.Converter != nil {
			converterCtx := ConverterCtx{
				Ctx:   &ctx,
				Arg:   &arg,
				Input: convert.(string),
			}
			output, ok := arg.Converter.Callback(&converterCtx)
			if !ok {
				return &result, ConversionError{
					ArgPosition: i + 1,
					Arg: &arg,
					Input: convert.(string),
				}
			}
			convert = output
		}

		processedArgs = append(processedArgs, convert)
	}

	result.Callback = func() {
		comm.Callback(&CommandCtx{
			Ctx:  &ctx,
			Args: processedArgs,
		})
	}

	return &result, nil
}
