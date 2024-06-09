package controller

import (
	"errors"
	"regexp"
	"strconv"
	"strings"

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

func isEmailValid(fl validator.FieldLevel) bool {
	email := fl.Field().String()
	email = strings.TrimSpace(email) // Remove any leading or trailing spaces
	// Improved regex for email validation
	regex := regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	return regex.MatchString(email)
}

// ValidateLatitude checks if a latitude value is valid.
func ValidateLatitude(lat float64) error {
	if lat < -90 || lat > 90 {
		return errors.New("invalid latitude: must be between -90 and 90")
	}
	return nil
}

// ValidateLongitude checks if a longitude value is valid.
func ValidateLongitude(long float64) error {
	if long < -180 || long > 180 {
		return errors.New("invalid longitude: must be between -180 and 180")
	}
	return nil
}

// ValidateLatLong checks if both latitude and longitude values are valid.
func ValidateLatLong(latStr, longStr string) error {
	lat, err := strconv.ParseFloat(latStr, 64)
	if err != nil {
		return errors.New("invalid latitude: must be a number")
	}

	long, err := strconv.ParseFloat(longStr, 64)
	if err != nil {
		return errors.New("invalid longitude: must be a number")
	}

	if err := ValidateLatitude(lat); err != nil {
		return err
	}
	if err := ValidateLongitude(long); err != nil {
		return err
	}
	return nil
}

// isValidUUID checks if a string is a valid UUID.
func isValidUUID(fl validator.FieldLevel) bool {
	uuidRegex := `^[0-9a-fA-F]{8}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{4}\-[0-9a-fA-F]{12}$`
	match, _ := regexp.MatchString(uuidRegex, fl.Field().String())
	return match
}
