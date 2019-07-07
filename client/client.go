package client

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"strings"

	"golang.org/x/xerrors"
)

const (
	endpointBase  = "https://%s.kibe.la/api/v1"
	userAgentBase = "Songmu-kibelasync/%s (+https://github.com/Songmu/kibelasync)"
)

type Doer interface {
	Do(*http.Request) (*http.Response, error)
}

type Client struct {
	token, endpoint string
	userAgent       string
	cli             Doer
	limiter         *rateLimitRoundTripper
}

type budget struct {
	Cost      int `json:"cost,string"`
	Consumed  int `json:"consumed,string"`
	Remaining int `json:"remaining,string"`
}

func New(ver, team, token string) (*Client, error) {
	cli := &Client{token: token}
	cli.endpoint = fmt.Sprintf(endpointBase, team)
	cli.limiter = newRateLimitRoundTripper()
	cli.cli = &http.Client{Transport: cli.limiter}
	cli.userAgent = fmt.Sprintf(userAgentBase, ver)
	return cli, nil
}

func (cli *Client) Do(pa *Payload) (json.RawMessage, error) {
	pa.Query = strings.TrimSpace(pa.Query)
	isQuery := !strings.HasPrefix(pa.Query, "mutation")
	if isQuery {
		// inject cost query (these fields are seems to be no cost now)
		q := pa.Query
		q = strings.TrimSuffix(q, "}")
		q += `
  budget {
    cost
    consumed
    remaining
  }
}`
		pa.Query = q
	}

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
	if isQuery && cli.limiter != nil {
		var res struct {
			Budget *budget `json:"budget"`
		}
		if err := json.Unmarshal(gResp.Data, &res); err != nil {
			log.Printf("failed to retrieve budgets from response: %s\n", err)
		}
		if res.Budget != nil {
			cli.limiter.announceRemainingCost(res.Budget.Remaining)
		}
	}
	return gResp.Data, resErr
}

type Payload struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables,omitempty"`
}

type response struct {
	Errors Errors          `json:"errors,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

// Test for create test client for using testing only
func Test(cli Doer) *Client {
	return &Client{cli: cli}
}
