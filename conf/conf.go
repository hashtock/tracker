package conf

import (
	"log"
	"time"
)

type Auth struct {
	ConsumerKey       string
	SecretKey         string
	AccessToken       string
	AccessTokenSecret string
	HMACSecret        string
}

type General struct {
	Timeout       time.Duration
	UpdateTime    time.Duration
	TagUpdateTime time.Duration
	SampingTime   time.Duration
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

func (r RemoteConfigs) names() []string {
	names := make([]string, 0, len(r))
	for name := range r {
		names = append(names, name)
	}
	return names
}

func (r *RemoteConfig) valid() bool {
	return r.URL != "" && r.HMACSecret != ""
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
