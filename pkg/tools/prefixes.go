package tools

import (
	"fmt"
	"strings"
)

func AddPrefix(str string, pre string) (string, error) {
	if len(strings.TrimSpace(str)) == 0 {
		return "", fmt.Errorf("%w: %s", ErrEmptyKey, str)
	}
	if len(strings.TrimSpace(pre)) == 0 {
		return "", fmt.Errorf("%w: %s", ErrEmptyKey, pre)
	}
	if strings.HasPrefix(str, pre) {
		return str, nil
	}
	return pre + str, nil
}
