package client

import (
	"errors"
	"strings"
)

// ErrorTooManyRequet is an error representing too many request
var ErrorTooManyRequet = errors.New("too many request")

// Errors is an error type contained multiple errors
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

// Error is for error type
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

// Error is error type for GraphQL request
type Error struct {
	Message    string          `json:"message"`
	Locations  []ErrorLocation `json:"locations,omitempty"`
	Path       []interface{}   `json:"path,omitempty"` // string or uint
	Extensions ErrorExtensions `json:"extensions"`
}

// ErrorLocation is represents error location
type ErrorLocation struct {
	Line   uint `json:"line"`
	Column uint `json:"column"`
}

// ErrorCode is for representing error types
type ErrorCode string

// ErrorCodes
const (
	RequestLimitExceeded ErrorCode = "REQUEST_LIMIT_EXCEEDED"
	TokenBudgetExhausted ErrorCode = "TOKEN_BUDGET_EXHAUSTED"
	TeamBudgetExhausted  ErrorCode = "TEAM_BUDGET_EXHAUSTED"
)

// ErrorExtensions is for extra information of errors
type ErrorExtensions struct {
	Code              ErrorCode `json:"code"`
	WaitMilliSecondes uint      `json:"waitMilliseconds,omitempty"`
}
