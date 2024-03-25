package middleware

// import (
// 	"log"
// 	"net/http"

// 	"github.com/gin-gonic/gin"
// 	"github.com/rahulgupta0409/authentication-service/helpers"
// )

// func AuthMiddleware(next http.Handler) http.Handler {
// 	return http.HandlerFunc(func(response http.ResponseWriter, request *http.Request) {
// 		// Do stuff here

// 		clientToken := request.Header.Get("token")

// 		if clientToken == "" {
// 			c.JSON(http.StatusInternalServerError, gin.H{"error": "No Authorization Header Provided"})
// 			response.Write([]byte("No Authorization Header Provided"))
// 			return
// 		}

// 		claims, err := helpers.ValidateToken(clientToken)
// 		if err != "" {
// 			response.Write([]byte(err))
// 			c.Abort()
// 			return
// 		}
// 		response.Set("email", claims.Email)
// 		response.WriteHeader()
// 		c.Set("firstname", claims.FirstName)
// 		c.Set("lastname", claims.LastName)
// 		c.Set("userid", claims.UserId)
// 		c.Set("usertype", claims.UserType)
// 		c.Next()

// 		log.Println(request.RequestURI)
// 		// Call the next handler, which can be another middleware in the chain, or the final handler.
// 		next.ServeHTTP(response, request)
// 	})
// }
