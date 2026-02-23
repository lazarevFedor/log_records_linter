package testdata

import (
	"log/slog"

	"go.uber.org/zap"
)

func testLowercaseStart() {
	slogger := slog.Default()
	slogger.Info("valid message starting with lowercase")   // ok
	slogger.Info("Invalid message starting with uppercase") // want "log message should start with lowercase letter"
	slogger.Info("123 starts with number")                  // ok

	zlogger, _ := zap.NewProduction()
	zlogger.Info("valid message starting with lowercase")   // ok
	zlogger.Info("Invalid message starting with uppercase") // want "log message should start with lowercase letter"
	zlogger.Info("123 starts with number")                  // ok
}

func testEnglishOnly() {
	slogger := slog.Default()
	slogger.Info("english only message")      // ok
	slogger.Info("message with русский text") // want "log message should be in English only"
	slogger.Info("message with 中文")           // want "log message should be in English only"

	zlogger, _ := zap.NewProduction()
	zlogger.Info("english only message")      // ok
	zlogger.Info("message with русский text") // want "log message should be in English only"
	zlogger.Info("message with 中文")           // want "log message should be in English only"
}

func testSpecialChars() {
	slogger := slog.Default()
	slogger.Info("message with dots and commas.") // ok
	slogger.Info("message with exclamation!")     // want "log message should not contain special characters or emoji"
	slogger.Info("message with @symbol")          // want "log message should not contain special characters or emoji"
	slogger.Info("message with emoji 😀")          // want "log message should not contain special characters or emoji"

	zlogger, _ := zap.NewProduction()
	zlogger.Info("message with dots and commas.") // ok
	zlogger.Info("message with exclamation!")     // want "log message should not contain special characters or emoji"
	zlogger.Info("message with @symbol")          // want "log message should not contain special characters or emoji"
	zlogger.Info("message with emoji 😀")          // want "log message should not contain special characters or emoji"
}

func testSensitiveData() {
	slogger := slog.Default()
	slogger.Info("operation completed successfully") // ok
	slogger.Info("password is incorrect")            // want "log message should not contain sensitive data"
	slogger.Info("token validation failed")          // want "log message should not contain sensitive data"
	slogger.Info("api_key is missing")               // want "log message should not contain sensitive data"
	slogger.Info("invalid bearer token")             // want "log message should not contain sensitive data"

	zlogger, _ := zap.NewProduction()
	zlogger.Info("operation completed successfully") // ok
	zlogger.Info("password is incorrect")            // want "log message should not contain sensitive data"
	zlogger.Info("token validation failed")          // want "log message should not contain sensitive data"
	zlogger.Info("api_key is missing")               // want "log message should not contain sensitive data"
	zlogger.Info("invalid bearer token")             // want "log message should not contain sensitive data"
}
