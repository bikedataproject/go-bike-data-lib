package strava

type StravaWebhookMessage struct {
	ObjectType     string `json:"object_type"`
	ObjectID       int32  `json:"object_id"`
	AspectType     string `json:"aspect_type"`
	OwnerID        int32  `json:"owner_id"`
	SubscriptionID int    `json:"subscription_id"`
	EventTime      int32  `json:"event_time"`
	// Using interface{} for Updates since it can be either a single update or an array of updates
	Updates interface{} `json:"updates"`
}

type StravaUpdateMessage struct {
	Title      string `json:"title"`
	Type       string `json:"type"`
	Private    bool   `json:"private"`
	Authorized bool   `json:"authorized"`
}

type StravaSubscribeRequest struct {
	ClientID     int    `json:"client_id"`
	ClientSecret string `json:"client_secret"`
	CallbackURL  string `json:"callback_url"`
	VerifyToken  string `json:"verify_token"`
}
