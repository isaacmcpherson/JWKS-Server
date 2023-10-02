package main

import (
	"crypto/rand"
	"crypto/rsa"
	"fmt"
	"net/http"
	"time"

	"github.com/dgrijalva/jwt-go"
	"github.com/gin-gonic/gin"
)

var (
	privateKey *rsa.PrivateKey
	publicKey  *rsa.PublicKey
)

func initKeys() error {
	// Generate RSA key pair
	privateKey, err := rsa.GenerateKey(rand.Reader, 2048)
	if err != nil {
		return err
	}
	publicKey = &privateKey.PublicKey
	return nil
}

func main() {
	if err := initKeys(); err != nil {
		fmt.Printf("Failed to generate keys: %v\n", err)
		return
	}

	r := gin.Default()

	// JWKS endpoint
	r.GET("/jwks", func(c *gin.Context) {
		jwks := gin.H{
			"keys": []gin.H{
				{
					"kid": "kid1",
					"kty": "RSA",
					"n":   publicKey.N.String(),
					"e":   publicKey.E,
					"use": "sig",
					"alg": "RS256",
				},
			},
		}
		c.JSON(http.StatusOK, jwks)
	})

	// Auth endpoint
	r.POST("/auth", func(c *gin.Context) {
		username := c.PostForm("username")
		password := c.PostForm("password")

		// Mock user authentication (replace with real authentication logic)
		if username == "userABC" && password == "password123" {
			// Check if the "expired" query parameter is present
			expired := c.Query("expired") == "true"

			// Set the key and expiry accordingly
			var signingKey *rsa.PrivateKey
			expiry := time.Now().Add(time.Hour) // Set an expiry time

			if expired {
				signingKey = privateKey             // Use the expired private key
				expiry = time.Now().Add(-time.Hour) // Set an expired expiry time
			} else {
				signingKey = privateKey // Use the current private key
			}

			// Create a JWT token
			token := jwt.New(jwt.SigningMethodRS256)
			claims := token.Claims.(jwt.MapClaims)
			claims["sub"] = username
			claims["exp"] = expiry.Unix()
			tokenString, err := token.SignedString(signingKey)
			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": "Error creating JWT"})
				return
			}

			c.JSON(http.StatusOK, gin.H{"token": tokenString})
		} else {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid username or password"})
		}
	})

	r.GET("/", func(c *gin.Context) {
		c.String(http.StatusOK, "Hello, World!")
	})

	// Run the web server on port 8080
	r.Run(":8080")
}
