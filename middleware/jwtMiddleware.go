package middleware

//
//import (
//	"GoMJTrainingCamp/service"
//	"fmt"
//
//	"github.com/gin-gonic/gin"
//	"net/http"
//	"strings"
//)
//
//// JWTAuthMiddleware checks the presence of the JWT token and validates it
//func JWTAuthMiddleware() gin.HandlerFunc {
//	return func(c *gin.Context) {
//		// Get the token from the Authorization header
//		tokenString := c.GetHeader("Authorization")
//		if tokenString == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"message": "Token is required"})
//			c.Abort()
//			return
//		}
//
//		// Remove the "Bearer " prefix
//		tokenString = strings.TrimPrefix(tokenString, "Bearer ")
//		if tokenString == "" {
//			c.JSON(http.StatusUnauthorized, gin.H{"message": "Bearer token is required"})
//			c.Abort()
//			return
//		}
//		fmt.Println(tokenString)
//		// Validate the JWT token
//		token, err := service.WithJWTAuth(tokenString)
//		if err != nil || !token.Valid {
//			c.JSON(http.StatusUnauthorized, gin.H{"message": "Invalid token"})
//			c.Abort()
//			return
//		}
//
//		// Proceed with the request
//		c.Next()
//	}
//}
