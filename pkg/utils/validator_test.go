package utils_test

import (
	"github.com/HunterGooD/voice_friend_user_service/pkg/utils"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestValidateEmail(t *testing.T) {
	tests := []struct {
		email   string
		isValid bool
	}{
		{"test@example.com", true},
		{"invalid-email", false},
		{"user@domain", false},
		{"user@domain.co.uk", true},
	}

	for _, tt := range tests {
		require.Equal(t, tt.isValid, utils.ValidateEmail(tt.email), "Not valid email %s", tt.email)
	}
}

func TestValidatePhone(t *testing.T) {
	tests := []struct {
		phone   string
		isValid bool
	}{
		{"+1234567890", true},
		{"1234567890", true},
		{"+1-800-555-5555", true},
		{"invalid-phone", false},
	}

	for _, tt := range tests {
		require.Equal(t, tt.isValid, utils.ValidatePhone(tt.phone), "Not valid phone %s", tt.phone)
	}
}
