package tools

import (
	"strconv"
	"strings"
)

func AddPrefix(str, pre string) string {
	if strings.HasPrefix(str, pre) {
		return str
	}
	return pre + str
}

func GetIntID(str, prefix string) (int, error) {
	return strconv.Atoi(strings.TrimPrefix(str, prefix))
}
