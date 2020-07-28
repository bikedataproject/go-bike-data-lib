# go-bike-data-lib 

<p align="center">
  <a href="https://github.com/bikedataproject/go-bike-data-lib">
    <img src="https://avatars3.githubusercontent.com/u/64870976?s=200&v=4" alt="Logo" width="80" height="80">
  </a>

  <h3 align="center">Go Bike Data Library</h3>

  <p align="center">
    This repository is used as a Go library used by our other Go services.
    <br />
    <a href="https://github.com/bikedataproject/go-bike-data-lib/issues">Report Bug</a>
    ·
    <a href="https://github.com/bikedataproject/go-bike-data-lib/issues">Request Feature</a>
  </p>
</p>

## About this repository

This repository is used as a Go library for the other Go services. It contains multiple packages, which separate functionalities or data structures required for certain services.

## Package: strava

This package contains datamodels for Strava services. These are mostly Structs representing JSON-data as a Go-native object.

The datamodels that have been implemented are:

- **Strava Activities**: an activity created by a certain user
  - Extended with ActivityMap: important information such as the activity PolyLine
- **Strava WebhookMessage**: a message sent by the Strava webhook service to inform us that there have been updates
- **Strava Webhook Validation Request**: a message sent by Strava to our service with a field named `hub.challenge` which has to be returned within 2 seconds to verify the functionality of the webhook listener.
- **Strava Refresh Message**: a response message from Strava returning important data in regards of access to a certain user his/her data
- **Strava Subscription Message**: a response message from Strava returning the ID of a webhook subscription on success

## Package: dbmodel

This package contains datamodels and functions to interact with the Postgres database.

The datamodels that have been implemented are:

- **User**: a single User object as stored in the database which stores access data and identifiers
- **UserContribution**: a single UserContribution object as stored in the database which defines the relation between a User and a Contribution
- **Contribution**: a single Contribution object as stored in the database which stores activity data

The database functionalities that have been implemented are:

- **Get connectionstring**: generates a connectionstring to be able to connect to the database
- **Check connection**: pings the database and verifies the connection
- **Get user data**: request data about a user with the UserID
- **Add user**: create a new user entry in the database
- **Add contribution**: create a new contribution in the database
- **Get expiring users**: Strava users have a limited access token duration (6 hours) and should be refreshed before timing out. This function fetches all users that are expiring in the next hour.
- **Update user**: modify an existing user with new data
- **Fetch new users**: new Strava users can have a lot of activity history which has not yet been fetched. The Strava daemon will fetch the activity history based on this function.
