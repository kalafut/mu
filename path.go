package mu

import (
	"path/filepath"
	"strings"
)

func SplitExt(s string) (string, string) {
	ext := filepath.Ext(s)
	head := strings.TrimSuffix(s, ext)

	return head, ext
}
