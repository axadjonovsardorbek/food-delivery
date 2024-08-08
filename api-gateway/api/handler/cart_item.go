package handler

import (
	"context"
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
// @Tags cart item
// @Accept json
// @Produce json
// @Param event body cp.CartItemCreateReq true "Cart data"
// @Success 200 {object} string "Cart item created"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart_item [post]
func (h *Handler) CartItemCreate(c *gin.Context) {
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

	cart_id, err := h.srvs.CartItem.GetCartId(context.Background(), &cp.ById{Id: id})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users cart not found"})
		log.Println("error: ", err)
		return
	}

	var req cp.CartItemCreateReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	req.UserId = id
	req.CartId = cart_id.Id

	_, err = h.srvs.CartItem.Create(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	amount, err := h.srvs.CartItem.GetTotalAmount(context.Background(), &cp.GetTotalAmountReq{
		UserId: req.UserId,
		CartId: req.CartId,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get price", "details": err.Error()})
		return
	}

	_, err = h.srvs.Cart.Update(context.Background(), &cp.CartUpdateReq{
		Id:          req.CartId,
		UserId:      req.UserId,
		TotalAmount: amount.TotalAmount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update cart", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart item created"})
}

// CartGetById handles the get a cart.
// @Summary Get cart
// @Description Get a cart
// @Tags cart item
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Success 200 {object} cp.CartItemGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart_item/{id} [get]
func (h *Handler) CartItemGetById(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.CartItem.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get cart item", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// CartGetAll handles getting all carts.
// @Summary Get all carts
// @Description Get all carts
// @Tags cart item
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Success 200 {object} cp.CartItemGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart_item/all [get]
func (h *Handler) CartItemGetAll(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := claims.(jwt.MapClaims)["user_id"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	cart_id, err := h.srvs.CartItem.GetCartId(context.Background(), &cp.ById{Id: userId})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"message": "cart is empty"})
		log.Println("error: ", err)
		return
	}

	req := cp.CartItemGetAllReq{
		UserId: userId,
		CartId: cart_id.Id,
		Filter: &cp.Filter{},
	}

	pageStr := c.Query("page")
	var page int
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

	res, err := h.srvs.CartItem.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get cart items", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// CartUpdate handles updating an existing cart.
// @Summary Update cart
// @Description Update an existing cart
// @Tags cart item
// @Accept json
// @Produce json
// @Param id path string true "Cart ID"
// @Param quantity query integer false "Quantity"
// @Success 200 {object} string "Cart updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Cart not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart_item/{id} [put]
func (h *Handler) CartItemUpdate(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	user_id := claims.(jwt.MapClaims)["user_id"].(string)
	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "user" {
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	id := c.Param("id")

	cart_id, err := h.srvs.CartItem.GetCartId(context.Background(), &cp.ById{Id: user_id})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users cart not found"})
		log.Println("error: ", err)
		return
	}

	quantityStr := c.Query("quantity")
	var quantity int
	if quantityStr == "" {
		quantity = 1
	} else {
		quantity, err = strconv.Atoi(quantityStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
			return
		}
	}

	_, err = h.srvs.CartItem.Update(context.Background(), &cp.CartItemUpdateReq{
		Id:       id,
		UserId:   user_id,
		Quantity: int32(quantity),
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update cart", "details": err.Error()})
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
		Id:          cart_id.Id,
		UserId:      user_id,
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
// @Tags cart item
// @Accept json
// @Produce json
// @Param id query string false "Cart ID"
// @Success 200 {object} string "Cart deleted"
// @Failure 400 {object} string "Invalid event ID"
// @Failure 404 {object} string "Cart not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /cart_item/{id} [delete]
func (h *Handler) CartItemDelete(c *gin.Context) {
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

	cart_id, err := h.srvs.CartItem.GetCartId(context.Background(), &cp.ById{Id: user_id})

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "users cart not found"})
		log.Println("error: ", err)
		return
	}

	id := &cp.ById{Id: c.Query("id")}

	_, err = h.srvs.CartItem.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete cart", "details": err.Error()})
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
		Id:          cart_id.Id,
		UserId:      user_id,
		TotalAmount: amount.TotalAmount,
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update cart", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Cart deleted"})
}
