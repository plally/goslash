package goslash

import (
	"github.com/bwmarrin/discordgo"
	"time"
)

type ApplicationCommandOptionType int

const (
	SUB_COMMAND ApplicationCommandOptionType = iota + 1
	SUB_COMMAND_GROUP
	STRING
	INTEGER
	BOOLEAN
	USER
	CHANNEL
	ROLE
)

func Choice(name string, value string) ApplicationCommandOptionChoice {
	return ApplicationCommandOptionChoice{
		Name:  name,
		Value: value,
	}
}
func User(name string, description string, required bool) ApplicationCommandOption {
	return ApplicationCommandOption{
		Type:        USER,
		Name:        name,
		Description: description,
		Required:    required,
	}
}

func Role(name string, description string, required bool) ApplicationCommandOption {
	return ApplicationCommandOption{
		Type:        ROLE,
		Name:        name,
		Description: description,
		Required:    required,
	}
}

func Group(name, description string, options []ApplicationCommandOption) ApplicationCommandOption {
	return ApplicationCommandOption{
		Type:        SUB_COMMAND_GROUP,
		Name:        name,
		Description: description,
		Options:     options,
	}
}

func SubCmd(name, description string, options []ApplicationCommandOption) ApplicationCommandOption {
	return ApplicationCommandOption{
		Type:        SUB_COMMAND,
		Name:        name,
		Description: description,
		Options:     options,
	}
}
func Int(name, description string, required bool) ApplicationCommandOption {
	return ApplicationCommandOption{
		Type:        INTEGER,
		Name:        name,
		Description: description,
		Default:     false,
		Required:    required,
	}
}

func Str(name, description string, required bool) ApplicationCommandOption {
	return ApplicationCommandOption{
		Type:        STRING,
		Name:        name,
		Description: description,
		Default:     false,
		Required:    required,
	}
}

type ApplicationCommand struct {
	ID            string                     `json:"id"`
	ApplicationID string                     `json:"ApplicationID"`
	Name          string                     `json:"name"`
	Description   string                     `json:"description"`
	Options       []ApplicationCommandOption `json:"options"`
}

type ApplicationCommandOption struct {
	Type        ApplicationCommandOptionType     `json:"type"`
	Name        string                           `json:"name"`
	Description string                           `json:"description"`
	Default     bool                             `json:"default"`
	Required    bool                             `json:"required"`
	Choices     []ApplicationCommandOptionChoice `json:"choices"`
	Options     []ApplicationCommandOption       `json:"options"`
}

type ApplicationCommandOptionChoice struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

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

type InteractionResponseType int

const (
	PONG InteractionResponseType = iota + 1
	ACK
	CHANNEL_MESSAGE
	CHANNEL_MESSAGE_WITH_SOURCE
	ACK_WITH_SOURCE
)

func Response(content string) *InteractionResponse {
	responseType := CHANNEL_MESSAGE
	return &InteractionResponse{
		Type: responseType,
		Data: &InteractionApplicationCommandCallbackData{
			Content: content,
		},
	}
}

func Acknowledge() *InteractionResponse {
	return &InteractionResponse{
		Type: ACK,
		Data: nil,
	}
}

func (resp *InteractionResponse) Embed(embed discordgo.MessageEmbed) *InteractionResponse {
	resp.Data.Embeds = append(resp.Data.Embeds, embed)
	return resp
}

func (resp *InteractionResponse) KeepSource() *InteractionResponse {
	if resp.Type == CHANNEL_MESSAGE {
		resp.Type = CHANNEL_MESSAGE_WITH_SOURCE
		return resp
	}

	resp.Type = ACK_WITH_SOURCE
	return resp
}

func (resp *InteractionResponse) OnlyAuthor() *InteractionResponse {
	resp.Data.Flags = 1 << 6
	return resp
}

type InteractionResponse struct {
	Type InteractionResponseType                    `json:"type"`
	Data *InteractionApplicationCommandCallbackData `json:"data,omitempty"`
}

type InteractionApplicationCommandCallbackData struct {
	TTS             bool                               `json:"tts"`
	Content         string                             `json:"content"`
	Embeds          []discordgo.MessageEmbed           `json:"embeds,omitempty"`
	AllowedMentions []discordgo.MessageAllowedMentions `json:"allowed_mentions,omitempty"`
	Flags           int                                `json:"flags"`
}
