package kibelasync

import (
	"fmt"
	"reflect"
	"testing"
)

func TestKibela_getGroupCount(t *testing.T) {
	expect := 353
	ki := testKibela(newClient([]string{fmt.Sprintf(`{
  "data": {
    "groups": {
      "totalCount": %d
    }
  }
}`, expect)}))
	cnt, err := ki.getGroupCount()
	if err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	if cnt != expect {
		t.Errorf("out: %d, expect: %d", cnt, expect)
	}
}

func TestKibela_getGroups(t *testing.T) {
	ki := testKibela(newClient([]string{`{
  "data": {
    "groups": {
      "nodes": [
        {
          "id": "R3JvdXAvMQ",
          "name": "Home"
        },
        {
          "id": "R3JvdXAvMg",
          "name": "Test"
        }
      ]
    }
  }
}`}))
	out, err := ki.getGroups()
	if err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}

	expect := []*group{{
		ID:   ID("R3JvdXAvMQ"),
		Name: "Home",
	}, {
		ID:   ID("R3JvdXAvMg"),
		Name: "Test",
	}}

	if !reflect.DeepEqual(out, expect) {
		t.Errorf("\n   out: %+v\nexpect: %+v", out, expect)
	}
}
