package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/xerrors"
)

const (
	endpointBase  = "https://%s.kibe.la/api/v1"
	userAgentBase = "Songmu-kibela/%s (+https://github.com/Songmu/kibela)"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	token, endpoint string
	userAgent       string
	cli             Doer
}

func New(ver string) (*Client, error) {
	cli := &Client{token: os.Getenv("KIBELA_TOKEN")}
	if cli.token == "" {
		return nil, fmt.Errorf("set token by KIBELA_TOKEN env value")
	}
	team := os.Getenv("KIBELA_TEAM")
	if team == "" {
		return nil, fmt.Errorf("set team name by KIBELA_TEAM env value")
	}
	cli.endpoint = fmt.Sprintf(endpointBase, team)
	cli.cli = &http.Client{Transport: newRateLimitRoundTripper()}
	cli.userAgent = fmt.Sprintf(userAgentBase, ver)
	return cli, nil
}

func (cli *Client) Do(pa *Payload) (json.RawMessage, error) {
	body := bytes.Buffer{}
	if err := json.NewEncoder(&body).Encode(pa); err != nil {
		return nil, err
	}
	req, err := http.NewRequest(http.MethodPost, cli.endpoint, &body)
	if err != nil {
		return nil, err
	}
	req.Header.Set("Authorization", fmt.Sprintf("Bearer %s", cli.token))
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")
	req.Header.Set("User-Agent", cli.userAgent)

	resp, err := cli.cli.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	if resp.StatusCode == http.StatusTooManyRequests {
		return nil, ErrorTooManyRequet
	}
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, xerrors.Errorf("API response with code: %d, %s", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("API response with code: %d, response: %s", resp.StatusCode, string(bs))
	}
	var gResp response
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return nil, err
	}
	var resErr error
	if len(gResp.Errors) > 0 {
		resErr = gResp.Errors
	}
	return gResp.Data, resErr
}

type Payload struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables,omitempty"`
}

type response struct {
	Errors Errors          `json:"message,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}
