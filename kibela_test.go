package kibela

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"reflect"
	"strings"
	"testing"

	"github.com/Songmu/kibela/client"
)

type testDoer struct {
	cursor        int
	responseTexts []string
}

func (td *testDoer) Do(req *http.Request) (*http.Response, error) {
	bodyText := td.responseTexts[td.cursor%len(td.responseTexts)]
	td.cursor++
	return &http.Response{
		Status:     "200 OK",
		StatusCode: http.StatusOK,
		Proto:      "HTTP/1.0",
		ProtoMajor: 1,
		Header:     make(http.Header),
		Close:      true,
		Body:       ioutil.NopCloser(bytes.NewReader([]byte(bodyText))),
		Request:    req,
	}, nil
}

var _ client.Doer = (*testDoer)(nil)

func newClient(responseTexts []string) *client.Client {
	return client.Test(&testDoer{responseTexts: responseTexts})
}

func testKibela(cli *client.Client) *kibela {
	return &kibela{cli: cli}
}

func TestKibela_setGroups(t *testing.T) {
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
	if err := ki.setGroups(); err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	expect := map[string]ID{
		"Home": ID("R3JvdXAvMQ"),
		"Test": ID("R3JvdXAvMg"),
	}
	if !reflect.DeepEqual(ki.groups, expect) {
		t.Errorf("got: %v, expect: %v", ki.groups, expect)
	}
}

func TestKibela_fetchGroupID(t *testing.T) {
	t.Run("fetch error", func(t *testing.T) {
		ki := testKibela(newClient([]string{`{
  "errors": [{
    "message": "error!"
  }]
}`}))
		_, err := ki.fetchGroupID("Home")
		if err == nil || !strings.Contains(err.Error(), "while setGroupID") {
			t.Errorf("error should be occurred and contains `while setGroupID`, but: %v", err)
		}
	})

	t.Run("fetch success", func(t *testing.T) {
		expect := "R3JvdXAvMQ"
		ki := testKibela(newClient([]string{fmt.Sprintf(`{
  "data": {
    "groups": {
      "nodes": [
        {
          "id": "%s",
          "name": "Home"
        },
        {
          "id": "R3JvdXAvMg",
          "name": "Test"
        }
      ]
    }
  }
}`, expect)}))
		id, err := ki.fetchGroupID("Home")
		if err != nil {
			t.Errorf("error should be nil, but: %s", err)
		}
		if id != ID(expect) {
			t.Errorf("got: %s, expect: %s", string(id), expect)
		}
		_, err = ki.fetchGroupID("Unknown")
		if err == nil || !strings.Contains(err.Error(), `group "Unknown"`) {
			t.Errorf("error should be occurred and contains `group \"Unknown\"`, but: %v", err)
		}
	})
}
