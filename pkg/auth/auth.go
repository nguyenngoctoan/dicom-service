package auth

import (
    "net/http"
    "github.com/gin-gonic/gin"
)

// Authentication middleware to validate access
func AuthMiddleware() gin.HandlerFunc {
    return func(c *gin.Context) {
        // Allow unauthenticated access to the health check endpoint
        if c.Request.URL.Path == "/health" {
            c.Next()
            return
        }

        token := c.GetHeader("Authorization")
        if token == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
            return
        }

        user, role := validateToken(token)
        if user == "" {
            c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
            return
        }
        c.Set("user", user)
        c.Set("role", role)
        c.Next()
    }
}

func validateToken(token string) (user, role string) {
    // Simulate token validation
    if token == "valid-token" {
        return "admin", "admin"
    }
    return "", ""
}
