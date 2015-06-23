package client

import (
	"encoding/json"
	"fmt"
	"net/http"

	"github.com/hashtock/auth/core"
)

type Client struct {
	serviceLocation string
	HttpClient      *http.Client
}

func NewClient(authServiceLocation string) *Client {
	return &Client{
		serviceLocation: authServiceLocation,
		HttpClient:      http.DefaultClient,
	}
}

func (c Client) Who(req *http.Request) (*core.User, error) {
	whoReq, err := http.NewRequest("GET", c.serviceLocation, nil)
	if err != nil {
		return nil, err
	}

	for _, cookie := range req.Cookies() {
		whoReq.AddCookie(cookie)
	}
	resp, err := c.HttpClient.Do(whoReq)

	if err != nil {
		return nil, err
	}

	if resp.StatusCode == http.StatusUnauthorized {
		return nil, core.ErrUserNotLoggedIn
	} else if resp.StatusCode >= 400 {
		return nil, fmt.Errorf("Problem with request: %s", resp.Status)
	}

	user := new(core.User)

	decoder := json.NewDecoder(resp.Body)
	err = decoder.Decode(*user)
	resp.Body.Close()

	return user, nil
}
