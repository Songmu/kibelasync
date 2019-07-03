package kibela

import (
	"fmt"
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
