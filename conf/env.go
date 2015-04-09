package conf

import (
	"fmt"
	"os"
	"regexp"
	"strings"
	"time"
)

const (
	ENV_PREFIX = "TRACKER_"
)

const (
	keyDB                  = "DB"
	keyTIMEOUT             = "TIMEOUT"
	keyUPDATE_TIME         = "UPDATE_TIME"
	keySAMPING_TIME        = "SAMPING_TIME"
	keyTAG_UPDATE_TIME     = "TAG_UPDATE_TIME"
	keyACCESS_TOKEN        = "ACCESS_TOKEN"
	keyACCESS_TOKEN_SECRET = "ACCESS_TOKEN_SECRET"
	keyCONSUMER_KEY        = "CONSUMER_KEY"
	keySECRET_KEY          = "SECRET_KEY"
	keySECRET              = "SECRET"
)

var cfgHelp = map[string]string{
	keyDB:                  "Location of MongoDB: mongodb://user:password@host:port/",
	keyTIMEOUT:             "How long to listen for, 0 for inifinite",
	keyUPDATE_TIME:         "How often push new counts to DB",
	keySAMPING_TIME:        "Store counts grouped by time",
	keyTAG_UPDATE_TIME:     "How often to check for new tags while listening",
	keyCONSUMER_KEY:        "Twitter App ConsumerKey",
	keySECRET_KEY:          "Twitter App SecretKey",
	keyACCESS_TOKEN:        "Twitter account access token",
	keyACCESS_TOKEN_SECRET: "Twitter account access token secret",
	keySECRET:              "Long random string used as shared secret",
}

var defaultDurations = map[string]time.Duration{
	keyTIMEOUT:         0,
	keyUPDATE_TIME:     5 * time.Second,
	keySAMPING_TIME:    15 * time.Minute,
	keyTAG_UPDATE_TIME: 10 * time.Second,
}

func makeKey(key string) string {
	return ENV_PREFIX + key
}

func getEnvOrDefaultDuration(key string) time.Duration {
	defaultValue, ok := defaultDurations[key]
	if !ok {
		panic(fmt.Sprintf("Default value for %v not available", key))
	}

	key = makeKey(key)
	durationStr := os.Getenv(key)
	if durationStr == "" {
		return defaultValue
	}

	duration, err := time.ParseDuration(durationStr)
	if err != nil {
		fmt.Printf("Could not get time duration using environment key: %v. Error: %v\n", key, err)
		os.Exit(1)
	}
	return duration
}

func mustHaveValue(key string) string {
	key = makeKey(key)
	value := os.Getenv(key)
	if value == "" {
		fmt.Printf("Could not get value using environment key: %v\n\n", key)

		PrintConfHelp()

		os.Exit(1)
	}
	return value
}

func loadConfig() {
	if cfg == nil {
		cfg = new(Config)
	}

	cfg.General.DB = mustHaveValue(keyDB)
	cfg.General.Timeout = getEnvOrDefaultDuration(keyTIMEOUT)
	cfg.General.UpdateTime = getEnvOrDefaultDuration(keyUPDATE_TIME)
	cfg.General.SampingTime = getEnvOrDefaultDuration(keySAMPING_TIME)
	cfg.General.TagUpdateTime = getEnvOrDefaultDuration(keyTAG_UPDATE_TIME)

	cfg.Auth.AccessToken = mustHaveValue(keyACCESS_TOKEN)
	cfg.Auth.AccessTokenSecret = mustHaveValue(keyACCESS_TOKEN_SECRET)
	cfg.Auth.ConsumerKey = mustHaveValue(keyCONSUMER_KEY)
	cfg.Auth.SecretKey = mustHaveValue(keySECRET_KEY)
	cfg.Auth.HMACSecret = mustHaveValue(keySECRET)
}

func loadRemoteConfigs() {
	rcfgs = make(RemoteConfigs, 0)

	remote_prefix := makeKey("REMOTE_")
	remote_regex := regexp.MustCompile(remote_prefix + `([A-Z\d_]+)_(URL|SECRET)$`)

	for _, e := range os.Environ() {
		pair := strings.Split(e, "=")

		remote := remote_regex.FindStringSubmatch(pair[0])
		if len(remote) < 2 {
			continue
		}

		name := strings.ToLower(remote[1])
		field := remote[2]

		rcfg, ok := rcfgs[name]
		if !ok {
			rcfg = RemoteConfig{}
		}

		if field == "URL" {
			rcfg.URL = pair[1]
		} else if field == "SECRET" {
			rcfg.HMACSecret = pair[1]
		}

		rcfgs[name] = rcfg
	}

	for key, rcfg := range rcfgs {
		if !rcfg.valid() {
			delete(rcfgs, key)
		}
	}
}

func PrintConfHelp() {
	keysOrder := []string{
		keyDB,
		keyTIMEOUT,
		keyUPDATE_TIME,
		keySAMPING_TIME,
		keyTAG_UPDATE_TIME,
		keyCONSUMER_KEY,
		keySECRET_KEY,
		keyACCESS_TOKEN,
		keyACCESS_TOKEN_SECRET,
		keySECRET,
	}

	fmt.Println("Environmental variables used in configuration")
	for _, key := range keysOrder {
		help := cfgHelp[key]

		envKey := makeKey(key)
		fmt.Println(envKey)
		value := os.Getenv(envKey)
		if value == "" {
			value = "not set"
			if duration, ok := defaultDurations[key]; ok {
				value = fmt.Sprintf("%v (default)", duration)
			}
		}
		fmt.Println("\tValue:", value)
		fmt.Println("\tHelp:", help)
		fmt.Println("")
	}
}
