package conf

import (
    "log"
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
    Timeout     string
    UpdateTime  string
    SampingTime string
    DB          string
}

type Config struct {
    Auth    Auth
    General General
}

var cfg *Config = nil

func loadConfig() {
    if cfg == nil {
        cfg = new(Config)
    }
    err := gcfg.ReadFileInto(cfg, "config.ini")
    if err != nil {
        log.Fatalln("Config error:", err.Error())
    }

    example_config := `[auth]
        ConsumerKey = "123"
        SecretKey   = "456"
        AccessToken = "679"
        AccessTokenSecret = "001"
    `

    if cfg.Auth.ConsumerKey == "" || cfg.Auth.SecretKey == "" || cfg.Auth.AccessToken == "" || cfg.Auth.AccessTokenSecret == "" {
        log.Fatalln("Twitter authentication missing!\nExpect:", example_config)
    }

    cfg.General.validate()
}

func GetConfig() *Config {
    if cfg == nil {
        loadConfig()
    }

    return cfg
}

func parseOrDie(duration_str string) time.Duration {
    duration, err := time.ParseDuration(duration_str)
    if err != nil {
        log.Fatal(err)
    }
    return duration
}

func (g *General) validate() {
    g.TimeoutD()
    g.UpdateTimeD()
    g.SampingTimeD()
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
