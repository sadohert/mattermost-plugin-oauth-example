package main

import (
	"fmt"
	"strings"

	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

func (p *Plugin) ExecuteCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	parts := strings.Fields(args.Command)
	if len(parts) < 2 {
		return &model.CommandResponse{
			Text: "Please provide a subcommand",
		}, nil
	}

	subcommand := parts[1]
	switch subcommand {
	case "connect":
		return p.runConnectCommand(c, args)
	case "create":
		return p.runCreateCommand(c, args)
	}

	return &model.CommandResponse{
		Text: "Invalid subcommand",
	}, nil
}

func (p *Plugin) runConnectCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	text := fmt.Sprintf("Click [here](%s) to connect your account", p.OAuther.GetConnectURL())

	return &model.CommandResponse{
		Text: text,
	}, nil
}

func (p *Plugin) runCreateCommand(c *plugin.Context, args *model.CommandArgs) (*model.CommandResponse, *model.AppError) {
	externalClient, err := p.getExternalClient(args.UserId)
	if err != nil {
		return &model.CommandResponse{
			Text: "Failed to get external client: " + err.Error(),
		}, nil
	}

	ticketID, err := externalClient.CreateTicket()
	if err != nil {
		return &model.CommandResponse{
			Text: "Failed to get create ticket: " + err.Error(),
		}, nil
	}

	return &model.CommandResponse{
		Text: fmt.Sprintf("Created a new ticket `%s`", ticketID),
	}, nil
}
