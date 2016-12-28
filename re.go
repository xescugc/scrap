package main

import "regexp"

var (
	searchRe = map[string]*regexp.Regexp{
		"twitter": regexp.MustCompile(`https://twitter.com/(\w+)`),
		//"email":   regexp.MustCompile(`([^@\s][!#$%&'*+-\/=?^{|}~\w]+@[^@\s][!#$%&'*+-\/=?^{|}~\w]+)`),
		//"email": regexp.MustCompile(`([^@\s][!#$%&'*+-\/=?^{|}~\w]+@[^@\s](?:(?:[^.][A-Za-z\d\-]+)|(?:\[(?:\d{1,3}?\.){3}\d{1,3}\])))`),
		//"email": regexp.MustCompile(`([\w+\-]+@[a-z\d\-]+(\.[a-z]+)*\.[a-z]+)`),
		"email": regexp.MustCompile(`((?:[a-z0-9!#$%&'*+/=?^_{|}~-]+(?:\.[a-z0-9!#$%&'*+/=?^_{|}~-]+)*|"(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21\x23-\x5b\x5d-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])*")@(?:(?:[a-z0-9](?:[a-z0-9-]*[a-z0-9])?\.)+[a-z0-9](?:[a-z0-9-]*[a-z0-9])?|\[(?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?|[a-z0-9-]*[a-z0-9]:(?:[\x01-\x08\x0b\x0c\x0e-\x1f\x21-\x5a\x53-\x7f]|\\[\x01-\x09\x0b\x0c\x0e-\x7f])+)\]))`),
	}
)

func applyRe(body string) ([]string, error) {
	matches := make(map[string]int)
	for _, match := range searchRe[*search].FindAllStringSubmatch(body, -1) {
		if len(match) > 1 {
			matches[match[1]] = 0
		}
	}
	result := make([]string, 0)
	for k, _ := range matches {
		result = append(result, k)
	}
	return result, nil
}
