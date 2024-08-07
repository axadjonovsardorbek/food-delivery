package handler

import (
	"context"
	"gateway/genproto/courier"
	cp "gateway/genproto/order"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// OrderAssign handles assign an existing order.
// @Summary Assign order
// @Description Assign an existing order
// @Tags admin
// @Accept json
// @Produce json
// @Param id query string false "Id"
// @Param courier_id query string false "CourierId"
// @Success 200 {object} string "Assigned order"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Order not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /admin/assign/order [post]
func (h *Handler) AssignOrder(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	courier_id := c.Query("courier_id")

	order := cp.OrderUpdateReq{
		Id:        c.Query("id"),
		CourierId: courier_id,
	}

	_, err := h.srvs.Order.Update(context.Background(), &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't assign order", "details": err.Error()})
		return
	}

	_, err = h.srvs.Notification.Create(context.Background(), &courier.NotificationCreateReq{
		UserId:  courier_id,
		Message: "Your order has been attached",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Assigned order"})
}

// Orders handles getting all order.
// @Summary Get all order
// @Description Get all order
// @Tags admin
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Success 200 {object} cp.OrderGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /admin/orders [get]
func (h *Handler) Orders(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "admin" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	req := cp.OrderGetAllReq{
		Filter: &cp.Filter{},
	}

	pageStr := c.Query("page")
	var page int
	var err error
	if pageStr == "" {
		page = 1
	} else {
		page, err = strconv.Atoi(pageStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
	}

	filter := cp.Filter{
		Page: int32(page),
	}

	req.Filter.Page = filter.Page

	res, err := h.srvs.Order.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}
