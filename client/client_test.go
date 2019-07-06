package client

import (
	"encoding/json"
	"testing"
)

func TestClient_unmarshalErrorResponse(t *testing.T) {
	var errResponse = `{
  "data": {
    "createNote": null
  },
  "errors": [
    {
      "message": "タイトルを入力してください。",
      "locations": [
        {
          "line": 2,
          "column": 3
        }
      ],
      "path": [
        "createNote"
      ],
      "extensions": {
        "code": "BAD_REQUEST",
        "id": null,
        "details": {
          "title": [
            {
              "error": "blank"
            }
          ]
        }
      }
    }
  ]
}`

	var gResp response
	if err := json.Unmarshal([]byte(errResponse), &gResp); err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	if len(gResp.Errors) != 1 ||
		gResp.Errors[0].Message != "タイトルを入力してください。" ||
		gResp.Errors[0].Path[0] != "createNote" ||
		gResp.Errors[0].Extensions.Code != "BAD_REQUEST" {
		t.Errorf("gResp.Errors something went wrong: %#v", gResp.Errors)
	}
}
