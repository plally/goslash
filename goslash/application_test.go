package goslash

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestApplication_HandleInteraction(t *testing.T) {
	app, _ := NewApp("123", "BOT Token")
	pingCommand := NewCommand("ping", "pong!")
	pingCommand.SetHandler("ping", pingHandler)
	app.AddCommand(pingCommand)

	resp := app.HandleInteraction(&Interaction{
		ID:        "123",
		Type:      2,
		ChannelID: "123",
		Data:      ApplicationCommandInteractionData{
			ID:      "123",
			Name:    "ping",
			Options: nil,
		},
		GuildID:   "123",
		Member:    GuildMember{},
		Token:     "abcdefg",
		Version:   1,
	})

	assert.Equal(t, resp, Response("pong!").KeepSource(), "resp should be equal to 'pong!'")

}

func pingHandler(ctx *InteractionContext) *InteractionResponse {
  return Response("pong!").KeepSource()
}