package controller

import (
	"regexp"

	"github.com/go-playground/validator/v10"
)

func customURL(fl validator.FieldLevel) bool {
	// Regular expression pattern for a valid URL
	// This pattern requires the URL to start with http:// or https://
	// followed by a valid domain name
	pattern := `^(http|https)://[a-zA-Z0-9\-\.]+\.[a-zA-Z]{2,}(\/\S*)?$`

	// Compile the regular expression
	regex := regexp.MustCompile(pattern)

	// Match the URL against the regular expression
	url := fl.Field().String()
	return regex.MatchString(url)
}
