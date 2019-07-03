package kibela

import "testing"

func TestID(t *testing.T) {
	id := newID("Blog", 366)

	expectString := "Blog/366"
	if id.String() != expectString {
		t.Errorf("id.String() = %q, expect: %q", id, expectString)
	}

	expectID := "QmxvZy8zNjY"
	if string(id) != expectID {
		t.Errorf("string(id) = %q, expect: %q", string(id), expectID)
	}

	if id.Type() != "Blog" {
		t.Errorf("id.String() = %q, expect: %q", id.Type(), "Blog")
	}

	num, err := id.Number()
	if err != nil {
		t.Errorf("error should be nil, but: %s", err)
	}
	if num != 366 {
		t.Errorf("id.Number() = %d, expect: %d", num, 366)
	}
}
