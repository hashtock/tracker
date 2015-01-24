package main

import (
    "log"
    "os"
    "time"

    "code.google.com/p/gcfg"
)

type auth struct {
    ConsumerKey       string
    SecretKey         string
    AccessToken       string
    AccessTokenSecret string
}

type config struct {
    Auth    auth
    General struct {
        Timeout time.Duration
        DB      string
    }
}

var cfg config

func init() {
    err := gcfg.ReadFileInto(&cfg, "config.ini")
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

    if err := startSession(cfg.General.DB); err != nil {
        log.Fatalln("Could not connect to DB.", err.Error())
    }
}

func main() {
    app := CliApp()
    app.Run(os.Args)
}
