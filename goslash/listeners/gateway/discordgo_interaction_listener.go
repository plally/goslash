package gateway

import (
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/plally/goslash/goslash"
	log "github.com/sirupsen/logrus"
)

const DISCORD_API_BASE_URL = "https://discord.com/api/v8"

type Listener struct {
	Session *discordgo.Session
	Handler goslash.InteractionHandler
}

func (listener *Listener) SetHandler(handler goslash.InteractionHandler) {
	listener.Handler = handler
}

func NewGatewayListener(botToken string) (*Listener, error) {
	session, err := discordgo.New("Bot " + botToken)
	if err != nil {
		return nil, err
	}

	listener := &Listener{
		Session: session,
		Handler: nil,
	}

	session.AddHandler(func(s *discordgo.Session, event *discordgo.Event) {
		if listener.Handler == nil || event.Type != "INTERACTION_CREATE" {
			return
		}

		var interaction goslash.Interaction
		err := json.Unmarshal(event.RawData, &interaction)
		if err != nil {
			log.WithField("error", err).Info("error unmarshalling gateway INTERACTION_CREATE data")
			return
		}

		response := listener.Handler(&interaction)
		if response != nil {
			url := fmt.Sprintf("%v/%v", DISCORD_API_BASE_URL, fmt.Sprintf("interactions/%v/%v/callback", interaction.ID, interaction.Token))

			_, err := session.Request("POST", url, response)
			if err != nil {
				log.Error(err)
			}
		}

	})
	err = session.Open()
	if err != nil {
		return nil, err
	}

	return listener, nil
}
