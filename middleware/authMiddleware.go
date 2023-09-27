package middleware

import (
	"fmt"
	"maisonsport/dao"
	"net/http"

	"github.com/gin-gonic/gin"
)

/**
*	中间件，从请求头中拿 token，然后通过 token  从 redis 里查出对应的  user_id
*
 */
func AuthMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		token := c.GetHeader("Authorization")

		// Check if token is provided
		if token == "" {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Authentication token is required"})
			c.Abort()
			return
		}

		// Retrieve user_id from Redis based on the token
		res := dao.RedisDB.Get(token)
		if res.Err() != nil {
			c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid token"})
			c.Abort()
			return
		}

		// 把 user_id  传下去
		c.Set("user_id", res.Val())
		fmt.Println("中间件在执行", res.Val())

		// Continue with the next handler
		c.Next()
	}
}
