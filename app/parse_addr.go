package app

import "regexp"

var parseAddrPattern = regexp.MustCompile(`(\w+)://([\w.]+):(\d+)`)

func parseAddr(addr string) (proto string, host string, port string, ok bool) {
	match := parseAddrPattern.FindStringSubmatch(addr)
	if len(match) != 4 {
		return "", "", "", false
	}

	return match[1], match[2], match[3], true
}
