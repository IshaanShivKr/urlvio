package auth

import (
	"net/http"

	"github.com/clerk/clerk-sdk-go/v2"
	clerkhttp "github.com/clerk/clerk-sdk-go/v2/http"
	"github.com/gin-gonic/gin"
)

const userIDKey = "userID"

func Middleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		next := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			claims, ok := clerk.SessionClaimsFromContext(r.Context())
			if !ok {
				c.AbortWithStatus(http.StatusUnauthorized)
				return
			}

			c.Set(userIDKey, claims.Subject)
			c.Next()
		})

		clerkhttp.WithHeaderAuthorization()(next).ServeHTTP(c.Writer, c.Request)
	}
}

func UserID(c *gin.Context) (string, bool) {
	value, exists := c.Get(userIDKey)
	if !exists {
		return "", false
	}

	userID, ok := value.(string)
	return userID, ok && userID != ""
}
