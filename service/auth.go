package service

import (
	"GoMJTrainingCamp/dbs/models/users"
	"GoMJTrainingCamp/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/golang-jwt/jwt/v5"
)

type contextKey string

const UserKey = "userID"

func WithJWTAuth(c *gin.Context) {
	// Step 1: Retrieve the token string from the request
	tokenString := utils.GetTokenFromRequest(c.Request)

	if tokenString == "" {
		log.Println("No token found in the request")
		permissionDenied(c)
		return
	}
	log.Printf("Token extracted: %s", tokenString)

	// Step 2: Remove "Bearer " prefix if present
	if strings.HasPrefix(tokenString, "Bearer ") {
		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
		log.Printf("Token after removing 'Bearer' prefix: %s", tokenString)
	}

	// Step 3: Validate the token
	token, err := ValidateJWT(tokenString)
	if err != nil {
		log.Printf("Failed to validate token: %v", err)
		permissionDenied(c)
		return
	}
	log.Println("Token validated successfully")

	// Step 4: Check if token is valid
	if !token.Valid {
		log.Println("Token is not valid")
		permissionDenied(c)
		return
	}

	// Step 5: Extract claims from token
	claims, ok := token.Claims.(jwt.MapClaims)
	if !ok {
		log.Println("Failed to parse token claims")
		permissionDenied(c)
		return
	}
	log.Println("Token claims extracted successfully")

	// Step 6: Retrieve userID from claims
	str, ok := claims["userID"].(string)
	if !ok {
		log.Println("userID claim is missing or not a string in token claims")
		permissionDenied(c)
		return
	}
	log.Printf("Extracted userID from token claims: %s", str)

	// Step 7: Convert userID to integer
	userID, err := strconv.Atoi(str)
	if err != nil {
		log.Printf("Failed to convert userID to int: %v", err)
		permissionDenied(c)
		return
	}
	log.Printf("Converted userID to int: %d", userID)

	// Step 8: Look up the user in the database
	user, err := models.GetUserByID(uint(userID)) // Pass db along with userID
	if err != nil {
		log.Printf("Failed to get user by id: %v", err)
		permissionDenied(c)
		return
	}
	log.Printf("User found with ID %d", user.IDUser)

	// Step 9: Add the user ID to the request context
	c.Set(UserKey, user.IDUser)
	log.Println("User ID added to request context")

	// Step 10: Proceed with the next middleware/handler
	c.Next()
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
func ValidateJWT(tokenString string) (*jwt.Token, error) {
	log.Println("Validating JWT...")
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		// Ensure the signing method is HMAC
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("unexpected signing method: %v", token.Header["alg"])
		}

		claims, ok := token.Claims.(jwt.MapClaims)
		if !ok {
			return nil, fmt.Errorf("could not parse claims")
		}

		if exp, ok := claims["exp"].(float64); ok {
			expirationTime := time.Unix(int64(exp), 0)
			if expirationTime.Before(time.Now()) {
				return nil, fmt.Errorf("token is expired")
			}
		}
		return []byte("my-secret-key"), nil
	})

	return token, err
}
func permissionDenied(c *gin.Context) {
	log.Println("Permission denied")

	// Create the API response struct
	response := utils.APIResponse{
		Status:  http.StatusForbidden,
		Message: "Permission denied",
		Data:    nil,
	}

	// Set the response content type and send the JSON response
	c.Header("Content-Type", "application/json")
	c.JSON(http.StatusForbidden, response)
}
