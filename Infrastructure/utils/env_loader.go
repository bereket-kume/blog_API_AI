package utils

import (
	"bufio"
	"log"
	"os"
	"strings"
)

// LoadEnv loads environment variables from .env file
func LoadEnv() {
	// Check if .env file exists
	if _, err := os.Stat(".env"); os.IsNotExist(err) {
		log.Println("No .env file found, using system environment variables")
		return
	}

	// Open .env file
	file, err := os.Open(".env")
	if err != nil {
		log.Printf("Error opening .env file: %v", err)
		return
	}
	defer file.Close()

	// Read file line by line
	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())

		// Skip empty lines and comments
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}

		// Split on first '=' to separate key and value
		parts := strings.SplitN(line, "=", 2)
		if len(parts) != 2 {
			continue
		}

		key := strings.TrimSpace(parts[0])
		value := strings.TrimSpace(parts[1])

		// Remove quotes if present
		if len(value) >= 2 && (value[0] == '"' && value[len(value)-1] == '"') {
			value = value[1 : len(value)-1]
		}

		// Set environment variable if not already set
		if os.Getenv(key) == "" {
			os.Setenv(key, value)
		}
	}

	if err := scanner.Err(); err != nil {
		log.Printf("Error reading .env file: %v", err)
	}

	log.Println("Environment variables loaded from .env file")
}
