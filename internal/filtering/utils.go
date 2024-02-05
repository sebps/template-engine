package filtering

import (
	"regexp"
	"strings"
)

// TODO: rewrite jsonPathRegex regexp
var jsonPathRegex = regexp.MustCompile(`^\$\[\?\((?P<path>@[^\s]+)\s[^\s]+\s[^\s]+\]$`)
var jsonPathPathPartRegex = regexp.MustCompile(`\@[^\s]+`)

func RewriteJsonPathFilterPathPart(jsonPathFilter string, additive string) string {
	finalJsonPathFilter := jsonPathFilter
	jsonPathFilterPath := GetJsonPathPathPart(jsonPathFilter)
	jsonPathFilterPathParts := strings.Split(jsonPathFilterPath, ".")
	newJsonPathFilterParts := make([]string, 0)
	for i, jsonPathFilterPathPart := range jsonPathFilterPathParts {
		newJsonPathFilterParts = append(newJsonPathFilterParts, jsonPathFilterPathPart)
		if i == 0 && jsonPathFilterPathPart == "@" {
			newJsonPathFilterParts = append(newJsonPathFilterParts, additive)
		}
	}
	finalJsonPathFilterPathPart := strings.Join(newJsonPathFilterParts, ".")
	finalJsonPathFilter = ReplaceJsonPathPathPart(finalJsonPathFilter, finalJsonPathFilterPathPart)

	return finalJsonPathFilter
}

// TODO: check regexp match
func IsJsonPathCompliant(input string) bool {
	return true
	// return jsonPathRegex.Match([]byte(input))
}

func GetJsonPathPathPart(input string) string {
	matches := jsonPathRegex.FindStringSubmatch(input)
	if len(matches) > 1 {
		return matches[1]
	}

	return ""
}

func ReplaceJsonPathPathPart(input string, replacment string) string {
	if IsJsonPathCompliant(input) {
		return jsonPathPathPartRegex.ReplaceAllString(input, replacment)
	}

	return input
}
