package strava

import (
	"time"

	geo "github.com/paulmach/go.geo"
)

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

// WebhookMessage : Body of incoming webhook messages
type WebhookMessage struct {
	ObjectType     string      `json:"object_type"`
	ObjectID       int         `json:"object_id"`
	AspectType     string      `json:"aspect_type"`
	OwnerID        int         `json:"owner_id"`
	SubscriptionID int         `json:"subscription_id"`
	EventTime      int         `json:"event_time"`
	Updates        interface{} `json:"updates"`
}

// Activity : Struct representing an activity from Strava
type Activity struct {
	Distance           float32   `json:"distance"`
	MovingTime         int       `json:"moving_time"`
	ElapsedTime        int       `json:"elapsed_time"`
	TotalElevationGain float64   `json:"total_elevation_gain"`
	Type               string    `json:"type"`
	WorkoutType        int       `json:"workout_type"`
	StartDateLocal     time.Time `json:"start_date_local"`
	EndDateLocal       time.Time
	PointsTime         []time.Time
	StartLatlng        []float64   `json:"start_latlng"`
	EndLatlng          []float64   `json:"end_latlng"`
	Map                ActivityMap `json:"map"`
	Commute            bool        `json:"commute"`
	LineString         *geo.Path
}

// ActivityMap : Struct representing the Map field in an activity message
type ActivityMap struct {
	ID              string `json:"id"`
	Polyline        string `json:"polyline"`
	ResourceState   int    `json:"resource_state"`
	SummaryPolyline string `json:"summary_polyline"`
}
