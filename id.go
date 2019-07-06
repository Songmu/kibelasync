package kibelasync

import (
	"encoding/base64"
	"fmt"
	"strconv"
	"strings"

	"golang.org/x/xerrors"
)

// ID represents kibela ID
type ID string

func newID(typ string, num int) ID {
	str := fmt.Sprintf("%s/%d", typ, num)
	return ID(base64.RawStdEncoding.EncodeToString([]byte(str)))
}

func (i ID) String() string {
	s, _ := base64.RawStdEncoding.DecodeString(string(i))
	return string(s)
}

func (i ID) Type() string {
	stuff := strings.Split(i.String(), "/")
	return stuff[0]
}

func (i ID) Number() (int, error) {
	stuff := strings.Split(i.String(), "/")
	if len(stuff) != 2 {
		return 0, fmt.Errorf("invalid id: %s", string(i))
	}
	num, err := strconv.Atoi(stuff[1])
	if err != nil {
		return 0, xerrors.Errorf("invalid id: %s, error: %w", i.String(), err)
	}
	return num, nil
}
