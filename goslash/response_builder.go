// helper functions to make building interaction responses easier

package goslash

import "github.com/bwmarrin/discordgo"

type InteractionResponse discordgo.InteractionResponse

func Pong() *InteractionResponse {
	return &InteractionResponse{
		Type: discordgo.InteractionResponsePong,
	}
}

func Response(content string) *InteractionResponse {
	return &InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionApplicationCommandResponseData{
			Content:         content,
			Flags:           0,
		},
	}
}

func Deferred() *InteractionResponse {
	return &InteractionResponse{
		Type: discordgo.InteractionResponseDeferredChannelMessageWithSource,
	}
}

func (response *InteractionResponse) Embed(embed discordgo.MessageEmbed) *InteractionResponse {
	if response.Data == nil {
		response.Data = &discordgo.InteractionApplicationCommandResponseData{}
	}
	response.Data.Embeds = append(response.Data.Embeds, &embed)
	return response
}

func (response *InteractionResponse) Ephemeral() *InteractionResponse {
	response.Data.Flags = response.Data.Flags | (1 << 6)

	return response
}

func (response *InteractionResponse) ToDiscordgo() *discordgo.InteractionResponse {
	if response == nil {
		return nil
	}

	newResp := discordgo.InteractionResponse(*response)
	return &newResp

}