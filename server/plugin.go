package main

import (
	"fmt"
	"net/http"
	"strings"
	"sync"

	pluginapi "github.com/mattermost/mattermost-plugin-api"
	"github.com/mattermost/mattermost-plugin-api/experimental/bot/logger"
	"github.com/mattermost/mattermost-plugin-api/experimental/oauther"
	"github.com/mattermost/mattermost-server/v6/model"
	"github.com/mattermost/mattermost-server/v6/plugin"
)

// Plugin implements the interface expected by the Mattermost server to communicate between the server and plugin processes.
type Plugin struct {
	plugin.MattermostPlugin

	// configurationLock synchronizes access to the configuration.
	configurationLock sync.RWMutex

	// configuration is the active plugin configuration. Consult getConfiguration and
	// setConfiguration for usage.
	configuration *configuration

	client  *pluginapi.Client
	OAuther oauther.OAuther

	botUserID string
}

func (p *Plugin) OnActivate() error {
	pluginAPIClient := pluginapi.NewClient(p.API, p.Driver)
	p.client = pluginAPIClient

	p.OAuther = oauther.NewFromClient(
		pluginAPIClient,
		p.getOAuthConfig(),
		p.onConnect,
		logger.New(p.API),
	)

	p.client.SlashCommand.Register(&model.Command{
		Trigger:      "oauth-example",
		AutoComplete: true,
		AutocompleteData: &model.AutocompleteData{
			Trigger: "oauth-example",
			SubCommands: []*model.AutocompleteData{
				{
					Trigger: "connect",
				},
				{
					Trigger: "create",
				},
			},
		},
	})

	botUserID, err := p.client.Bot.EnsureBot(&model.Bot{
		Username: "oauth-bot",
	})
	if err != nil {
		return err
	}

	p.botUserID = botUserID

	return nil
}

// ServeHTTP demonstrates a plugin that handles HTTP requests by greeting the world.
func (p *Plugin) ServeHTTP(c *plugin.Context, w http.ResponseWriter, r *http.Request) {
	if strings.HasPrefix(r.URL.Path, "/oauth2") {
		p.OAuther.ServeHTTP(w, r)
		return
	}

	fmt.Fprint(w, "Hello, world!")
}

// See https://developers.mattermost.com/extend/plugins/server/reference/
