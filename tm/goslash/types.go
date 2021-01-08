package main

import (
	"time"
)

type InteractionType int

const (
	PING InteractionType = iota + 1
	APPLICATION_COMMAND
)

type GuildMember struct {
	Deaf         bool        `json:"deaf"`
	IsPending    bool        `json:"is_pending"`
	JoinedAt     time.Time   `json:"joined_at"`
	Mute         bool        `json:"mute"`
	Nick         interface{} `json:"nick"`
	Pending      bool        `json:"pending"`
	Permissions  string      `json:"permissions"`
	PremiumSince interface{} `json:"premium_since"`
	Roles        []string    `json:"roles"`
	User         struct {
		Avatar        string `json:"avatar"`
		Discriminator string `json:"discriminator"`
		ID            string `json:"id"`
		PublicFlags   int    `json:"public_flags"`
		Username      string `json:"username"`
	} `json:"user"`
}

type ApplicationCommandInteractionDataOption struct {
	Name    string                                    `json:"name"`
	Value   interface{}                               `json:"value"`
	Options []ApplicationCommandInteractionDataOption `json:"options"`
}

type ApplicationCommandInteractionData struct {
	ID      string                                    `json:"id"`
	Name    string                                    `json:"name"`
	Options []ApplicationCommandInteractionDataOption `json:"options"`
}

type Interaction struct {
	ID        string                            `json:"id"`
	Type      InteractionType                   `json:"type"`
	ChannelID string                            `json:"channel_id"`
	Data      ApplicationCommandInteractionData `json:"data"`
	GuildID   string                            `json:"guild_id"`
	Member    GuildMember                       `json:"member"`
	Token     string                            `json:"token"`
	Version   int                               `json:"version"`
}
