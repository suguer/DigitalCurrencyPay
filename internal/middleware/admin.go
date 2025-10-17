package middleware

import "github.com/gin-gonic/gin"

func AdminMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		user_id, exists := c.Get("user_id")
		if !exists {
			c.JSON(401, gin.H{"error": "未登录"})
			c.Abort()
			return
		}
		if user_id.(uint) != 1 {
			c.JSON(403, gin.H{"error": "权限不足"})
			c.Abort()
			return
		}
	}
}
