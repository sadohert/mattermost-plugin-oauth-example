package main

import "golang.org/x/oauth2"

type ExternalServiceClient struct {
	token *oauth2.Token
}

func (s *ExternalServiceClient) GetMe() (string, error) {
	return "some user id for token " + s.token.AccessToken, nil
}

func (s *ExternalServiceClient) CreateTicket() (string, error) {
	// This function would be defined in an external library. Magic happens here.
	return "ticket-12345", nil
}

func makeClientFromExternalIntegration(token *oauth2.Token) *ExternalServiceClient {
	return &ExternalServiceClient{token}
}

func (p *Plugin) getExternalClient(userID string) (*ExternalServiceClient, error) {
	token, err := p.OAuther.GetToken(userID)
	if err != nil {
		return nil, err
	}

	client := makeClientFromExternalIntegration(token)
	return client, nil
}

func getOAuthScopesForExternalService() []string {
	return []string{"user:read", "user:write"}
}
