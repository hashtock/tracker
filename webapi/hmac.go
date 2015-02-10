package webapi

import (
    "crypto"
    _ "crypto/sha1"
    "log"
    "time"

    "github.com/auroratechnologies/vangoh"

    "github.com/hashtock/tracker/conf"
)

type constSecretProvider struct {
    secret []byte
}

func (c *constSecretProvider) GetSecretKey(id []byte) ([]byte, error) {
    return c.secret, nil
}

func newVanGoh() *vangoh.VanGoH {
    cfg := conf.GetConfig()

    if cfg.Auth.HMACSecret == "" {
        log.Fatalln("HMACSecret key not present.")
    }

    secretProvider := &constSecretProvider{
        secret: []byte(cfg.Auth.HMACSecret),
    }

    vg := vangoh.NewSingleProvider(secretProvider)
    vg.SetAlgorithm(crypto.SHA1.New)
    vg.SetMaxTimeSkew(time.Minute * 15)

    return vg
}
