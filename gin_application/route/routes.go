package route

import (
	"GinLearning/gin_application/controller"
	"GinLearning/gin_application/middleware"
	"github.com/gin-gonic/gin"
)

func CollectRoute(r *gin.Engine) *gin.Engine {
	r.Use(middleware.CORSMiddleware(), middleware.RecoveryMiddleware())
	r.POST("/api/auth/register", controller.Register)                     // 注册
	r.POST("/api/auth/login", controller.Login)                           //  登录
	r.GET("/api/auth/info", middleware.AuthMiddleware(), controller.Info) // 传递数据
	return r
}
