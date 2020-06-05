package orders

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

// TestValidateLatitude
func TestValidateLatitude(t *testing.T) {
	assert.True(t, validateLatitude("22.281980123456"))
}

// TestValidateLatitudeOutOfRange with out of range latitude
func TestValidateLatitudeOutOfRange(t *testing.T) {
	assert.False(t, validateLatitude("122.281980123456"))
}

// TestValidateLongitude
func TestValidateLongitude(t *testing.T) {
	assert.True(t, validateLongitude("114.161370123456"))
}

// TestValidateLongitudeOutOfRange with out of range longitude
func TestValidateLongitudeOutOfRange(t *testing.T) {
	assert.False(t, validateLongitude("214.161370123456"))
}
