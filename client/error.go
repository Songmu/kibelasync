package client

import (
	"errors"
	"strings"
)

var ErrorTooManyRequet = errors.New("too many request")

/*
{
  "errors": [
    {
      "message": "Unexpected end of document",
      "extensions": {
        "code": "PARSE_ERROR"
      }
    }
  ]
}
*/
type Errors []Error

func (e Errors) Error() string {
	errs := make([]string, len(e))
	for i, e := range e {
		errs[i] = e.Message
	}
	if len(errs) == 0 {
		return "unknown error"
	}
	return strings.Join(errs, "\n")
}

type Error struct {
	Message    string          `json:"message"`
	Locations  []ErrorLocation `json:"locations,omitempty"`
	Path       []interface{}   `json:"path,omitempty"` // string or uint
	Extensions ErrorExtensions `json:"extensions"`
}

type ErrorLocation struct {
	Line   uint `json:"line"`
	Column uint `json:"column"`
}

type ErrorCode string

const (
	RequestLimitExceeded ErrorCode = "REQUEST_LIMIT_EXCEEDED"
	TokenBudgetExhausted ErrorCode = "TOKEN_BUDGET_EXHAUSTED"
	TeamBudgetExhausted  ErrorCode = "TEAM_BUDGET_EXHAUSTED"
)

type ErrorExtensions struct {
	Code              ErrorCode `json:"code"`
	WaitMilliSecondes uint      `json:"waitMilliseconds,omitempty"`
}
