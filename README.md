# Tracker

Service tracking occurances of hashtags.
Currently only on Twitter.

## Requirenments

- [Golang](https://golang.org/)
- [MongoDB](https://www.mongodb.org/)
- [Twitter API key](https://apps.twitter.com/)

## Instalation

To build from source:

```bash
go get github.com/hashtock/tracker
```

## Usage

Run `./tracker` to listen for new hashtags and serve web app

## Configuration

Tracker expects to read configuration from environment vairables.
Running Trucker without mandatory configuration value or executing `config` command (`tracker config`) will cause it to print out help message like:

```
Environmental variables used in configuration
TRACKER_DB
    Value: mongodb://admin:123456@1.2.3.4:27017/
    Help: Location of MongoDB: mongodb://user:password@host:port/

TRACKER_TIMEOUT
    Value: 0 (default)
    Help: How long to listen for, 0 for inifinite

TRACKER_UPDATE_TIME
    Value: 5s (default)
    Help: How often push new counts to DB

TRACKER_SAMPING_TIME
    Value: 15m0s (default)
    Help: Store counts grouped by time

TRACKER_TAG_UPDATE_TIME
    Value: 10s (default)
    Help: How often to check for new tags while listening

TRACKER_CONSUMER_KEY
    Value: not set
    Help: Twitter App ConsumerKey

TRACKER_SECRET_KEY
    Value: not set
    Help: Twitter App SecretKey

TRACKER_ACCESS_TOKEN
    Value: not set
    Help: Twitter account access token

TRACKER_ACCESS_TOKEN_SECRET
    Value: not set
    Help: Twitter account access token secret

TRACKER_SECRET
    Value: not set
    Help: Long random string used as shared secret
```

