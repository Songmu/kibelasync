package kibela

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

type gqError struct {
	Message    string            `json:"message"`
	Locations  []gqErrorLocation `json:"locations,omitempty"`
	Path       []interface{}     `json:"path,omitempty"` // string or uint
	Extensions gqErrorExtensions `json:"extensions"`
}

type gqErrorLocation struct {
	Line   uint `json:"line"`
	Column uint `json:"column"`
}

type gqErrorCode string

const (
	requestLimitExceeded gqErrorCode = "REQUEST_LIMIT_EXCEEDED"
	tokenBudgetExhausted gqErrorCode = "TOKEN_BUDGET_EXHAUSTED"
	teamBudgetExhausted  gqErrorCode = "TEAM_BUDGET_EXHAUSTED"
)

type gqErrorExtensions struct {
	Code              gqErrorCode `json:"code"`
	WaitMilliSecondes uint        `json:"waitMilliseconds,omitempty"`
}
