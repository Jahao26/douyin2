package middleware

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

func Automiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		tokenstring := c.Query("token") // 获取加密的鉴权token
		if tokenstring == "" {
			tokenstring = c.PostForm("token")
		}
		token, claim, err := ParseToken(tokenstring) // 解析token
		if err != nil || !token.Valid {
			c.JSON(http.StatusOK, gin.H{"code": 401, "msg": "权限不足"})
			c.Abort() // 确保这个请求的其他函数不会被调用，例如router中的第二个handlefunc
			return
		}

		c.Set("uid", claim.UserId)
		c.Next()
	}
}
