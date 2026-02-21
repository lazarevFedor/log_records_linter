package pkg

import (
	"go/token"
	"regexp"
	"strings"
	"unicode"

	"golang.org/x/tools/go/analysis"
)

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

func checkLowercaseStart(pass *analysis.Pass, pos token.Pos, msg string) {
	msg = strings.TrimSpace(msg)
	if msg == "" {
		return
	}

	for _, r := range msg {
		if unicode.IsLetter(r) {
			if unicode.IsUpper(r) {
				pass.Reportf(pos, "log message should start with lowercase letter")
			}
			return
		}
	}
}

func checkEnglishOnly(pass *analysis.Pass, pos token.Pos, msg string) {
	for _, r := range msg {
		if !unicode.IsLetter(r) {
			continue
		}

		if !((r >= 'A' && r <= 'Z') || (r >= 'a' && r <= 'z')) {
			pass.Reportf(pos, "log message should be in English only (Latin alphabet)")
			return
		}
	}
}

// checkNoSpecialChars проверяет, что сообщение не содержит недопустимые спецсимволы и эмодзи.
func checkNoSpecialChars(pass *analysis.Pass, pos token.Pos, msg string) {
	for _, r := range msg {
		if unicode.IsLetter(r) || unicode.IsDigit(r) || unicode.IsSpace(r) {
			continue
		}

		if r == '.' || r == ',' || r == ':' || r == ';' || r == '-' || r == '_' || r == '\'' || r == '"' {
			continue
		}

		if r == '!' || r == '?' || r == '…' || r > 127 {
			pass.Reportf(pos, "log message should not contain special characters or emoji")
			return
		}
	}
}

func checkNoSensitiveData(pass *analysis.Pass, pos token.Pos, msg string) {
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
			pass.Reportf(pos, "log message should not contain sensitive data keywords like %q", keyword)
			return
		}
	}

	for _, sp := range secretPatterns {
		if sp.pattern.MatchString(msg) {
			pass.Reportf(pos, "log message appears to contain secret data (%s)", sp.name)
			return
		}
	}
}
