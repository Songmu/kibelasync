package kibelasync

import (
	"encoding/json"
	"fmt"
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
	// omit {"updatedAt": "2019-06-23T17:22:38.496Z"} for testing
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

func TestKibela_getNotesCount(t *testing.T) {
	expect := 353
	ki := testKibela(newClient([]string{fmt.Sprintf(`{
  "data": {
    "notes": {
      "totalCount": %d
    }
  }
}`, expect)}))
	cnt, err := ki.getNotesCount()
	if err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	if cnt != expect {
		t.Errorf("out: %d, expect: %d", cnt, expect)
	}
}
