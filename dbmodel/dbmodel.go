package dbmodel

import (
	"database/sql"
	"fmt"
	"time"

	"github.com/lib/pq"
	geo "github.com/paulmach/go.geo"
	log "github.com/sirupsen/logrus"
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
	IsHistoryFetched  bool
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

// Database : Struct to hold the database connection
type Database struct {
	PostgresHost       string
	PostgresUser       string
	PostgresPassword   string
	PostgresPort       int64
	PostgresDb         string
	PostgresRequireSSL string
}

// getDBConnectionString : Generate connectionstring
func (db Database) getDBConnectionString() string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=%v", db.PostgresHost, db.PostgresPort, db.PostgresUser, db.PostgresPassword, db.PostgresDb, db.PostgresRequireSSL)
}

// checkConnection : Check if the database can be reached
func (db Database) checkConnection() bool {
	// Connect to database
	connection, err := sql.Open("postgres", db.getDBConnectionString())
	if err != nil {
		return false
	}
	// Ping database
	err = connection.Ping()
	return err == nil
}

// VerifyConnection : Connect to Postgres
func (db Database) VerifyConnection() (err error) {
	if db.checkConnection() {
		log.Info("Database is reachable")
	} else {
		log.Fatal("Database is unreachable")
	}
	return
}

// GetUserData : Request a user token for the ID
func (db Database) GetUserData(userID string) (usr User, err error) {
	// Connect to database
	connection, err := sql.Open("postgres", db.getDBConnectionString())
	if err != nil {
		defer connection.Close()
		return
	}

	// Query user
	response := connection.QueryRow(`
	SELECT "Id", "UserIdentifier", "Provider", "ProviderUser", "AccessToken", "RefreshToken", "TokenCreationDate", "ExpiresAt", "ExpiresIn", "IsHistoryFetched"
	FROM "Users"
	WHERE "ProviderUser"=$1;
	`, userID)

	// Load data in struct
	err = response.Scan(&usr.ID, &usr.UserIdentifier, &usr.Provider, &usr.ProviderUser, &usr.AccessToken, &usr.RefreshToken, &usr.TokenCreationDate, &usr.ExpiresAt, &usr.ExpiresIn, &usr.IsHistoryFetched)
	defer connection.Close()
	return
}

// AddUser : Create new user in the database
func (db Database) AddUser(user *User) (newUser User, err error) {
	// Connect to database
	connection, err := sql.Open("postgres", db.getDBConnectionString())
	if err != nil {
		defer connection.Close()
		err = fmt.Errorf("Could not create database connection: %v", err)
		return
	}

	// Create new user & fetch ID
	query := `
	INSERT INTO "Users"
	("UserIdentifier", "Provider", "ProviderUser", "TokenCreationDate", "ExpiresAt", "ExpiresIn", "IsHistoryFetched")
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING "Id";
	`
	response := connection.QueryRow(query, &user.UserIdentifier, &user.Provider, &user.ProviderUser, &user.TokenCreationDate, &user.ExpiresAt, &user.ExpiresIn, &user.IsHistoryFetched)

	// Load data in struct
	err = response.Scan(&newUser.ID)
	if err != nil {
		return
	}

	// Load old data
	newUser.ExpiresAt = user.ExpiresAt
	newUser.ExpiresIn = user.ExpiresIn
	newUser.ProviderUser = user.ProviderUser
	newUser.Provider = user.Provider
	newUser.UserIdentifier = user.UserIdentifier

	return
}

