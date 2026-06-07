package domain

import "strings"

func DisplayName(name, nickname string) string {
	nickname = strings.TrimSpace(nickname)
	if nickname != "" {
		return nickname
	}
	return name
}
