package kibela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"os"

	"golang.org/x/xerrors"
)

const endpointBase = "https://%s.kibe.la/api/v1"

var defaultUserAgent string

func init() {
	defaultUserAgent = "Songmu-kibela/" + version + " (+https://github.com/Songmu/kibela)"
}

type doer interface {
	Do(*http.Request) (*http.Response, error)
}

type client struct {
	token, endpoint string
	userAgent       string
	cli             doer
}

func newClient() (*client, error) {
	cli := &client{token: os.Getenv("KIBELA_TOKEN")}
	if cli.token == "" {
		return nil, fmt.Errorf("set token by KIBELA_TOKEN env value")
	}
	team := os.Getenv("KIBELA_TEAM")
	if team == "" {
		return nil, fmt.Errorf("set team name by KIBELA_TEAM env value")
	}
	cli.endpoint = fmt.Sprintf(endpointBase, team)
	cli.cli = &http.Client{Transport: newRateLimitRoundTripper()}
	cli.userAgent = defaultUserAgent
	return cli, nil
}

func (cli *client) Do(pa *payload) (*gqResponse, error) {
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
	if resp.StatusCode < 200 || resp.StatusCode > 299 {
		bs, err := ioutil.ReadAll(resp.Body)
		if err != nil {
			return nil, fmt.Errorf("API response with code: %d, %s", resp.StatusCode, err)
		}
		return nil, fmt.Errorf("API response with code: %d, response: %s", resp.StatusCode, string(bs))
	}
	var gResp gqResponse
	if err := json.NewDecoder(resp.Body).Decode(&gResp); err != nil {
		return nil, err
	}
	if gResp.Data == nil {
		return nil, gResp.Errors
	}
	return &gResp, nil
}

type payload struct {
	Query     string      `json:"query"`
	Variables interface{} `json:"variables,omitempty"`
}

type gqResponse struct {
	Errors gqErrors        `json:"message,omitempty"`
	Data   json.RawMessage `json:"data,omitempty"`
}

/*
{
  "data": {
    "notes": {
      "totalCount": 353
    }
  }
}
*/
// OK
func (cli *client) getNotesCount() (int, error) {
	gResp, err := cli.Do(&payload{Query: totalCountQuery})
	if err != nil {
		return 0, xerrors.Errorf("failed to cli.getNotesCount: %w", err)
	}
	var res struct {
		Notes struct {
			TotalCount int `json:"totalCount"`
		} `json:"notes"`
	}
	if err := json.Unmarshal(gResp.Data, &res); err != nil {
		return 0, xerrors.Errorf("failed to cli.getNotesCount: %w", err)
	}
	return res.Notes.TotalCount, nil
}
