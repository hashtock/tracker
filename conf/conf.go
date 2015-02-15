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
    Timeout     string
    UpdateTime  string
    SampingTime string
    DB          string
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

var cfg *Config = nil
var rcfgs RemoteConfigs = nil

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

func loadRemoteConfigs() {
    rcfgs = make(RemoteConfigs, 0)

    example_config := `[remote "host1"]
URL = "www.tracker.com:80"
HMACSecret = "shared secret with host1"
    `

    tmp_rcfgs := struct {
        Remote map[string]*RemoteConfig
    }{}

    err := gcfg.ReadFileInto(&tmp_rcfgs, "remotes.ini")
    if err != nil {
        if os.IsNotExist(err) {
            log.Fatalf("Could not find remote tracker configuration. Expected remotes.ini with content:\n%v", example_config)
        } else {
            log.Fatalln("Remote config error:", err.Error())
        }
    }

    for key, config := range tmp_rcfgs.Remote {
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
