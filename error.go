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

type gqErrors struct {
	Errors []gqError   `json:"message"`
	Data   interface{} `json:"data,omitempty"`
}

type gqError struct {
	Message    string            `json:"message"`
	Extensions gqErrorExtensions `json:"extensions"`
}

type gqErrorExtensions struct {
	Code string `json:"code"`
}
