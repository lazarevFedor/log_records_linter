package pkg

import (
	"testing"
)

// TestIsLowercaseStartValid tests the isLowercaseStartValid function
// to ensure it correctly identifies messages that start with a lowercase letter
func TestIsLowercaseStartValid(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		expectError bool
	}{
		{
			name:        "lowercase_start",
			msg:         "fedya",
			expectError: false,
		},
		{
			name:        "uppercase_start",
			msg:         "Fedya",
			expectError: true,
		},
		{
			name:        "number_start",
			msg:         "21 Fedya",
			expectError: false,
		},
		{
			name:        "space_lowercase",
			msg:         " fedya ",
			expectError: false,
		},
		{
			name:        "space_uppercase",
			msg:         " Fedya ",
			expectError: true,
		},
		{
			name:        "empty_string",
			msg:         "",
			expectError: false,
		},
		{
			name:        "only_spaces",
			msg:         "  ",
			expectError: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := isLowercaseStartValid(tt.msg)

			if tt.expectError && !result {
				t.Errorf("expected error, but got none: %s", tt.name)
			}
			if !tt.expectError && result {
				t.Errorf("expected no error, but got one: %s", tt.name)
			}
		})
	}
}

// TestIsEnglishOnlyValid tests the isEnglishOnlyValid function
// to ensure it correctly identifies messages that contain only English text
func TestIsEnglishOnlyValid(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		expectError bool
	}{
		{
			name:        "english",
			msg:         "fedor is a good golang developer",
			expectError: false,
		},
		{
			name:        "english_numbers",
			msg:         "fedor is 21 years old",
			expectError: false,
		},
		{
			name:        "english_punctuation",
			msg:         "fedor write code: clean, efficient, and well-documented.",
			expectError: false,
		},
		{
			name:        "cyrillic",
			msg:         "Федя плохой и использовал кириллицу в сообщении",
			expectError: true,
		},
		{
			name:        "chinese",
			msg:         "中文",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := !(isEnglishOnlyValid(tt.msg))

			if tt.expectError && !result {
				t.Errorf("expected error, but got none: %s", tt.name)
			}
			if !tt.expectError && result {
				t.Errorf("expected no error, but got one: %s", tt.name)
			}
		})
	}
}

// TestIsNoSpecialCharsValid tests the isNoSpecialCharsValid function
// to ensure it correctly identifies messages that contain only allowed characters
func TestIsNoSpecialCharsValid(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		expectError bool
	}{
		{
			name:        "no_special_chars",
			msg:         "selectel is the best company I ever know",
			expectError: false,
		},
		{
			name:        "allowed_punctuation",
			msg:         "message with. commas, colons: semicolons; dashes- and underscores_",
			expectError: false,
		},
		{
			name:        "allowed_quotes",
			msg:         `message with 'single' and "double" quotes`,
			expectError: false,
		},
		{
			name:        "exclamation_mark",
			msg:         "do not scream, I am not deaf!",
			expectError: true,
		},
		{
			name:        "question_mark",
			msg:         "i wanted to ask you about?",
			expectError: true,
		},
		{
			name:        "ellipsis",
			msg:         "tough…",
			expectError: true,
		},
		{
			name:        "emoji",
			msg:         "hi 😀",
			expectError: true,
		},
		{
			name:        "at_symbol",
			msg:         "catmail@dogmail.com",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := isNoSpecialCharsValid(tt.msg)

			if tt.expectError && !result {
				t.Errorf("expected error, but got none: %s", tt.name)
			}
			if !tt.expectError && result {
				t.Errorf("expected no error, but got one: %s", tt.name)
			}
		})
	}
}

// TestIsNoSensitiveDataValid tests the isNoSensitiveDataValid function to ensure
// it correctly identifies messages containing sensitive data.
func TestIsNoSensitiveDataValid(t *testing.T) {
	tests := []struct {
		name        string
		msg         string
		expectError bool
	}{
		{
			name:        "clean_message",
			msg:         "all good there",
			expectError: false,
		},
		{
			name:        "password_keyword",
			msg:         "password is incorrect",
			expectError: true,
		},
		{
			name:        "passwd_keyword",
			msg:         "user passwd expired",
			expectError: true,
		},
		{
			name:        "pwd_keyword",
			msg:         "pwd changed successfully",
			expectError: true,
		},
		{
			name:        "token_keyword",
			msg:         "invalid token provided",
			expectError: true,
		},
		{
			name:        "api_key_keyword",
			msg:         "api_key is missing",
			expectError: true,
		},
		{
			name:        "secret_keyword",
			msg:         "secret not found",
			expectError: true,
		},
		{
			name:        "private_key_keyword",
			msg:         "private_key is required",
			expectError: true,
		},
		{
			name:        "access_key_keyword",
			msg:         "access_key validation failed",
			expectError: true,
		},
		{
			name:        "client_secret_keyword",
			msg:         "client_secret is invalid",
			expectError: true,
		},
		{
			name:        "bearer_keyword",
			msg:         "bearer token required",
			expectError: true,
		},
		{
			name:        "password_uppercase",
			msg:         "PASSWORD is incorrect",
			expectError: true,
		},
		{
			name:        "jwt_token",
			msg:         "token eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIn0.dozjgNryP4J3jVmNHl0w5N_XgL0n3I9PlFUP0THsR8U",
			expectError: true,
		},
		{
			name:        "github_pat",
			msg:         "token ghp_1234567890123456789012345678901234",
			expectError: true,
		},
		{
			name:        "uuid_pattern",
			msg:         "id 550e8400-e29b-41d4-a716-446655440000",
			expectError: true,
		},
		{
			name:        "bearer_token_pattern",
			msg:         "auth bearer='abcd1234efgh5678ijkl9012'",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			t.Parallel()
			result := isNoSensitiveDataValid(tt.msg)

			if tt.expectError && !result {
				t.Errorf("expected error, but got none: %s", tt.name)
			}
			if !tt.expectError && result {
				t.Errorf("expected no error, but got one: %s", tt.name)
			}
		})
	}
}
