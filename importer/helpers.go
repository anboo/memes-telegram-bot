package importer

import (
	"strings"
)

var blackList = []string{
	"http://",
	"https://",
}

func BlackListed(text string) bool {
	if len(text) >= 1024 {
		return false
	}

	for _, bad := range blackList {
		if strings.Contains(text, bad) {
			return true
		}
	}
	return false
}
