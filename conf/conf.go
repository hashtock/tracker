package conf

import (
	"log"
	"os"
	"time"

	"code.google.com/p/gcfg"
)

type Auth struct {
	ConsumerKey       string
	SecretKey         string
	AccessToken       string
	AccessTokenSecret string
	HMACSecret        string
}

type General struct {
	Timeout       string
	UpdateTime    string
	TagUpdateTime string
	SampingTime   string
	DB            string
}

type Config struct {
	Auth    Auth
	General General
}

type RemoteConfig struct {
	URL        string
	HMACSecret string
}

type RemoteConfigs map[string]RemoteConfig

var cfg *Config
var rcfgs RemoteConfigs

const exampleConfig = `[general]
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
`

const exampleRemoteConfig = `[remote "host1"]
URL = "www.tracker.com:80"
HMACSecret = "shared secret with host1"
`

func (r RemoteConfigs) names() []string {
	names := make([]string, 0, len(r))
	for name := range r {
		names = append(names, name)
	}
	return names
}

func loadConfig() {
	if cfg == nil {
		cfg = new(Config)
	}
	err := gcfg.ReadFileInto(cfg, "config.ini")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Could not find remote tracker configuration. Expected config.ini with content:\n%v\n", exampleConfig)
		} else {
			log.Fatalln("Config error:", err.Error())
		}
	}

	if cfg.Auth.ConsumerKey == "" || cfg.Auth.SecretKey == "" || cfg.Auth.AccessToken == "" || cfg.Auth.AccessTokenSecret == "" {
		log.Fatalln("Twitter authentication missing!\nExpect:", exampleConfig)
	}

	cfg.General.validate()
}

func loadRemoteConfigs() {
	rcfgs = make(RemoteConfigs, 0)

	tmpRCfgs := struct {
		Remote map[string]*RemoteConfig
	}{}

	err := gcfg.ReadFileInto(&tmpRCfgs, "remotes.ini")
	if err != nil {
		if os.IsNotExist(err) {
			log.Fatalf("Could not find remote tracker configuration. Expected remotes.ini with content:\n%v\n", exampleRemoteConfig)
		} else {
			log.Fatalln("Remote config error:", err.Error())
		}
	}

	for key, config := range tmpRCfgs.Remote {
		rcfgs[key] = *config
	}
}

func GetConfig() *Config {
	if cfg == nil {
		loadConfig()
	}

	return cfg
}

func ListRemoteConfigs() []string {
	if rcfgs == nil {
		loadRemoteConfigs()
	}

	return rcfgs.names()
}

func GetRemoteConfig(remote string) RemoteConfig {
	if rcfgs == nil {
		loadRemoteConfigs()
	}

	config, ok := rcfgs[remote]
	if !ok {
		log.Fatalf("Could not find config configuration for: %v. Available configurations: %v", remote, rcfgs.names())
	}
	return config
}

func parseOrDie(durationStr string) time.Duration {
	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		log.Fatalf("Could not parse %#v as duration. Expected config: \n%s\n%v", durationStr, exampleConfig, err)
	}
	return duration
}

func (g *General) validate() {
	g.TimeoutD()
	g.UpdateTimeD()
	g.SampingTimeD()
	g.TagUpdateTimeD()
}

func (g *General) TimeoutD() time.Duration {
	return parseOrDie(g.Timeout)
}

func (g *General) UpdateTimeD() time.Duration {
	return parseOrDie(g.UpdateTime)
}

func (g *General) SampingTimeD() time.Duration {
	return parseOrDie(g.SampingTime)
}

func (g *General) TagUpdateTimeD() time.Duration {
	return parseOrDie(g.TagUpdateTime)
}
