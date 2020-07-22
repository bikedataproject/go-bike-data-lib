package dbmodel

import (
	"time"

	geo "github.com/paulmach/go.geo"
)

// User : Struct to respresent a user object from the database
type User struct {
	ID                string
	UserIdentifier    string
	Provider          string
	ProviderUser      string
	AccessToken       string
	RefreshToken      string
	TokenCreationDate time.Time
	ExpiresAt         int
	ExpiresIn         int
}

// Contribution : Struct to respresent a contribution object from the database
type Contribution struct {
	ContributionID string
	UserAgent      string
	Distance       int
	TimeStampStart time.Time
	TimeStampStop  time.Time
	Duration       int
	PointsGeom     *geo.Path
	PointsTime     []time.Time
}

// UserContribution : Struct to respresent a UserContribution object from the database
type UserContribution struct {
	UserContributionID string
	UserID             string
	ContributionID     string
}
