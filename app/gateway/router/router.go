package router

import (
	"github.com/gin-gonic/gin"
	"github.com/leexsh/go-todolist/app/gateway/http"
	middleware1 "github.com/leexsh/go-todolist/app/gateway/middleware"
)

func NewRouter() *gin.Engine {
	ginRouter := gin.Default()
	ginRouter.Use(middleware1.ErrorMiddleware(), middleware1.JWT())
	store := cookie.NewStore([]byte("something-very-secret"))
	ginRouter.Use(sessions.Sessions("mysession", store))

	v1 := ginRouter.Group("/api/v1")
	{
		v1.GET("ping", func(context *gin.Context) {
			context.JSON(200, "success")
		})
		// 用户服务
		v1.POST("/user/register", http.UserRegister)
		v1.POST("/user/login", http.UserLogin)

		// 需要登录保护
		authed := v1.Group("/")
		authed.Use(middleware1.JWT())
		{
			// 任务模块
			authed.GET("task", http.GetTaskList)
			authed.POST("task", http.CreateTask)
			authed.PUT("task", http.UpdateTask)
			authed.DELETE("task", http.DeleteTask)
		}
	}
	return ginRouter
}
