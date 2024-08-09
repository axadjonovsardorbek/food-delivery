package api

import (
	"gateway/api/handler"
	"gateway/api/middleware"

	_ "gateway/api/docs"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(h *handler.Handler) *gin.Engine {

	router := gin.Default()

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected := router.Group("/", middleware.JWTMiddleware())

	courier := protected.Group("/courier")
	courier.POST("/accept", h.AcceptOrder)
	courier.PUT("/status/:id", h.UpdateOrder)
	courier.GET("/order/history", h.OrdersHistory)

	admin := protected.Group("/admin")
	admin.POST("/assign/order", h.AssignOrder)
	task := admin.Group("/task")
	task.POST("/", h.TaskCreate)
	task.GET("/:id", h.TaskGetById)
	task.GET("/all", h.TaskGetAll)
	task.PUT("/:id", h.TaskUpdate)
	task.DELETE("/:id", h.TaskDelete)

	cart := protected.Group("/cart")
	cart.POST("/", h.CartCreate)
	cart.GET("/:id", h.CartGetById)
	cart.GET("/all", h.CartGetAll)
	cart.PUT("/:id", h.CartUpdate)
	cart.DELETE("/:id", h.CartDelete)

	cart_item := protected.Group("/cart_item")
	cart_item.POST("/", h.CartItemCreate)
	cart_item.GET("/:id", h.CartItemGetById)
	cart_item.GET("/all", h.CartItemGetAll)
	cart_item.PUT("/:id", h.CartItemUpdate)
	cart_item.DELETE("/:id", h.CartItemDelete)

	order := protected.Group("/order")
	order.POST("/", h.OrderCreate)
	order.GET("/:id", h.OrderGetById)
	order.GET("/all", h.OrderGetAll)
	order.DELETE("/:id", h.OrderDelete)

	order_item := protected.Group("/order_item")
	order_item.GET("/:id", h.OrderItemGetById)
	order_item.GET("/all", h.OrderItemGetAll)

	product := protected.Group("/product")
	product.POST("/", h.ProductCreate)
	product.GET("/:id", h.ProductGetById)
	product.GET("/all", h.ProductGetAll)
	product.PUT("/:id", h.ProductUpdate)
	product.DELETE("/:id", h.ProductDelete)

	notification := protected.Group("/notification")
	notification.POST("/", h.NotificationCreate)
	notification.GET("/:id", h.NotificationGetById)
	notification.GET("/all", h.NotificationGetAll)
	notification.PUT("/:id", h.NotificationUpdate)
	notification.DELETE("/:id", h.NotificationDelete)

	return router
}
