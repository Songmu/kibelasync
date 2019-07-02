package kibela

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
)

func TestNoteUnmarshalJSON(t *testing.T) {
	input := `{
     "id": "QmxvZy8zNjY",
     "title": "APIテストpublic",
     "content": "コンテント!\nコンテント",
     "coediting": true,
     "folderName": "testtop/testsub1",
     "groups": [
       {
         "name": "Home",
         "id": "R3JvdXAvMQ"
       }
     ],
     "author": {
       "account": "Songmu"
     }
   }`
	var n note
	if err := json.NewDecoder(strings.NewReader(input)).Decode(&n); err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	expect := note{
		ID:        ID("QmxvZy8zNjY"),
		Title:     "APIテストpublic",
		Content:   "コンテント!\nコンテント",
		CoEditing: true,
		Folder:    "testtop/testsub1",
		Groups: []*group{
			{
				ID:   ID("R3JvdXAvMQ"),
				Name: "Home",
			},
		},
		Author: struct {
			Account string `json:"account"`
		}{
			Account: "Songmu",
		},
	}

	if !reflect.DeepEqual(n, expect) {
		t.Errorf("got:\n%#v expect:\n%#v", n, expect)
	}
}
