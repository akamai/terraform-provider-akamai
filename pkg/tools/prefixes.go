package tools

import (
	"strconv"
	"strings"
)

// AddPrefix will add prefix to given string.
func AddPrefix(str, pre string) string {
	if str == "" {
		return ""
	}
	if strings.HasPrefix(str, pre) {
		return str
	}
	return pre + str
}

// TrimPrefix will extract prefix from the given string if it is exist.
func TrimPrefix(str, prefix string) string {
	if strings.HasPrefix(str, prefix) {
		return strings.TrimPrefix(str, prefix)
	}
	return str
}

// GetIntID is used to get the id out from the string.
func GetIntID(str, prefix string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(str, prefix))
}
