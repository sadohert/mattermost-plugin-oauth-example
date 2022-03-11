package main

import (
	"github.com/mattermost/mattermost-server/v6/model"
	"golang.org/x/oauth2"
)

func (p *Plugin) onConnect(userID string, token oauth2.Token, payload []byte) {
	client := makeClientFromExternalIntegration(&token)

	externalUserID, err := client.GetMe()
	if err != nil {
		p.client.Log.Error(err.Error())
		return
	}

	keyToStoreExternalUserInfo := userID + "_external"
	externalUserIDBytes := []byte(externalUserID)
	p.client.KV.Set(keyToStoreExternalUserInfo, externalUserIDBytes)

	p.client.Post.DM(p.botUserID, userID, &model.Post{
		Message: "Hey you're connected now!",
	})
}

func (p *Plugin) getOAuthConfig() oauth2.Config {
	conf := p.getConfiguration()

	// Test server run with https://www.npmjs.com/package/oauth2-mock-server
	// Replace this with the appropriate values from your external integration
	authURL := "http://localhost:8080/authorize"
	tokenURL := "http://localhost:8080/token"

	// Arbitrary scopes supported by external service
	scopes := getOAuthScopesForExternalService()

	return oauth2.Config{
		ClientID:     conf.OAuth2ClientId,
		ClientSecret: conf.OAuth2ClientSecret,
		Scopes:       scopes,
		Endpoint: oauth2.Endpoint{
			AuthURL:  authURL,
			TokenURL: tokenURL,
		},
	}
}
