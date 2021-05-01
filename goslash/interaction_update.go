package goslash

import (
	"github.com/bwmarrin/discordgo"
	log "github.com/sirupsen/logrus"
)

// InteractionUpdate contains an Interaction and other relevant information handlers need
type InteractionUpdate struct {
	*discordgo.Interaction
	App *Application

	InvokedCommands []string
}

func (interaction *InteractionUpdate) IsDM() bool {
	return interaction.Interaction.GuildID == ""
}

func (interaction *InteractionUpdate) GetUser() *discordgo.User {
	if interaction.Member != nil {
		return interaction.Member.User
	}

	return interaction.User
}

// call the /callback endpoint for this interaction, not useful for http interaction handlers
func (interaction *InteractionUpdate) Callback(response *InteractionResponse) {
	err := interaction.App.Session.InteractionRespond(interaction.Interaction, response.ToDiscordgo())
	if err != nil {
		log.Error("Error responding to reaction: ", response)
	}
}

func (interaction *InteractionUpdate) Followup(params *discordgo.WebhookParams) (*discordgo.Message, error) {
	return interaction.App.Session.FollowupMessageCreate(interaction.App.ClientID, interaction.Interaction, true, params)
}

func (interaction *InteractionUpdate) EditOriginal(params *discordgo.WebhookEdit) error {
	return interaction.App.Session.InteractionResponseEdit(interaction.App.ClientID, interaction.Interaction, params)
}
func (interaction *InteractionUpdate) DeleteOriginal() error {
	return interaction.App.Session.InteractionResponseDelete(interaction.App.ClientID, interaction.Interaction)
}

// InteractionUpdate.GetOption returns an ApplicationCommandInteractionDataOption from the interactions list of options
func (interaction *InteractionUpdate) GetOption(nameToFind string) *discordgo.ApplicationCommandInteractionDataOption {
	options := interaction.Data.Options

	return getOption(interaction.InvokedCommands[1:], nameToFind, options)
}

func getOption(invokedCommands []string, nameToFind string, options []*discordgo.ApplicationCommandInteractionDataOption) *discordgo.ApplicationCommandInteractionDataOption {
	for _, option := range options {
		if option.Name == nameToFind && len(invokedCommands) == 0 {
			return option
		}
		if len(invokedCommands) > 0 && option.Name == invokedCommands[0] {
			return getOption(invokedCommands[1:], nameToFind, option.Options)
		}
	}

	return nil
}

func (interaction *InteractionUpdate) GetString(name string) string {
	option := interaction.GetOption(name)
	if option == nil {
		return ""
	}

	return option.Value.(string)
}

func (interaction *InteractionUpdate) GetInt(name string) int {
	option := interaction.GetOption(name)
	if option == nil {
		return 0
	}
	value := option.Value.(float64)
	return int(value)
}

func (interaction *InteractionUpdate) GetBool(name string) bool {
	option := interaction.GetOption(name)
	if option == nil {
		return false
	}

	return option.Value.(bool)
}