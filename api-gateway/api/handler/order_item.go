package handler

import (
	"context"
	cp "gateway/genproto/order"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// OrderItemGetById handles the get a item.
// @Summary Get item
// @Description Get a item
// @Tags order
// @Accept json
// @Produce json
// @Param id path string true "OrderItem ID"
// @Success 200 {object} cp.OrderItemGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /order/item/{id} [get]
func (h *Handler) OrderItemGetById(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role == "courier" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.OrderItem.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get item", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// OrderItemGetAll handles getting all item.
// @Summary Get all item
// @Description Get all item
// @Tags order
// @Accept json
// @Produce json
// @Param order_id query string false "OrderId"
// @Success 200 {object} cp.OrderItemGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /order/item/all [get]
func (h *Handler) OrderItemGetAll(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role == "courier" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	req := cp.OrderItemGetAllReq{
		OrderId: c.Query("order_id"),
	}

	res, err := h.srvs.OrderItem.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get items", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
