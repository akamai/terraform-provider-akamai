package str

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

// GetIntID is used to get the id out from the string.
func GetIntID(str, prefix string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(str, prefix))
}
