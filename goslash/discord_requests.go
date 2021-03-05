package goslash

import (
	"encoding/json"
	"fmt"
	log "github.com/sirupsen/logrus"
)

const DISCORD_API_BASE_URL = "https://discord.com/api/v8"

func (app Application) PostJson(path string, obj interface{}) ([]byte, error) {
	url := fmt.Sprintf("%v/%v", DISCORD_API_BASE_URL, path)

	return app.Session.Request("POST", url, obj)
}

func (app *Application) AddCommand(command *Command) {
	app.Commands[command.Name] = command
}

func (app *Application) RegisterGlobal(command *Command) (*Command, error) {
	data, err := app.PostJson(fmt.Sprintf("applications/%v/commands", app.ClientID), command)
	if err != nil {
		log.WithField("error", err).Info("Error occurred while registering global command")
		return nil, err
	}

	err = json.Unmarshal(data, &command.ApplicationCommand)
	app.Commands[command.Name] = command
	command.isGlobal = true
	if err != nil {
		log.WithField("error", err).Info("Error occurred while registering global command")
	}
	return command, err
}

func (app *Application) RegisterGuild(guildid string, command *Command) (*Command, error) {
	data, err := app.PostJson(fmt.Sprintf("applications/%v/guilds/%v/commands", app.ClientID, guildid), command)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(data, &command.ApplicationCommand)
	return command, err
}

func (app *Application) RegisterAllGuild(guildid string) error {
	for _, command := range app.Commands {
		_, err := app.RegisterGuild(guildid, command)
		if err != nil {
			return err
		}
	}
	return nil
}

func (app *Application) RegisterAllGlobal() error {
	for _, command := range app.Commands {
		_, err := app.RegisterGlobal(command)
		if err != nil {
			return err
		}
	}
	return nil
}
