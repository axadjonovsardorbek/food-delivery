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

// CartCreate handles the creation of a new cart.
// @Summary Create cart
// @Description Create a new cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} string "Cart created"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart [post]
func (h *Handler) CartCreate(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	id := claims.(jwt.MapClaims)["user_id"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user"{
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	var req cp.CartCreateReq

	req.UserId = id

	_, err := h.srvs.Cart.Create(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	_, err = h.srvs.Notification.Create(context.Background(), &courier.NotificationCreateReq{
		UserId: id,
		Message: "Cart created for you",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart created"})
}

// CartGetById handles the get a cart.
// @Summary Get cart
// @Description Get a cart
// @Tags cart
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} cp.CartGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart/{id} [get]
func (h *Handler) CartGetById(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user"{
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.Cart.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get cart", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// CartGetAll handles getting all carts.
// @Summary Get all carts
// @Description Get all carts
// @Tags cart
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Success 200 {object} cp.CartGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart/all [get]
func (h *Handler) CartGetAll(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := claims.(jwt.MapClaims)["user_id"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user"{
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	req := cp.CartGetAllReq{
		UserId: userId,
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

	res, err := h.srvs.Cart.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get carts", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// CartUpdate handles updating an existing cart.
// @Summary Update cart
// @Description Update an existing cart
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} string "Cart updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Cart not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart/{id} [put]
func (h *Handler) CartUpdate(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user_id := claims.(jwt.MapClaims)["user_id"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user"{
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	cart_id, err := h.srvs.CartItem.GetCartId(context.Background(), &cp.ById{Id: user_id})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users cart not found"})
		log.Println("error: ", err)
		return
	}

	amount, err := h.srvs.CartItem.GetTotalAmount(context.Background(), &cp.GetTotalAmountReq{
		UserId: user_id,
		CartId: cart_id.Id,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get price", "details": err.Error()})
		return
	}

	_, err = h.srvs.Cart.Update(context.Background(), &cp.CartUpdateReq{
		Id: cart_id.Id,
		UserId: user_id,
		TotalAmount: amount.TotalAmount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update cart", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart updated"})
}

// CartDelete handles deleting a cart by ID.
// @Summary Delete cart
// @Description Delete a cart by ID
// @Tags cart
// @Accept json
// @Produce json
// @Success 200 {object} string "Cart deleted"
// @Failure 400 {object} string "Invalid event ID"
// @Failure 404 {object} string "Cart not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart/{id} [delete]
func (h *Handler) CartDelete(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)
	user_id := claims.(jwt.MapClaims)["user_id"].(string)

	if role != "user"{
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	cart_id, err := h.srvs.CartItem.GetCartId(context.Background(), &cp.ById{Id: user_id})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users cart not found"})
		log.Println("error: ", err)
		return
	}

	_, err = h.srvs.Cart.Delete(context.Background(), cart_id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete cart", "details": err.Error()})
		return
	}

	_, err = h.srvs.Cart.Update(context.Background(), &cp.CartUpdateReq{
		Id: cart_id.Id,
		UserId: user_id,
		TotalAmount: 0,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update cart", "details": err.Error()})
		return
	}

	_, err = h.srvs.Notification.Create(context.Background(), &courier.NotificationCreateReq{
		UserId: user_id,
		Message: "Your cart has been emptied",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted"})
}
