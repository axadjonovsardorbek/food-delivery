package api

import (
	"gateway/api/handler"
	"gateway/api/middleware"

	"github.com/gin-gonic/gin"
	ginSwagger "github.com/swaggo/gin-swagger"
	swaggerFiles "github.com/swaggo/files"
	_ "gateway/api/docs"
)

// @securityDefinitions.apikey BearerAuth
// @in header
// @name Authorization
func NewRouter(h *handler.Handler) *gin.Engine {

	router := gin.Default()

	router.GET("/api/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	protected := router.Group("/", middleware.JWTMiddleware())

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

	product := protected.Group("/product")
	product.POST("/", h.ProductCreate)
	product.GET("/:id", h.ProductGetById)
	product.GET("/all", h.ProductGetAll)
	product.PUT("/:id", h.ProductUpdate)
	product.DELETE("/:id", h.ProductDelete)

	task := protected.Group("/task")
	task.POST("/", h.TaskCreate)
	task.GET("/:id", h.TaskGetById)
	task.GET("/all", h.TaskGetAll)
	task.PUT("/:id", h.TaskUpdate)
	task.DELETE("/:id", h.TaskDelete)

	notification := protected.Group("/notification")
	notification.POST("/", h.NotificationCreate)
	notification.GET("/:id", h.NotificationGetById)
	notification.GET("/all", h.NotificationGetAll)
	notification.PUT("/:id", h.NotificationUpdate)
	notification.DELETE("/:id", h.NotificationDelete)

	// event := protected.Group("/timeline/custom-event")
	// event.POST("/", h.EventCreate)
	// event.GET("/:id", h.EventGetById)
	// event.GET("/all", h.EventGetAll)
	// event.PUT("/:id", h.EventUpdate)
	// event.DELETE("/:id", h.EventDelete)

	// milestone := protected.Group("/timeline/milestone")
	// milestone.POST("/", h.MilestoneCreate)
	// milestone.GET("/:id", h.MilestoneGetById)
	// milestone.GET("/all", h.MilestoneGetAll)
	// milestone.PUT("/:id", h.MilestoneUpdate)
	// milestone.DELETE("/:id", h.MilestoneDelete)

	// historical := protected.Group("/timeline/historical")
	// historical.POST("/", h.HistoricalEventCreate)
	// historical.GET("/:id", h.HistoricalEventGetById)
	// historical.GET("/all", h.HistoricalEventGetAll)
	// historical.PUT("/:id", h.HistoricalEventUpdate)
	// historical.DELETE("/:id", h.HistoricalEventDelete)

	// personal := protected.Group("/timeline/personal")
	// personal.POST("/", h.PersonalEventCreate)
	// personal.GET("/:id", h.PersonalEventGetById)
	// personal.GET("/all", h.PersonalEventGetAll)
	// personal.PUT("/:id", h.PersonalEventUpdate)
	// personal.DELETE("/:id", h.PersonalEventDelete)
	
	// context := protected.Group("/timeline/context")
	// context.GET("/:date", h.Context)

	return router
}
