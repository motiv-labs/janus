package godog

import (
	"strings"

	"github.com/cucumber/messages-go/v10"
)

// based on http://behat.readthedocs.org/en/v2.5/guides/6.cli.html#gherkin-filters
func matchesTags(filter string, tags []*messages.Pickle_PickleTag) (ok bool) {
	ok = true

	for _, andTags := range strings.Split(filter, "&&") {
		var okComma bool

		for _, tag := range strings.Split(andTags, ",") {
			tag = strings.Replace(strings.TrimSpace(tag), "@", "", -1)

			okComma = hasTag(tags, tag) || okComma
			if tag[0] == '~' {
				tag = tag[1:]
				okComma = !hasTag(tags, tag) || okComma
			}
		}

		ok = ok && okComma
	}

	return
}

func hasTag(tags []*messages.Pickle_PickleTag, tag string) bool {
	for _, t := range tags {
		tName := strings.Replace(t.Name, "@", "", -1)

		if tName == tag {
			return true
		}
	}

	return false
}
