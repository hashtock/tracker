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

To get list of available commands run `tracker --help`

Basic use case would be to:
* Add some tags to track ```tracker add android golang mognodb```
* Listen
  * No web API ```tracker listen```
  * With web API ```tracker web```
* Lookup counts (lists only completed samples. See SampingTime value in config)
  * Sum for time period ```tacker counts 1h```
  * Counts over time for time period ```tacker counts_detailed 2h```

## Configuration

Tracker expects `config.ini` file in place from where you will execute it.
Running Trucker without config will cause it to print out examplary config, like one below:

```ini
[general]
Timeout = 60s ; How long to listen for, 0 for inifinite
UpdateTime = 5s ; How often push new counts to DB
SampingTime = 15m ; Store counts grouped by time
TagUpdateTime = 1m ; How often to check for new tags while listening
DB = "mongodb://user:password@host:port/"

[auth]
ConsumerKey = "Twitter App ConsumerKey"
SecretKey   = "Twitter App SecretKey"
AccessToken = "Twitter account access token"
AccessTokenSecret = "Twitter account access token secret"
HMACSecret = "Long Random String"
```

