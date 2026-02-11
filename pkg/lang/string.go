package lang

import (
	"fmt"
	"strings"
	"unicode/utf8"
)

const (
	Empty       = ""
	Space       = ' '
	SpaceString = " "
	Slash       = "/"
)

func IsLengthLongerThen(p *string, l int) bool {
	return p != nil && len(*p) > l
}

func IsLengthLessThen(p *string, l int) bool {
	return p != nil && len(*p) < l
}

func IsRuneCountLongerThen(p *string, l int) bool {
	return p != nil && utf8.RuneCountInString(*p) > l
}

func IsRuneCountLessThen(p *string, l int) bool {
	return p != nil && utf8.RuneCountInString(*p) < l
}

func FormatL(length int, p string, args ...any) string {
	if p == "" {
		return ""
	}
	s := fmt.Sprintf(p, args...)
	if 0 < length && length < len(s) {
		return s[:length]
	}
	return s
}

func Substring(p *string, start int, end int) string {
	if p == nil || start < 0 || end < 0 || start >= end || start >= len(*p) {
		return ""
	}
	if end > len(*p) {
		end = len(*p)
	}
	return (*p)[start:end]
}

// IsBlank checks if a string is blank or not.
func IsBlank(str string) bool {
	return strings.TrimSpace(str) == ""
}

// JoinWith concatenates strings from the parts slice, separated by sep.
// It skips nil pointers in the parts slice.
func JoinWith(sep string, parts ...*string) string {
	var builder strings.Builder
	for _, part := range parts {
		if part != nil {
			if builder.Len() > 0 {
				builder.WriteString(sep)
			}
			builder.WriteString(*part)
		}
	}
	return builder.String()
}

// TrimJoin will concatenates strings in rune compitable mode after trimming
// space and separator, both left and right, from each part.
func TrimJoin(s string, parts ...*string) string {
	var builder strings.Builder
	sep := Ternary(len(s) > 0, []rune(s), []rune{})
	seplen := len(sep)
	var ss, se *rune
	if seplen > 0 {
		ss, se = Dup(sep[0]), Dup(sep[seplen-1])
	}
	for _, part := range parts {
		if part == nil || len(*part) == 0 || (*part) == SpaceString {
			continue
		}
		runes := []rune(*part)
		beg, end, last := 0, len(runes), len(runes)-1
		// prefix
		for c := runes[beg]; beg < end; c = runes[beg] {
			if nbeg := beg + seplen; ss != nil && c == *ss && nbeg < end && string(runes[beg:nbeg]) == s {
				// sep match prior
				beg += seplen
			} else if c == Space {
				beg += 1
			} else {
				break
			}
		}
		// suffix
		for c := runes[last]; last >= 0; c = runes[last] {
			// 0 <= last+1-seplen => -1 <= last-seplen
			if nbeg, nend := last+1-seplen, last+1; se != nil && c == *se && 0 <= nbeg && string(runes[nbeg:nend]) == s {
				// sep match prior
				last -= seplen
			} else if c == Space {
				last -= 1
			} else {
				break
			}
		}
		if beg <= last {
			if builder.Len() > 0 {
				builder.WriteString(s)
			}
			builder.WriteString(string(runes[beg : last+1]))
		}
	}
	return builder.String()
}

// Join concatenates string pointers with a slash separator.
func Join(parts ...*string) string {
	return JoinWith(Slash, parts...)
}