// AddContribution : Create new user contribution
func (db Database) AddContribution(contribution *Contribution, user *User) (err error) {
	// Connect to database
	connection, err := sql.Open("postgres", db.getDBConnectionString())
	if err != nil {
		defer connection.Close()
		return fmt.Errorf("Could not create database connection: %v", err)
	}

	// Write Contribution
	query := `
	INSERT INTO "Contributions"
	("UserAgent", "Distance", "TimeStampStart", "TimeStampStop", "Duration", "PointsGeom", "PointsTime")
	VALUES ($1, $2, $3, $4, $5, $6, $7)
	RETURNING "ContributionId";
	`
	response := connection.QueryRow(query, contribution.UserAgent, contribution.Distance, contribution.TimeStampStart, contribution.TimeStampStop, contribution.Duration, contribution.PointsGeom.ToWKT(), pq.Array(contribution.PointsTime))
	defer connection.Close()

	// Create contributions
	userContrib := UserContribution{
		UserID: user.ID,
	}

	err = response.Scan(&userContrib.ContributionID)
	if err != nil {
		err = fmt.Errorf("Could not extract contributionID: %v", err)
		return
	}

	// Write WriteUserContribution
	query = `
	INSERT INTO "UserContributions"
	("UserId", "ContributionId")
	VALUES ($1, $2);
	`
	if _, err = connection.Exec(query, userContrib.UserID, &userContrib.ContributionID); err != nil {
		defer connection.Close()
		return fmt.Errorf("Could not insert value into contributions: %s", err)
	}

	defer connection.Close()
	return
}

// GetExpiringUsers : Get users which are expiring within half an hour
func (db Database) GetExpiringUsers() (users []User, err error) {
	// Connect to database
	connection, err := sql.Open("postgres", db.getDBConnectionString())
	if err != nil {
		defer connection.Close()
		return
	}

	// Fetch expiring users
	response, err := connection.Query(`
	SELECT "Id", "RefreshToken", "UserIdentifier" FROM "Users"
	WHERE "ExpiresAt" <= $1 and "Provider" = 'web/Strava';
	`, time.Now().Add(30*time.Minute).Unix())
	if err != nil {
		defer connection.Close()
		return
	}

	// Convert sql.Rows into User objects
	for response.Next() {
		var user User
		err = response.Scan(&user.ID, &user.RefreshToken, &user.UserIdentifier)
		if err != nil {
			log.Warnf("Could not add expiring user to result: %v", err)
		}
		users = append(users, user)
	}

	defer connection.Close()
	return
}

// UpdateUser : Update an existing user
func (db Database) UpdateUser(user *User) (err error) {
	// Connect to database
	connection, err := sql.Open("postgres", db.getDBConnectionString())
	if err != nil {
		defer connection.Close()
		return
	}

	// Update user in database
	_, err = connection.Exec(`
	UPDATE "Users"
	SET "ExpiresAt" = $1,
		"ExpiresIn" = $2,
		"AccessToken" = $3,
		"RefreshToken" = $4,
		"IsHistoryFetched" = $5
	WHERE "UserIdentifier" = $6;
	`, &user.ExpiresAt, &user.ExpiresIn, &user.AccessToken, &user.RefreshToken, &user.IsHistoryFetched, &user.UserIdentifier)

	defer connection.Close()
	return
}

// FetchNewUsers : Fetch an array of new users that have not yet fetched their old data
func (db Database) FetchNewUsers() (users []User, err error) {
	// Connect to database
	connection, err := sql.Open("postgres", db.getDBConnectionString())
	if err != nil {
		defer connection.Close()
		return
	}

	// Fetch new users
	response, err := connection.Query(`
	SELECT "Id", "UserIdentifier", "Provider", "ProviderUser", "AccessToken", "RefreshToken", "TokenCreationDate", "ExpiresAt", "ExpiresIn", "IsHistoryFetched"
	FROM "Users"
	WHERE "Provider" = 'web/Strava'
	AND NOT "IsHistoryFetched";
	`)
	if err != nil {
		defer connection.Close()
		return
	}

	// Convert sql.Rows into User objects
	for response.Next() {
		var user User
		err = response.Scan(&user.ID, &user.UserIdentifier, &user.Provider, &user.ProviderUser, &user.AccessToken, &user.RefreshToken, &user.TokenCreationDate, &user.ExpiresAt, &user.ExpiresIn, &user.IsHistoryFetched)
		if err != nil {
			log.Warnf("Could not add expiring user to result: %v", err)
		}
		users = append(users, user)
	}

	defer connection.Close()
	return
}
