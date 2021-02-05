package goslash

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// InteractionContext contains an Interaction and other relevant information handlers need
type InteractionContext struct {
	*Interaction
	App *Application

	InvokedCommands []string
}

func (ctx *InteractionContext) Respond(response *InteractionResponse) {
	_, err := ctx.App.PostJson(fmt.Sprintf("interactions/%v/%v/callback", ctx.Interaction.ID, ctx.Interaction.Token), response)
	if err != nil {
		log.Error("Error responding to reaction: ", response)
	}
}

func (ctx *InteractionContext) Send(message *discordgo.WebhookParams) (*discordgo.Message, error){
	return ctx.App.Session.WebhookExecute(ctx.App.ClientID, ctx.Token, true, message)
}

// InteractionContext.GetOption returns an ApplicationCommandInteractionDataOption from the interactions list of options
func (ctx *InteractionContext) GetOption(nameToFind string) *ApplicationCommandInteractionDataOption {
	options := ctx.Data.Options

	return getOption(ctx.InvokedCommands[1:], nameToFind, options)
}

func getOption(invokedCommands []string, nameToFind string, options []ApplicationCommandInteractionDataOption) *ApplicationCommandInteractionDataOption {
	for _, option := range options {
		if option.Name == nameToFind && len(invokedCommands) == 0 {
			return &option
		}
		if len(invokedCommands) > 0 && option.Name == invokedCommands[0] {
			return getOption(invokedCommands[1:], nameToFind, option.Options)
		}
	}

	return nil
}

func (ctx *InteractionContext) GetString(name string) string {
	option := ctx.GetOption(name)
	if option == nil {
		return ""
	}

	return option.Value.(string)
}
func (ctx *InteractionContext) GetInt(name string) int {
	option := ctx.GetOption(name)
	if option == nil {
		return 0
	}
	value := option.Value.(float64)
	return int(value)
}
