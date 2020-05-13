package utils

import (
	"bytes"
	"strings"
)

func splitVersion(v string) [3][]byte {
	var x [3][]byte
	s := strings.Split(v, ".")
	sLen := len(s)
	if sLen > 3 {
		s = s[:3]
	}
	for k, v := range s {
		x[k] = []byte(v)
	}
	return x
}

func VersionCompare(v1, v2 string) (result int) {
	x1 := splitVersion(v1)
	x2 := splitVersion(v2)

	for k := range x1 {
		result = bytes.Compare(x1[k], x2[k])
		if result != 0 {
			break
		}
	}
	return
}
