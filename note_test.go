package kibela

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
     },
     "createdAt": "2019-06-23T16:54:09.447+09:00",
     "publishedAt": "2019-06-23T16:54:09.444+09:00",
     "contentUpdatedAt": "2019-06-23T16:54:09.445+09:00",
     "updatedAt": "2019-06-23T17:22:38.496+09:00"
   }`
	var n note
	if err := json.NewDecoder(strings.NewReader(input)).Decode(&n); err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	expect := note{}

	if !reflect.DeepEqual(n, expect) {
		t.Errorf("got:\n%#v expect:\n%#v", n, expect)
	}

	fmt.Println(n.ContentUpdatedAt.String())
}
