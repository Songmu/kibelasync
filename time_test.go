package kibela

import (
	"bytes"
	"encoding/json"
	"fmt"
	"testing"
	"time"
)

func mustTime(tstr string) Time {
	tt, err := time.Parse(rfc3339Milli, tstr)
	if err != nil {
		panic(err)
	}
	return Time{Time: tt}
}

type testTime struct {
	Time Time `json:"time"`
}

func TestTime_MarshalJSON(t *testing.T) {
	in := "2019-06-23T16:54:09.447+09:00"
	tt := testTime{Time: mustTime(in)}
	buf := &bytes.Buffer{}
	if err := json.NewEncoder(buf).Encode(tt); err != nil {
		t.Errorf("error should be nil but: %s", err)
	}
	expect := fmt.Sprintf(`{"time":"%s"}`+"\n", in)
	out := buf.String()
	if expect != out {
		t.Errorf("\nexpect: %s\n   out: %s", expect, out)
	}
}
