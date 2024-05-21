package storage

import (
    "crypto/sha256"
    "encoding/hex"
    "errors"
    "io/ioutil"
    "os"
    "path/filepath"
)

const storagePath = "storage" // Define the directory for storage

// Initialize sets up necessary directories for storage.
func Initialize() error {
    if err := os.MkdirAll(storagePath, 0755); err != nil {
        return err
    }
    return nil
}

// SaveFile saves a file from a provided path to the storage directory.
// Returns a unique file identifier for retrieving the file.
func SaveFile(filePath string) (string, error) {
    data, err := ioutil.ReadFile(filePath)
    if err != nil {
        return "", err
    }

    // Generate a unique identifier based on the file content
    fileID := generateFileID(data)
    savePath := filepath.Join(storagePath, fileID)

    // Save the file data to the new location
    if err := ioutil.WriteFile(savePath, data, 0644); err != nil {
        return "", err
    }

    return fileID, nil
}

// RetrieveFile retrieves a file based on its file ID.
func RetrieveFile(fileID string) (string, error) {
    filePath := filepath.Join(storagePath, fileID)
    if _, err := os.Stat(filePath); os.IsNotExist(err) {
        return "", errors.New("file not found")
    }
    return filePath, nil
}

// generateFileID creates a unique identifier for a file using SHA-256 hashing.
func generateFileID(data []byte) string {
    hasher := sha256.New()
    hasher.Write(data)
    return hex.EncodeToString(hasher.Sum(nil))
}
