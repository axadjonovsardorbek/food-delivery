package handler

import (
	"context"
	"fmt"
	cp "gateway/genproto/order"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// ProductCreate handles the creation of a new product.
// @Summary Create product
// @Description Create a new product
// @Tags product
// @Accept json
// @Produce json
// @Param product body cp.ProductCreateReq true "Product data"
// @Success 200 {object} string "Product created"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /product [post]
func (h *Handler) ProductCreate(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role != "admin"{
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	var req cp.ProductCreateReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}
	
	_, err := h.srvs.Product.Create(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Product created"})
}

// ProductGetById handles the get a product.
// @Summary Get product
// @Description Get a product
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} cp.ProductGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /product/{id} [get]
func (h *Handler) ProductGetById(c *gin.Context) {
	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.Product.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get product", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// ProductGetAll handles getting all product.
// @Summary Get all product
// @Description Get all product
// @Tags product
// @Accept json
// @Produce json
// @Param name query integer false "Name"
// @Param price query integer false "Price"
// @Param page query integer false "Page"
// @Success 200 {object} cp.ProductGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /product/all [get]
func (h *Handler) ProductGetAll(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	role := claims.(jwt.MapClaims)["role"].(string)

	if role == "courier"{
		c.JSON(http.StatusForbidden, gin.H{"error": "This page forbidden for you"})
		return
	}

	priceStr := c.Query("price")

	var price int
	var err error
	if priceStr == "" {
		price = 0
	} else {
		price, err = strconv.Atoi(priceStr)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid price parameter"})
			return
		}
	}

	req := cp.ProductGetAllReq{
		Name:   c.Query("name"),
		Price: int32(price),
		Filter:   &cp.Filter{},
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

	res, err := h.srvs.Product.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get products", "details": err.Error()})
		return
	}
	fmt.Println(res)
	c.JSON(http.StatusOK, res)
}

// ProductUpdate handles updating an existing product.
// @Summary Update product
// @Description Update an existing product
// @Tags product
// @Accept json
// @Produce json
// @Param id query string false "Id"
// @Param product body cp.ProductCreateReq true "Product data"
// @Success 200 {object} string "Product updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Product not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /product/{id} [put]
func (h *Handler) ProductUpdate(c *gin.Context) {
	req := cp.ProductCreateReq{}
	
	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	product := cp.ProductUpdateReq{
		Id:      c.Query("id"),
		Product: &req,
	}

	_, err := h.srvs.Product.Update(context.Background(), &product)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update product", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product updated"})
}

// ProductDelete handles deleting a product by ID.
// @Summary Delete product
// @Description Delete a product by ID
// @Tags product
// @Accept json
// @Produce json
// @Param id path string true "Product ID"
// @Success 200 {object} string "Product deleted"
// @Failure 400 {object} string "Invalid product ID"
// @Failure 404 {object} string "Product not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /product/{id} [delete]
func (h *Handler) ProductDelete(c *gin.Context) {
	id := &cp.ById{Id: c.Param("id")}
	_, err := h.srvs.Product.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete product", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Product deleted"})
}
