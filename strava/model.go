package strava

// WebhookMessage : Data sent from a Strava webhook message
type WebhookMessage struct {
	ObjectType     string `json:"object_type"`
	ObjectID       int32  `json:"object_id"`
	AspectType     string `json:"aspect_type"`
	OwnerID        int32  `json:"owner_id"`
	SubscriptionID int    `json:"subscription_id"`
	EventTime      int32  `json:"event_time"`
	// Using interface{} for Updates since it can be either a single update or an array of updates
	Updates interface{} `json:"updates"`
}

// UpdateMessage : Data sent from Stava about updates through webhooks
type UpdateMessage struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	Private    bool   `json:"private"`
	Authorized bool   `json:"authorized"`
}

// SubscribeRequest : Data model to send to Strava in order to create a webhook subscription
type SubscribeRequest struct {
	ClientID     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
	VerifyToken  string `json:"verify_token"`
}
