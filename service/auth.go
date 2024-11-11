package service

import (
	"GoMJTrainingCamp/dbs/models/users"
	"GoMJTrainingCamp/utils"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey contextKey = "userID"

// WithJWTAuth is a middleware function for handling JWT-based authentication
func WithJWTAuth(handlerFunc http.HandlerFunc) http.HandlerFunc {

	return func(w http.ResponseWriter, r *http.Request) {
		// Step 1: Retrieve the token string from the request
		tokenString := utils.GetTokenFromRequest(r)

		if tokenString == "" {
			log.Println("No token found in the request")
			permissionDenied(w)
			return
		}
		log.Printf("Token extracted: %s", tokenString)

		// Step 2: Remove "Bearer " prefix if present
		if strings.HasPrefix(tokenString, "Bearer ") {
			tokenString = strings.TrimPrefix(tokenString, "Bearer ")
			log.Printf("Token after removing 'Bearer' prefix: %s", tokenString)
		}

		// Step 3: Validate the token
		token, err := validateJWT(tokenString)
		if err != nil {
			log.Printf("Failed to validate token: %v", err)
			permissionDenied(w)
			return
		}
		log.Println("Token validated successfully")

		// Step 4: Check if token is valid
		if !token.Valid {
			log.Println("Token is not valid")
			permissionDenied(w)
			return
		}

		// Step 5: Extract claims from token
		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			log.Println("Failed to parse token claims")
			permissionDenied(w)
			return
		}
		log.Println("Token claims extracted successfully")

		// Step 6: Retrieve userID from claims
		str, ok := claims["userID"].(string)
		if !ok {
			log.Println("userID claim is missing or not a string in token claims")
			permissionDenied(w)
			return
		}
		log.Printf("Extracted userID from token claims: %s", str)

		// Step 7: Convert userID to integer
		userID, err := strconv.Atoi(str)
		if err != nil {
			log.Printf("Failed to convert userID to int: %v", err)
			permissionDenied(w)
			return
		}
		log.Printf("Converted userID to int: %d", userID)

		// Step 8: Look up the user in the database
		user, err := models.GetUserByID(uint(userID)) // Pass db along with userID
		if err != nil {
			log.Printf("Failed to get user by id: %v", err)
			permissionDenied(w)
			return
		}
		log.Printf("User found with ID %d", user.IDUser)

		// Step 9: Add the user ID to the request context
		ctx := r.Context()
		ctx = context.WithValue(ctx, UserKey, user.IDUser)
		r = r.WithContext(ctx)
		log.Println("User ID added to request context")

		// Step 10: Proceed with the original handler function
		handlerFunc(w, r)
		log.Println("Handler function executed successfully")
	}
}

// CreateJWT generates a JWT token with a specified expiration and userID claim
func CreateJWT(secret []byte, userID int) (string, error) {
	expiration := time.Second * time.Duration(3600*24*7)
	log.Printf("Creating JWT with expiration: %d seconds", 3600*24*7)

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"userID":    strconv.Itoa(userID),
		"expiresAt": time.Now().Add(expiration).Unix(),
	})

	tokenString, err := token.SignedString(secret)
	if err != nil {
		log.Printf("Error signing token: %v", err)
		return "", err
	}

	log.Printf("Generated token: %s", tokenString)
	return tokenString, nil
}

// validateJWT validates the provided JWT token string using the configured secret
func validateJWT(tokenString string) (*jwt.Token, error) {
	log.Printf("Validating JWT with secret: %s", "not-secret-secret-anymore?")

	return jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		return []byte("not-secret-secret-anymore?"), nil
	})
}

// permissionDenied handles permission errors by sending a 403 response
func permissionDenied(w http.ResponseWriter) {
	log.Println("Permission denied")
	// Send a consistent error response using the generic APIResponse format
	response := utils.APIResponse{
		Status:  http.StatusForbidden,
		Message: "Permission denied",
		Data:    nil, // No data in error responses
	}

	// Write the response in JSON format
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusForbidden)
	json.NewEncoder(w).Encode(response)
}
