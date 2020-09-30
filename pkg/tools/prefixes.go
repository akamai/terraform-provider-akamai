package tools

import (
	"errors"
	"fmt"
	"strings"
)

func AddPrefix(str string, pre string) (string, error) {
	if len(strings.TrimSpace(str)) == 0 {
		return "", fmt.Errorf("%w: %s", errors.New("Prefix string cannot be blank"), str)
	}
	if len(strings.TrimSpace(pre)) == 0 {
		return "", fmt.Errorf("%w: %s", errors.New("Prefix cannot be blank"), pre)
	}
	if strings.HasPrefix(str, pre) {
		return str, nil
	}
	return pre + str, nil
}
