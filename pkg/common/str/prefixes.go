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

// GetInt64ID is used to get the id out from the string as int64.
func GetInt64ID(str, prefix string) (int64, error) {
	return strconv.ParseInt(strings.TrimPrefix(str, prefix), 10, 64)
}
