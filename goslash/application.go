package goslash

import (
	"github.com/bwmarrin/discordgo"
	"net/http"
)

type Application struct {
	http.Client

	Commands        map[string]*Command
	AuthHeader      string
	ClientID        string
	DefaultResponse *InteractionResponse

	Session *discordgo.Session // TODO dont depend on discordgo
}
type InteractionHandler func(interaction *Interaction) *InteractionResponse

// a Listener handles receiving and responding to an interaction while letting an InteractionHandler decide the response
// some listeners supported included with goslash are lambda.Listener (for aws lambda), httplistener.Listener and gateway.Listener (for receiving interactions from the discord gateway)
type Listener interface {
	SetHandler(handler InteractionHandler)
}

func (app *Application) SetListener(listener Listener) {
	listener.SetHandler(func(interaction *Interaction) *InteractionResponse {
		return app.HandleInteraction(interaction)
	})
}

func NewApp(clientId, auth string) (*Application, error) {
	session, err := discordgo.New(auth)
	if err != nil {
		return nil, err
	}

	return &Application{
		Client:          http.Client{},
		Commands:        make(map[string]*Command),
		AuthHeader:      auth,
		ClientID:        clientId,
		Session:         session,
		DefaultResponse: Response("Sorry, a response for that command could not be found").OnlyAuthor(),
	}, nil
}

func (app *Application) GetCommand(name string) *Command {
	if cmd, ok := app.Commands[name]; ok {
		return cmd
	}

	return nil
}

func (app *Application) HandleInteraction(interaction *Interaction) *InteractionResponse {
	resp := app.getResponse(interaction)
	if resp == nil {
		resp = app.DefaultResponse
	}

	return resp
}

func (app *Application) getResponse(interaction *Interaction) *InteractionResponse {

	rootCommand := app.GetCommand(interaction.Data.Name)
	if rootCommand == nil {
		return nil
	}

	ctx := &InteractionContext{interaction, app, []string{rootCommand.Name}}

	if resp := rootCommand.Handle(ctx); resp != nil {
		return resp
	}

	return handleOptions(rootCommand, ctx, interaction.Data.Options)
}

func handleOptions(rootCommand *Command, ctx *InteractionContext, options []ApplicationCommandInteractionDataOption) *InteractionResponse {
	for _, option := range options {
		ctx.InvokedCommands = append(ctx.InvokedCommands, option.Name)

		if resp := rootCommand.Handle(ctx); resp != nil {
			return resp
		}

		if option.Options == nil || len(options) < 1 {
			continue
		}

		if resp := handleOptions(rootCommand, ctx, option.Options); resp != nil {
			return resp
		}
	}

	return nil
}
