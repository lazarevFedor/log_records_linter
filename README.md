# Log Records Linter

Линтер для Go, совместимый с **golangci-lint**, проверяющий форматирование лог-записей в коде.

## Линтер предназначен для проверки лог-сообщений по следующим правилам:

1. Лог-сообщения должны начинаться со строчной буквы
    ```go
    //❌Неправильно
    log.Info("Starting server on port 8080")
    slog.Error("Failed to connect to database")
    
    //✅Правильно
    log.Info("starting server on port 8080")
    slog.Error("failed to connect to database")
    ```

2. Лог-сообщения должны быть только на английском языке
   ```go
    //❌Неправильно
    log.Info("запуск сервера")
    log.Error("ошибка подключения к базе данных")
    
    //✅Правильно
    log.Info("starting server")
    log.Error("failed to connect to database")
   ```

3. Лог-сообщения не должны содержать спецсимволы или эмодзи
   ```go
    //❌Неправильно
    log.Info("server started!🚀")
    log.Error("connection failed!!!")
    log.Warn("warning: something went wrong...")
    
    //✅Правильно
    log.Info("server started")
    log.Error("connection failed")
    log.Warn("something went wrong")
   ```

4. Лог-сообщения не должны содержать потенциально чувствительные данные
   ```go
    //❌Неправильно
    log.Info("user password: " + password)
    log.Debug("api_key=" + apiKey)
    log.Info("token: " + token)
    
    //✅Правильно
    log.Info("user authenticated successfully")
    log.Debug("api request completed")
    log.Info("token validated")
   ```

## Поддерживаемые логгеры

Линтер работает со следующими логирующими библиотеками:
- `log/slog`
- `go.uber.org/zap`

## Инструкции по использованию

### Как CLI инструмент (через go vet)

```bash
# Сборка CLI инструмента
make build-cli

# Запуск линтера на текущем проекте
make vet

# Запуск линтера с поддержкой Suggested Fixes
make vet-fix
```

### Как плагин golangci-lint через Module Plugin System

```bash
# Сборка custom-gcl с плагином
make build-plugin

# Запустить линтера на текущем проекте
make lint

# Запуск линтера с поддержкой Suggested Fixes
make lint-fix
```

## Конфигурация

### Файл конфига правил

По умолчанию используется `configs/config_rules.json`:

```json
{
  "enable_lowercase_start": true,
  "enable_english_only": true,
  "enable_no_special_chars": true,
  "enable_sensitive_patterns": true
}
```

**Параметры:**
- `enable_lowercase_start` — проверять на строчную букву в начале
- `enable_english_only` — проверять на английский язык
- `enable_no_special_chars` — проверять на отсутствие спецсимволов
- `enable_sensitive_patterns` — проверять на чувствительные данные

### Конфиг golangci-lint

При использовании плагина через golangci-lint, конфигурируйте в `.golangci.yaml`:

```yaml
version: "2"

linters:
  default: none
  enable:
    - logs

  settings:
    custom:
      logs:
        type: module
        description: "Checks log records correct formatting"
        original-url: "github.com/lazarevFedor/log_records_linter"
        settings:
          config: "configs/config_rules.json"
```

### Конфиг сборки плагина

Содержимое `.custom-gcl.yml`:

```yaml
version: v2.10.1

plugins:
  - module: 'log_records_linter'
    import: 'log_records_linter/pkg'
    path: '.'
```

## Команды Makefile

```bash
make help             # Показать справку по командам
make build-cli        # Собрать CLI инструмент (logs)
make vet              # Проверить код через go vet
make vet-fix          # Проверить и исправить через go vet
make build-plugin     # Собрать custom-gcl с плагином
make lint             # Запустить проверку через custom-gcl
make lint-fix         # Запустить проверку и исправления
make test             # Запустить тесты
```

## Опции запуска

### Для CLI версии

```bash
# С конфигом по умолчанию
./logs ./...

# С кастомным конфигом
./logs -config=/path/to/config.json ./...
```

### Для плагина golangci-lint

```bash
# Обычный запуск линтера
./custom-gcl run

# Запуск с подробным выводом
./custom-gcl run -v

# Запуск с Suggested Fixes
./custom-gcl run --fix

# Запуск линтера на конкретном файле
./custom-gcl run ./testdata/test_logger.go
```

## Тестирование

```bash
make test
```

## Структура проекта

```
log_records_linter/
├── cmd/
│   └── main.go                # CLI точка входа
├── pkg/
│   ├── analyzer.go            # Основной анализатор
│   ├── checking_rules.go      # Правила проверки
│   ├── config.go              # Работа с конфигом
│   ├── analyzer_test.go       # Тесты анализатора
│   └── checking_rules_test.go # Тесты правил
├── configs/
│   └── config_rules.json      # Конфиг правил
├── testdata/
│   ├── test_logger.go         # Примеры для тестов
│   └── src/                   # Мок внешней библиотеки для тестов
├── .golangci.yaml             # Конфиг golangci-lint
├── .custom-gcl.yml            # Конфиг для сборки плагина
├── Makefile                   # Команды для сборки и тестирования
├── go.mod                     # Зависимости
└── README.md                  # Документация
```
## Примеры использования

### Пример 1: Проверка существующего проекта

```bash
# Клонируйте или перейдите в директорию проекта
cd /path/to/your/project

# Соберите кастомный golangci-lint с плагином
make build-plugin

# Запустите проверку всего проекта
./custom-gcl run
```

**Пример вывода:**
```
testdata/test_logger.go:12:15: log message should start with lowercase letter (logs)
        slogger.Info("Invalid message starting with uppercase")
                     ^
testdata/test_logger.go:24:15: log message should be in English only (logs)
        slogger.Info("message with русский text")
                     ^
testdata/test_logger.go:36:15: log message should not contain special characters or emoji (logs)
        slogger.Info("message with exclamation!")
                     ^
```

### Пример 2: Автоматическое исправление ошибок

Линтер поддерживает автоматические исправления для некоторых правил:

```bash
# Запустите с флагом --fix для автоматических исправлений
./custom-gcl run --fix

# Или через Makefile
make lint-fix
```

**До исправления:**
```go
package main

import "log/slog"

func main() {
    logger := slog.Default()
    logger.Info("Server Started on port 8080")
    logger.Error("Connection failed!!!")
    logger.Warn("Ошибка подключения")
}
```

**После исправления:**
```go
package main

import "log/slog"

func main() {
    logger := slog.Default()
    logger.Info("server Started on port 8080")
    logger.Error("connection failed")
    logger.Warn("")
}
```