package config

import (
    "log"

    "dicom-service/pkg/storage"
)

// Load initializes the configuration for the application.
func Load() {
    // Perform necessary configuration loading here
    // Example: loading environment variables or configuration files

    // Initialize storage
    if err := storage.Initialize(); err != nil {
        log.Fatalf("Failed to initialize storage: %v", err)
    }
}
