// re for extend regexp

package re

import (
	"regexp"
)

//RegexpGroup: regexp group to map like python re.groups()
func RegexpGroup(exp *regexp.Regexp, text string) map[string]string {
	match := exp.FindStringSubmatch(text)
	result := make(map[string]string)
	for i, name := range exp.SubexpNames() {
		if i != 0 && name != "" {
			result[name] = match[i]
		}
	}
	return result
}
