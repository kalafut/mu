package mu

import "strings"

func StrPtr(s string) *string {
	return &s
}

func StrPtrOrNil(s string) *string {
	s = strings.TrimSpace(s)
	if s == "" {
		return nil
	}
	return &s
}

func IntPtr(i int) *int {
	return &i
}

func StrFromPtr(s *string) string {
	if s != nil {
		return *s
	}

	return ""
}

func IntFromPtr(s *int) (result int) {
	if s != nil {
		return *s
	}

	return 0
}
