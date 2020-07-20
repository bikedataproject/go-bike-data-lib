package strava

// SubscriptionMessage : Struct that holds the ID of an individual webhook subscription
type SubscriptionMessage struct {
	ID int `json:"id"`
}

// RefreshMessage : Struct that holds the response when refreshing strava access
type RefreshMessage struct {
	TokenType    string `json:"token_type"`
	AccessToken  string `json:"access_token"`
	ExpiresAt    int    `json:"expires_at"`
	ExpiresIn    int    `json:"expires_in"`
	RefreshToken string `json:"refresh_token"`
}

// WebhookValidationRequest : Body of the incoming GET request to verify the endpoint
type WebhookValidationRequest struct {
	HubChallenge string `json:"hub.challenge"`
}
