package goenvsubst_test

import (
	"fmt"
	"os"

	"github.com/iamolegga/goenvsubst"
)

// ExampleDo demonstrates basic usage with a struct
func ExampleDo() {
	// Set up environment variables
	os.Setenv("DATABASE_URL", "postgres://localhost:5432/mydb")
	os.Setenv("API_KEY", "secret123")
	defer func() {
		os.Unsetenv("DATABASE_URL")
		os.Unsetenv("API_KEY")
	}()

	config := &struct {
		DatabaseURL string
		APIKey      string
		Debug       bool
	}{
		DatabaseURL: "$DATABASE_URL",
		APIKey:      "$API_KEY",
		Debug:       true,
	}

	err := goenvsubst.Do(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Database URL: %s\n", config.DatabaseURL)
	fmt.Printf("API Key: %s\n", config.APIKey)
	fmt.Printf("Debug: %t\n", config.Debug)

	// Output:
	// Database URL: postgres://localhost:5432/mydb
	// API Key: secret123
	// Debug: true
}

// ExampleDo_slice demonstrates usage with slices
func ExampleDo_slice() {
	// Set up environment variables
	os.Setenv("SERVICE1", "auth-service")
	os.Setenv("SERVICE2", "payment-service")
	defer func() {
		os.Unsetenv("SERVICE1")
		os.Unsetenv("SERVICE2")
	}()

	services := []string{"$SERVICE1", "$SERVICE2", "static-service"}

	err := goenvsubst.Do(&services)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, service := range services {
		fmt.Printf("Service %d: %s\n", i+1, service)
	}

	// Output:
	// Service 1: auth-service
	// Service 2: payment-service
	// Service 3: static-service
}

// ExampleDo_map demonstrates usage with maps
func ExampleDo_map() {
	// Set up environment variables
	os.Setenv("REDIS_URL", "redis://localhost:6379")
	os.Setenv("MONGO_URL", "mongodb://localhost:27017")
	defer func() {
		os.Unsetenv("REDIS_URL")
		os.Unsetenv("MONGO_URL")
	}()

	config := map[string]string{
		"redis":    "$REDIS_URL",
		"mongodb":  "$MONGO_URL",
		"static":   "https://cdn.example.com",
		"$ENV_KEY": "this-key-wont-change", // Keys are not substituted
	}

	err := goenvsubst.Do(&config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Redis: %s\n", config["redis"])
	fmt.Printf("MongoDB: %s\n", config["mongodb"])
	fmt.Printf("Static: %s\n", config["static"])

	// Output:
	// Redis: redis://localhost:6379
	// MongoDB: mongodb://localhost:27017
	// Static: https://cdn.example.com
}

// ExampleDo_nestedStructure demonstrates usage with complex nested structures
func ExampleDo_nestedStructure() {
	// Set up environment variables
	os.Setenv("DB_HOST", "localhost")
	os.Setenv("DB_PORT", "5432")
	os.Setenv("API_ENDPOINT", "https://api.example.com")
	defer func() {
		os.Unsetenv("DB_HOST")
		os.Unsetenv("DB_PORT")
		os.Unsetenv("API_ENDPOINT")
	}()

	type DatabaseConfig struct {
		Host string
		Port string
	}

	type AppConfig struct {
		Database  DatabaseConfig
		Services  []string
		Endpoints map[string]string
	}

	config := &AppConfig{
		Database: DatabaseConfig{
			Host: "$DB_HOST",
			Port: "$DB_PORT",
		},
		Services: []string{"$SERVICE1", "static-service"},
		Endpoints: map[string]string{
			"api":    "$API_ENDPOINT",
			"health": "/health",
		},
	}

	err := goenvsubst.Do(config)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Database Host: %s\n", config.Database.Host)
	fmt.Printf("Database Port: %s\n", config.Database.Port)
	fmt.Printf("API Endpoint: %s\n", config.Endpoints["api"])
	fmt.Printf("Health Endpoint: %s\n", config.Endpoints["health"])
	if config.Services[0] == "" {
		fmt.Printf("Service 1:\n")
	} else {
		fmt.Printf("Service 1: %s\n", config.Services[0])
	}
	fmt.Printf("Service 2: %s\n", config.Services[1])

	// Output:
	// Database Host: localhost
	// Database Port: 5432
	// API Endpoint: https://api.example.com
	// Health Endpoint: /health
	// Service 1:
	// Service 2: static-service
}

// ExampleDo_pointers demonstrates usage with pointers
func ExampleDo_pointers() {
	// Set up environment variables
	os.Setenv("SECRET_KEY", "my-secret-key")
	defer os.Unsetenv("SECRET_KEY")

	// Pointer to string
	secretValue := "$SECRET_KEY"
	err := goenvsubst.Do(&secretValue)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Secret: %s\n", secretValue)

	// Nil pointer (safe to process)
	var nilPtr *string
	err = goenvsubst.Do(&nilPtr)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	fmt.Printf("Nil pointer handled safely: %v\n", nilPtr == nil)

	// Output:
	// Secret: my-secret-key
	// Nil pointer handled safely: true
}

// ExampleDo_array demonstrates usage with arrays
func ExampleDo_array() {
	// Set up environment variables
	os.Setenv("SERVER1", "server1.example.com")
	os.Setenv("SERVER2", "server2.example.com")
	defer func() {
		os.Unsetenv("SERVER1")
		os.Unsetenv("SERVER2")
	}()

	servers := [3]string{"$SERVER1", "$SERVER2", "server3.example.com"}

	err := goenvsubst.Do(&servers)
	if err != nil {
		fmt.Printf("Error: %v\n", err)
		return
	}

	for i, server := range servers {
		fmt.Printf("Server %d: %s\n", i+1, server)
	}

	// Output:
	// Server 1: server1.example.com
	// Server 2: server2.example.com
	// Server 3: server3.example.com
}
