package stringtoint

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"
)

// Parse returns the amount of bytes for the human notation.
// B for bytes
// K for kilobytes
// M for megabytes
// G for gigabytes
// T for Terabytes
func Parse(s string) (int64, error) {
	reInt := regexp.MustCompile("^([0-9]+)$")
	reHuman := regexp.MustCompile("(?i)^([0-9]+)(B|K|M|G|T)$")
	if reInt.MatchString(s) {
		s = s + "b"
	} else if reHuman.MatchString(s) {
		s = strings.ToLower(s)
	} else {
		return -1, fmt.Errorf("unknown format")
	}

	switch s[len(s)-1:] {
	case "b":
		return calculateBytes(s, 1)
	case "k":
		return calculateBytes(s, 1024)
	case "m":
		return calculateBytes(s, 1024*1024)
	case "g":
		return calculateBytes(s, 1024*1024*1024)
	case "t":
		return calculateBytes(s, 1024*1024*1024*1024)
	default:
		return -1, fmt.Errorf("unknown format")
	}
}

func calculateBytes(s string, multiplier int64) (int64, error) {
	s = s[0 : len(s)-1]
	number, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return -1, fmt.Errorf("could not convert string to int got %s with error %s", s, err)
	}

	return number * int64(multiplier), nil
}
