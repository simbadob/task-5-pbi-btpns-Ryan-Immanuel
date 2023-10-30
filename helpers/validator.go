package helpers

import (
	"regexp"
	"strings"
)

func IsURLValidWithGivenExtensions(url string, desiredExtensions []string) bool {
	urlPattern := `^(https?://)?([a-zA-Z0-9.-]+(\.[a-zA-Z]{2,})+(/.*)?)?$`

	regex, err := regexp.Compile(urlPattern)
	if err != nil {
		return false
	}

	parts := strings.Split(url, "/")
	if len(parts) == 0 {
		return false
	}

	filename := parts[len(parts)-1]

	if !regex.MatchString(url) {
		return false
	}

	for _, ext := range desiredExtensions {
		if ContainsTargetExtension(filename, ext) {
			return true
		}
	}

	return false
}

func ContainsTargetExtension(filename string, desiredExtension string) bool {
	parts := strings.Split(filename, ".")
	if len(parts) < 2 {
		return false
	}

	extension := parts[len(parts)-1]

	return strings.EqualFold(extension, desiredExtension)
}
