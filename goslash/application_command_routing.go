package goslash

import "github.com/bwmarrin/discordgo"

func (app *Application) HandleInteraction(interaction *discordgo.Interaction) *discordgo.InteractionResponse {
	resp := app.getResponse(interaction)
	if resp == nil {
		resp = app.DefaultResponse
	}

	return resp.ToDiscordgo()
}

func (app *Application) getResponse(interaction *discordgo.Interaction) *InteractionResponse {

	rootCommand := app.GetCommand(interaction.Data.Name)
	if rootCommand == nil {
		return nil
	}

	ctx := &InteractionUpdate{interaction, app, []string{rootCommand.Name}}

	if resp := rootCommand.Handle(ctx); resp != nil {
		return resp
	}

	return handleOptions(rootCommand, ctx, interaction.Data.Options)
}

func handleOptions(rootCommand *Command, ctx *InteractionUpdate, options []*discordgo.ApplicationCommandInteractionDataOption) *InteractionResponse {
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
