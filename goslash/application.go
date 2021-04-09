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

type InteractionHandler func(interaction *discordgo.Interaction) *discordgo.InteractionResponse

// a Listener handles receiving and responding to an interaction while letting an InteractionHandler decide the response
// some listeners supported included with goslash are lambda.Listener (for aws lambda), httplistener.Listener and gateway.Listener (for receiving interactions from the discord gateway)
type Listener interface {
	SetHandler(handler InteractionHandler)
}

func (app *Application) SetListener(listener Listener) {
	listener.SetHandler(func(interaction *discordgo.Interaction) *discordgo.InteractionResponse {
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
		DefaultResponse: Response("A response for that command could not be found"),
	}, nil
}

func (app *Application) GetCommand(name string) *Command {
	if cmd, ok := app.Commands[name]; ok {
		return cmd
	}

	return nil
}

func (app *Application) AddCommand(command *Command) {
	app.Commands[command.Name] = command
}
