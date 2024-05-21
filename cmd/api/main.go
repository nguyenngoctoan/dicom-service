package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "dicom-service/pkg/dicomhandler"
    "dicom-service/internal/config"
    "dicom-service/pkg/auth"
)

func main() {
    config.Load()
    router := gin.Default()

    // Apply authentication middleware globally, except for the health check endpoint
    router.GET("/health", dicomhandler.HealthCheck)  // Handle /health without authentication

    // Apply authentication middleware to all other routes
    router.Use(auth.AuthMiddleware())
    
    router.POST("/upload", dicomhandler.UploadDICOM)
    router.GET("/attribute", dicomhandler.GetDICOMAttribute)
    router.GET("/convert", dicomhandler.ConvertDICOMToPNG)

    log.Fatal(router.Run(":8080"))
}
