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
}

type Config struct {
    Auth    Auth
    General struct {
        Timeout    time.Duration
        UpdateTime time.Duration
        DB         string
    }
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
}

func GetConfig() *Config {
    if cfg == nil {
        loadConfig()
    }

    return cfg
}
