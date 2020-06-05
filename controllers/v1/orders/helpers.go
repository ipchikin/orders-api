package orders

import (
	"errors"
	"regexp"

	"github.com/gin-gonic/gin"
)

func abortWithErrorJSON(c *gin.Context, code int, message string) {
	c.Error(errors.New(message))
	c.AbortWithStatusJSON(code, gin.H{"error": message})
}

// validateLatitude
func validateLatitude(lat string) bool {
	matched, err := regexp.MatchString(`^(\+|-)?(?:90(?:(?:\.0+)?)|(?:[0-9]|[1-8][0-9])(?:(?:\.[0-9]+)?))$`, lat)
	if err != nil {
		return false
	}

	return matched
}

// validateLongitude
func validateLongitude(long string) bool {
	matched, err := regexp.MatchString(`^(\+|-)?(?:180(?:(?:\.0+)?)|(?:[0-9]|[1-9][0-9]|1[0-7][0-9])(?:(?:\.[0-9]+)?))$`, long)
	if err != nil {
		return false
	}

	return matched
}
