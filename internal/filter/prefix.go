package filter

import (
	"bytes"
)

func Prefix(prefix []byte) FilterFunc {
	return func(line []byte) (bool, error) {
		return bytes.HasPrefix(line, prefix), nil
	}
}
