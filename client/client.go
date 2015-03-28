package client

import (
	"bytes"
	"crypto"
	"crypto/hmac"
	"crypto/md5"
	_ "crypto/sha1" // Register hashing alg
	"encoding/base64"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"time"

	"github.com/hashtock/tracker/conf"
	"github.com/hashtock/tracker/core"
)

type Tracker struct {
	HMACSecret string
	Host       string
	Client     *http.Client
}

func NewTracker(conf conf.RemoteConfig) (*Tracker, error) {
	return &Tracker{
		HMACSecret: conf.HMACSecret,
		Host:       conf.URL,
		Client:     http.DefaultClient,
	}, nil
}

func NewTrackerPlain(HMACSecret, URL string) (*Tracker, error) {
	return NewTracker(conf.RemoteConfig{
		HMACSecret: HMACSecret,
		URL:        URL,
	})
}

func (t *Tracker) Tags() (tags []core.Tag, err error) {
	res, lerr := t.doSignedRequest("GET", "/api/tag/", nil)
	if lerr != nil {
		err = lerr
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()

	if err := json.Unmarshal(body, &tags); err != nil {
		log.Fatalln(err)
	}

	return
}

func (t *Tracker) Counts(since, until time.Time) (counts []core.TagCount, err error) {
	query := url.Values{}
	query.Set("since", since.Format(time.RFC3339))
	query.Set("until", until.Format(time.RFC3339))

	res, lerr := t.doSignedRequest("GET", "/api/counts/", query)
	if lerr != nil {
		err = lerr
		return
	}

	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		log.Fatalln(err)
	}
	res.Body.Close()

	if err := json.Unmarshal(body, &counts); err != nil {
		log.Fatalln(err)
	}
	return
}

func (t *Tracker) AddTag(tag string) (err error) {
	uri := fmt.Sprintf("/api/tag/%s/", tag)
	_, err = t.doSignedRequest("PUT", uri, nil)
	return
}

func (t *Tracker) TagTrends(tag string, since, until time.Time, sampling core.Sampling) (tagCounts core.TagCountTrend, err error) {
	rc, err := t.trendRequest(tag, since, until, sampling)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(rc)
	err = decoder.Decode(&tagCounts)
	rc.Close()

	return
}

func (t *Tracker) Trends(since, until time.Time) (trends []core.TagCountTrend, err error) {
	rc, err := t.trendRequest("", since, until, core.SamplingRaw)
	if err != nil {
		return
	}

	decoder := json.NewDecoder(rc)
	err = decoder.Decode(&trends)
	rc.Close()

	return
}

func (t *Tracker) trendRequest(tag string, since, until time.Time, sampling core.Sampling) (rc io.ReadCloser, err error) {
	query := url.Values{}
	query.Set("since", since.Format(time.RFC3339))
	query.Set("until", until.Format(time.RFC3339))

	if sampling != core.SamplingRaw {
		query.Set("sampling", sampling.String())
	}

	if tag == "" && sampling.String() != "" {
		err = errors.New("Sampling can be specified only for single tag")
		return
	}

	url := "/api/trends/"
	if tag != "" {
		url = fmt.Sprintf("%v%v/", url, tag)
	}

	rec, err := t.doSignedRequest("GET", url, query)
	if err != nil {
		return
	}
	rc = rec.Body
	return
}

func (t Tracker) doSignedRequest(method, path string, query url.Values) (*http.Response, error) {
	url := url.URL{
		Scheme:   "http",
		Host:     t.Host,
		Path:     path,
		RawQuery: query.Encode(),
	}

	req, err := http.NewRequest(method, url.String(), nil)
	if err != nil {
		return nil, err
	}

	h := md5.New()
	if req.Body != nil {
		io.Copy(h, req.Body)
	}

	req.Header.Add("Content-Type", "application/json")
	req.Header.Add("Content-MD5", fmt.Sprintf("%x", h.Sum(nil)))
	req.Header.Add("Date", time.Now().Format(time.ANSIC))
	sig := "HashTock tracker:" + t.generateSignature(req)
	req.Header.Add("Authorization", sig)

	res, err := t.Client.Do(req)
	if err != nil {
		return nil, err
	}

	if res.StatusCode >= 400 {
		return nil, fmt.Errorf("Problem with request: %s", res.Status)
	}

	return res, nil
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
