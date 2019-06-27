package kibela

import (
	"encoding/json"
	"reflect"
	"strings"
	"testing"
	"time"
)

func mustTime(tstr string) Time {
	tt, err := time.Parse(rfc3339Milli, "2019-06-23T17:22:38.496Z")
	if err != nil {
		panic(err)
	}
	return Time{Time: tt}
}

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
		Groups: []struct {
			ID   `json:"id"`
			Name string `json:"name"`
		}{
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
		// CreatedAt:        mustTime("2019-06-23T16:54:09.447+09:00"),
		// PublishedAt:      mustTime("2019-06-23T16:54:09.444+09:00"),
		// ContentUpdatedAt: mustTime("2019-06-23T16:54:09.445+09:00"),
		// UpdatedAt:        mustTime("2019-06-23T17:22:38.496+09:00"),
	}

	if !reflect.DeepEqual(n, expect) {
		t.Errorf("got:\n%#v expect:\n%#v", n, expect)
	}
}
