package fits

import (
	"strconv"
)

type Results struct {
	Types   []Type   `json:"type"`
	Methods []Method `json:"method"`
}

func toStr(f float64) string {
	return strconv.FormatFloat(f, 'g', -1, 64)
}
