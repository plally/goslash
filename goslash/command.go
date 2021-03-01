package goslash

import (
	"strings"
)

type CommandHandler func(interaction *InteractionContext) *InteractionResponse

// Command stores an ApplicationCommand and a handler and adds some helper methods
type Command struct {
	ApplicationCommand
	checks   map[string][]CommandHandler
	handlers map[string]CommandHandler
	isGlobal bool
}

func (cmd *Command) IsGlobal() bool {
	return cmd.isGlobal
}

func (cmd *Command) SetOptions(options ...ApplicationCommandOption) *Command {
	cmd.Options = options
	return cmd
}

func (cmd *Command) AddSubCheck(name string, command CommandHandler) *Command{
	cmd.checks[name] = append(cmd.checks[name], command)

	return cmd
}

func (cmd *Command) AddCheck(command CommandHandler) *Command {
	cmd.AddSubCheck(cmd.Name, command)

	return cmd
}

func (cmd *Command) SetHandler(command CommandHandler) *Command {
	cmd.SetSubHandler(cmd.Name, command)

	return cmd
}

func (cmd *Command) SetSubHandler(name string, command CommandHandler) *Command {
	cmd.handlers[name] = command

	return cmd
}

func (cmd *Command) GetHandler(name string) CommandHandler {
	return cmd.handlers[name]
}

func (cmd *Command) Handle(interaction *InteractionContext) *InteractionResponse {
	name := strings.Join(interaction.InvokedCommands, " ")

	handler := cmd.GetHandler(name)

	if handler == nil {
		return nil
	}
	checks := cmd.checks[name]

	for _, check := range checks {

		if checkResponse := check(interaction); checkResponse != nil {
			return checkResponse
		}
	}

	return handler(interaction)
}

func NewCommand(name, description string) *Command {
	return &Command{
		ApplicationCommand: ApplicationCommand{
			Name:        name,
			Description: description,
			Options:     nil,
		},
		handlers: make(map[string]CommandHandler),
		checks:   make(map[string][]CommandHandler),
	}
}
