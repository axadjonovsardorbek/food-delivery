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

// OrderAccept handles accept an existing order.
// @Summary Accept order
// @Description Accept an existing order
// @Tags courier
// @Accept json
// @Produce json
// @Param id query string false "Id"
// @Param courier_id query string false "CourierId"
// @Success 200 {object} string "Accepted order"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Order not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/accept [post]
func (h *Handler) AcceptOrder(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)
	id := claims.(jwt.MapClaims)["user_id"].(string)

	if role != "courier" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	order := cp.OrderUpdateReq{
		Id:        c.Query("id"),
		CourierId: id,
	}

	_, err := h.srvs.Order.Update(context.Background(), &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't accept order", "details": err.Error()})
		return
	}

	_, err = h.srvs.Notification.Create(context.Background(), &courier.NotificationCreateReq{
		UserId:  id,
		Message: "Your order has been attached",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accepted order"})
}

// OrderUpdate handles update an existing order.
// @Summary Update order
// @Description Update an existing order
// @Tags courier
// @Accept json
// @Produce json
// @Param id query string false "Id"
// @Param status query string false "Status"
// @Success 200 {object} string "Accepted order"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Order not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/status/{id} [put]
func (h *Handler) UpdateOrder(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "courier" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	order := cp.OrderUpdateReq{
		Id:     c.Query("id"),
		Status: c.Query("status"),
	}

	_, err := h.srvs.Order.Update(context.Background(), &order)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't accept order", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Accepted order"})
}

// OrderGetAll handles getting all order.
// @Summary Get all order
// @Description Get all order
// @Tags courier
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Param status query string false "Status"
// @Success 200 {object} cp.OrderGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/order/history [get]
func (h *Handler) OrdersHistory(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)

	req := cp.OrderGetAllReq{
		CourierId: id,
		Status:    c.Query("status"),
		Filter:    &cp.Filter{},
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
