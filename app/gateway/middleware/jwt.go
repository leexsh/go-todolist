package middleware1

import (
	"time"

	"github.com/gin-gonic/gin"
	"github.com/leexsh/go-todolist/pkg/ctl"
	"github.com/leexsh/go-todolist/pkg/myerr"
	"github.com/leexsh/go-todolist/util/jwt"
)

func JWT() gin.HandlerFunc {
	return func(c *gin.Context) {
		var code int
		var data interface{}
		code = 200
		token := c.GetHeader("Authorization")
		if token == "" {
			code = 404
			c.JSON(200, gin.H{
				"status": code,
				"msg":    myerr.GetMsg(code),
				"data":   data,
			})
			c.Abort()
		}
		claims, err := jwt.ParseToken(token)
		if err != nil {
			code = myerr.ErrorAuthCheckTokenFail
		} else if time.Now().Unix() > claims.ExpiresAt {
			code = myerr.ErrorAuthCheckTokenTimeout
		}
		if code != myerr.SUCCESS {
			c.JSON(200, gin.H{
				"status": code,
				"msg":    myerr.GetMsg(code),
				"data":   data,
			})
			c.Abort()
			return
		}
		c.Request = c.Request.WithContext(ctl.NewContext(c.Request.Context(), &ctl.UserInfo{Id: claims.UserID}))
		ctl.InitUserInfo(c.Request.Context())
		c.Next()
	}
}
