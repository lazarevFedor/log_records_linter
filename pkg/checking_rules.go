package pkg

import (
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

// secretPatterns contains regex patterns for common secrets that should not be included in log messages.
var secretPatterns = []struct {
	pattern *regexp.Regexp
	name    string
}{
	// JWT Tokens
	{regexp.MustCompile(`eyJ[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}\.[A-Za-z0-9_-]{10,}`), "JWT token"},

	// GitHub
	{regexp.MustCompile(`ghp_[0-9a-zA-Z]{36}`), "GitHub Personal Access Token"},
	{regexp.MustCompile(`gho_[0-9a-zA-Z]{36}`), "GitHub OAuth Access Token"},
	{regexp.MustCompile(`ghr_[0-9a-zA-Z]{36}`), "GitHub Refresh Token"},

	// Private Keys
	{regexp.MustCompile(`-----BEGIN (RSA |EC |DSA |OPENSSH |PGP )?PRIVATE KEY( BLOCK)?-----`), "Private Key"},

	// Generic patterns
	{regexp.MustCompile(`[0-9a-f]{8}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{4}-[0-9a-f]{12}`), "UUID (potential secret)"},

	// Tokens
	{regexp.MustCompile(`(?i)(bearer|token)['"]?\s*[:=]\s*['"]?[0-9a-zA-Z\-_.]{20,}`), "Bearer/Auth Token"},
}

// isLowercaseStartValid checks if the log message starts with a lowercase letter.
func isLowercaseStartValid(msg string) bool {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return false
	}

	if unicode.IsLetter(rune(msg[0])) {
		return unicode.IsUpper(rune(msg[0]))
	}
	return false
}

// isEnglishOnlyValid checks if the log message contains only English letters, digits, spaces, and allowed punctuation.
func isEnglishOnlyValid(msg string) bool {
	for _, r := range msg {
		if unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		if r == '.' || r == ',' || r == ':' || r == ';' || r == '-' || r == '_' || r == '\'' || r == '"' {
			continue
		}
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			continue
		}
		return true
	}
	return false
}

// isNoSpecialCharsValid checks if the log message contains any special characters or emoji that are not allowed.
func isNoSpecialCharsValid(msg string) bool {
	for _, r := range msg {
		if unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}
		if r == '.' || r == ',' || r == ':' || r == ';' || r == '-' || r == '_' || r == '\'' || r == '"' {
			continue
		}
		if (r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z') {
			continue
		}
		return true
	}
	return false
}

// isNoSensitiveDataValid checks if the log message contains any sensitive data based on keywords and regex patterns.
func isNoSensitiveDataValid(msg string) bool {
	lowerMsg := strings.ToLower(msg)
	sensitiveKeywords := []string{
		"password", "passwd", "pwd",
		"token", "api_key", "apikey",
		"secret", "private_key", "privatekey",
		"access_key", "accesskey",
		"client_secret", "clientsecret",
		"bearer",
	}

	for _, keyword := range sensitiveKeywords {
		if strings.Contains(lowerMsg, keyword) {
			return true
		}
	}
	for _, sp := range secretPatterns {
		if sp.pattern.MatchString(msg) {
			return true
		}
	}
	return false
}

// checkLowercaseStart checks if the log message starts with a lowercase letter
// and reports an issue if it does not.
func checkLowercaseStart(pass *analysis.Pass, pos token.Pos, msg string) {
	if isLowercaseStartValid(msg) {
		pass.Reportf(pos, "log message should start with lowercase letter")
	}
}

// checkEnglishOnly checks if the log message contains only English letters, digits,
// spaces, and allowed punctuation, and reports an issue if it does not.
func checkEnglishOnly(pass *analysis.Pass, pos token.Pos, msg string) {
	if isEnglishOnlyValid(msg) {
		pass.Reportf(pos, "log message should be in English only (Latin alphabet)")
	}
}

// checkNoSpecialChars checks if the log message contains any special characters
// or emoji that are not allowed, and reports an issue if it does.
func checkNoSpecialChars(pass *analysis.Pass, pos token.Pos, msg string) {
	if isNoSpecialCharsValid(msg) {
		pass.Reportf(pos, "log message should not contain special characters or emoji")
	}
}

// checkNoSensitiveData checks if the log message contains any sensitive data based on keywords
// and regex patterns, and reports an issue if it does.
func checkNoSensitiveData(pass *analysis.Pass, pos token.Pos, msg string) {
	if isNoSensitiveDataValid(msg) {
		pass.Reportf(pos, "log message should not contain sensitive data")
	}
}
