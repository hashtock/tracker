package client

import (
    "bytes"
    "crypto"
    "crypto/hmac"
    "crypto/md5"
    _ "crypto/sha1"
    "encoding/base64"
    "io"
    "io/ioutil"
    "log"
    "net/http"
    "net/url"
    "time"
)

type Tracker struct {
    HMACSecret string
    Host       string
    Client     *http.Client
}

func NewTracker(secret string, host string) Tracker {
    return Tracker{
        HMACSecret: secret,
        Host:       host,
        Client:     http.DefaultClient,
    }
}

func (t *Tracker) GetTagList() (tags []string, err error) {
    res, lerr := t.doSignedRequest("GET", "/api/tag/")
    if lerr != nil {
        err = lerr
        return
    }

    log.Println("Status:", res.Status)
    body, err := ioutil.ReadAll(res.Body)
    res.Body.Close()
    log.Println("Body:", string(body))
    return
}

func (t Tracker) doSignedRequest(method string, path string) (*http.Response, error) {
    url := url.URL{
        Scheme: "http",
        Host:   t.Host,
        Path:   path,
    }

    req, err := http.NewRequest(method, url.String(), nil)
    if err != nil {
        log.Fatalln(err)
        return nil, err
    }

    h := md5.New()
    if req.Body != nil {
        io.Copy(h, req.Body)
    }

    req.Header.Add("Content-Type", "application/json")
    req.Header.Add("Content-MD5", string(h.Sum(nil)))
    req.Header.Add("Date", time.Now().Format(time.ANSIC))
    sig := "HashTock tracker:" + t.generateSignature(req)
    req.Header.Add("Authorization", sig)

    log.Println("Sig:", sig)

    return t.Client.Do(req)
}

func (t *Tracker) generateSignature(req *http.Request) string {
    var newline = "\u000A"
    var buffer bytes.Buffer

    buffer.WriteString(req.Method)
    buffer.WriteString(newline)

    buffer.WriteString(req.Header.Get("Content-MD5"))
    buffer.WriteString(newline)

    buffer.WriteString(req.Header.Get("Content-Type"))
    buffer.WriteString(newline)

    buffer.WriteString(req.Header.Get("Date"))
    buffer.WriteString(newline)

    buffer.WriteString(req.URL.Path)

    signingString := buffer.String()

    mac := hmac.New(crypto.SHA1.New, []byte(t.HMACSecret))
    mac.Write([]byte(signingString))

    signature := base64.StdEncoding.EncodeToString(mac.Sum(nil))

    return string(signature)
}
