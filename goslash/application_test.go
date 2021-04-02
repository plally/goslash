package goslash

import (
	"github.com/bwmarrin/discordgo"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplication_HandleInteraction(t *testing.T) {
	app, _ := NewApp("123", "BOT Token")
	pingCommand := NewCommand("ping", "pong!")
	pingCommand.SetHandler(pingHandler)
	app.AddCommand(pingCommand)

	resp := app.HandleInteraction(&discordgo.Interaction{
		ID:        "123",
		Type:      2,
		ChannelID: "123",
		Data: discordgo.ApplicationCommandInteractionData{
			ID:      "123",
			Name:    "ping",
			Options: nil,
		},
		GuildID: "123",
		Member:  &discordgo.Member{},
		Token:   "abcdefg",
		Version: 1,
	})

	assert.Equal(t, resp, Response("pong!").ToDiscordgo(), "resp should be equal to 'pong!'")

}

func pingHandler(ctx *InteractionUpdate) *InteractionResponse {
	return Response("pong!")
}
