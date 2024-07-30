package api

import (
	"github.com/gin-gonic/gin"

	swaggerFiles "github.com/swaggo/files"

	ginSwagger "github.com/swaggo/gin-swagger"

	_ "auth/api/docs"
	"auth/api/handler"
	"auth/api/middleware"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewApi(h *handler.Handler) *gin.Engine {
	router := gin.Default()

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	router.POST("/user/register", h.UserRegister)
	router.POST("/courier/register", h.CourierRegister)

	router.POST("/login", h.Login)

	protected := router.Group("/", middleware.JWTMiddleware())
	protected.GET("/profile", h.Profile)
	protected.GET("/refresh-token", h.RefreshToken)
	protected.DELETE("/profile/delete", h.DeleteProfile)
	protected.PUT("/profile/update", h.UpdateProfile)
	protected.PUT("/change-password", h.ChangePassword)
	protected.POST("/forgot-password", h.ForgotPassword)
	protected.POST("/reset-password", h.ResetPassword)

	return router
}
