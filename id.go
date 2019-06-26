package kibela

import "encoding/base64"

type ID string

func (i ID) String() string {
	s, _ := base64.RawStdEncoding.DecodeString(string(i))
	return string(s)
}
