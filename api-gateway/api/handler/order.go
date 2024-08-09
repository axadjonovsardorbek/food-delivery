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

// OrderCreate handles the creation of a new order.
// @Summary Create order
// @Description Create a new order
// @Tags order
// @Accept json
// @Produce json
// @Param order body cp.OrderCreateReq true "Order data"
// @Success 200 {object} string "Order created"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /order [post]
func (h *Handler) OrderCreate(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	var req cp.OrderCreateReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	cart_id, err := h.srvs.CartItem.GetCartId(context.Background(), &cp.ById{Id: id})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users cart not found"})
		log.Println("error: ", err)
		return
	}

	amount, err := h.srvs.CartItem.GetTotalAmount(context.Background(), &cp.GetTotalAmountReq{
		UserId: id,
		CartId: cart_id.Id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get price", "details": err.Error()})
		return
	}

	cart_items, err := h.srvs.CartItem.GetAll(context.Background(), &cp.CartItemGetAllReq{UserId: id, CartId: cart_id.Id})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get cart items", "details": err.Error()})
		return
	}

	if amount.TotalAmount <= 0 {
		c.JSON(http.StatusBadRequest, gin.H{"message": "Cart is empty"})
		return
	}

	req.UserId = id
	req.TotalAmount = amount.TotalAmount

	order_id, err := h.srvs.Order.Create(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	if order_id.Id == "" || order_id.Id == "string" {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "invalid order id"})
		return
	}

	for i := 0; i < len(cart_items.CartItems); i++ {
		item := cart_items.CartItems[i]

		order_item := cp.OrderItemCreateReq{
			OrderId:   order_id.Id,
			ProductId: item.ProductId,
			Quantity:  item.Quantity,
		}

		_, err = h.srvs.OrderItem.Create(context.Background(), &order_item)

		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			log.Println("error: ", err)
			return
		}
	}

	_, err = h.srvs.Cart.Delete(context.Background(), cart_id)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	_, err = h.srvs.Cart.Update(context.Background(), &cp.CartUpdateReq{
		Id:          cart_id.Id,
		UserId:      req.UserId,
		TotalAmount: amount.TotalAmount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update cart", "details": err.Error()})
		return
	}

	_, err = h.srvs.Notification.Create(context.Background(), &courier.NotificationCreateReq{
		UserId:  id,
		Message: "Your order has been confirmated",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order created"})
}

// OrderGetById handles the get a order.
// @Summary Get order
// @Description Get a order
// @Tags order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} cp.OrderGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /order/{id} [get]
func (h *Handler) OrderGetById(c *gin.Context) {
	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.Order.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get order", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// OrderGetAll handles getting all order.
// @Summary Get all order
// @Description Get all order
// @Tags order
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Param status query string false "Status"
// @Success 200 {object} cp.OrderGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /order/all [get]
func (h *Handler) OrderGetAll(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)

	req := cp.OrderGetAllReq{
		UserId: id,
		Status: c.Query("status"),
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

// OrderDelete handles deleting a order by ID.
// @Summary Delete order
// @Description Delete a order by ID
// @Tags order
// @Accept json
// @Produce json
// @Param id path string true "Order ID"
// @Success 200 {object} string "Order deleted"
// @Failure 400 {object} string "Invalid order ID"
// @Failure 404 {object} string "Order not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /order/{id} [delete]
func (h *Handler) OrderDelete(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)
	user_id := claims.(jwt.MapClaims)["user_id"].(string)

	if role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	id := &cp.ById{Id: c.Param("id")}
	_, err := h.srvs.Order.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete order", "details": err.Error()})
		return
	}

	_, err = h.srvs.Notification.Create(context.Background(), &courier.NotificationCreateReq{
		UserId:  user_id,
		Message: "Your order has been cancelled",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Order deleted"})
}
