/*
Package goenvsubst provides functionality to recursively replace environment variable
references in Go data structures with their actual values from the environment.

The package supports various Go data types including structs, slices, maps, arrays,
and pointers, both as top-level inputs and nested within other structures.

Environment variables should be referenced in the format $VAR_NAME. If an environment
variable is not set or is empty, it will be replaced with an empty string.

# Basic Usage

The main function Do() accepts any Go data structure and modifies it in-place:

	import "github.com/iamolegga/goenvsubst"

	// Set environment variable
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/mydb")

	config := &struct {
		DatabaseURL string
		Debug       bool
	}{
		DatabaseURL: "$DATABASE_URL",
		Debug:       true,
	}

	err := goenvsubst.Do(config)
	if err != nil {
		log.Fatal(err)
	}
	// config.DatabaseURL is now "postgres://localhost:5432/mydb"

# Struct Fields

Environment variable substitution works with any string field in a struct:

	type Config struct {
		Host     string
		Port     string
		Username string
		Password string
	}

	config := &Config{
		Host:     "$DB_HOST",
		Port:     "$DB_PORT",
		Username: "$DB_USER",
		Password: "$DB_PASS",
	}

	goenvsubst.Do(config)

# Slices and Arrays

String elements in slices and arrays are processed:

	// Slice example
	urls := []string{"$API_URL", "$BACKUP_URL", "https://static.example.com"}
	goenvsubst.Do(&urls)

	// Array example
	servers := [3]string{"$SERVER1", "$SERVER2", "$SERVER3"}
	goenvsubst.Do(&servers)

# Maps

Map values (but not keys) are processed for environment variable substitution:

	config := map[string]string{
		"database_url": "$DATABASE_URL",
		"redis_url":    "$REDIS_URL",
		"api_key":      "$API_KEY",
	}
	goenvsubst.Do(&config)

# Nested Structures

The package handles deeply nested structures:

	type DatabaseConfig struct {
		URL      string
		Username string
		Password string
	}

	type AppConfig struct {
		Database DatabaseConfig
		Services []string
		Env      map[string]string
	}

	config := &AppConfig{
		Database: DatabaseConfig{
			URL:      "$DATABASE_URL",
			Username: "$DB_USER",
			Password: "$DB_PASS",
		},
		Services: []string{"$SERVICE1", "$SERVICE2"},
		Env: map[string]string{
			"LOG_LEVEL": "$LOG_LEVEL",
			"DEBUG":     "$DEBUG_MODE",
		},
	}

	goenvsubst.Do(config)

# Pointers

The package safely handles pointers, including nil pointers:

	var config *struct {
		Value string
	}

	// Safe to call with nil pointer
	goenvsubst.Do(&config) // No error, no operation

	// With actual pointer
	config = &struct{ Value string }{"$MY_VALUE"}
	goenvsubst.Do(config)

# Complex Example

A real-world configuration structure:

	type ServerConfig struct {
		Host string `json:"host"`
		Port string `json:"port"`
	}

	type AppConfig struct {
		Server    ServerConfig            `json:"server"`
		Database  string                  `json:"database"`
		Redis     string                  `json:"redis"`
		Services  []string                `json:"services"`
		Features  map[string]bool         `json:"features"`
		Endpoints map[string]string       `json:"endpoints"`
		Secrets   map[string]*string      `json:"secrets"`
	}

	// Set environment variables
	os.Setenv("APP_HOST", "0.0.0.0")
	os.Setenv("APP_PORT", "8080")
	os.Setenv("DATABASE_URL", "postgres://localhost/myapp")
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("API_ENDPOINT", "https://api.example.com")

	config := &AppConfig{
		Server: ServerConfig{
			Host: "$APP_HOST",
			Port: "$APP_PORT",
		},
		Database: "$DATABASE_URL",
		Redis:    "$REDIS_URL",
		Services: []string{"$SERVICE_AUTH", "$SERVICE_PAYMENT"},
		Features: map[string]bool{
			"feature_a": true,
			"feature_b": false,
		},
		Endpoints: map[string]string{
			"api":     "$API_ENDPOINT",
			"webhook": "$WEBHOOK_URL",
		},
		Secrets: map[string]*string{
			"jwt_secret": func() *string { s := "$JWT_SECRET"; return &s }(),
		},
	}

	err := goenvsubst.Do(config)
	if err != nil {
		log.Fatal(err)
	}

# Error Handling

The Do function returns an error if there are issues during processing.
Currently, the function is designed to be robust and typically returns nil,
but error handling is provided for future extensibility:

	err := goenvsubst.Do(config)
	if err != nil {
		log.Printf("Failed to substitute environment variables: %v", err)
		return
	}

# Important Notes

- Only string values are processed for environment variable substitution
- Map keys are never modified, only values
- Missing or empty environment variables are replaced with empty strings
- The function modifies the input data structure in-place
- Nil pointers are handled safely without causing panics
- The function is safe for concurrent use as it doesn't modify global state
*/
package goenvsubst
