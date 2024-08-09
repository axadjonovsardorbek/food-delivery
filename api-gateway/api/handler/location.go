package handler

import (
	"context"
	cp "gateway/genproto/courier"
	"log"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/golang-jwt/jwt"
)

// LocationCreate handles the creation of a new location.
// @Summary Create location
// @Description Create a new location
// @Tags location
// @Accept json
// @Produce json
// @Param location body cp.LocationCreateReq true "Location data"
// @Success 200 {object} string "Location created"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/location [post]
func (h *Handler) LocationCreate(c *gin.Context) {
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

	var req cp.LocationCreateReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	req.CourierId = id

	_, err := h.srvs.Location.Create(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Location created"})
}

// LocationGetById handles the get a location.
// @Summary Get location
// @Description Get a location
// @Tags location
// @Accept json
// @Produce json
// @Param id path string true "Location ID"
// @Success 200 {object} cp.LocationGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/location/{id} [get]
func (h *Handler) LocationGetById(c *gin.Context) {
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

	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.Location.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get location", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// LocationGetAll handles getting all location.
// @Summary Get all location
// @Description Get all location
// @Tags location
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Success 200 {object} cp.LocationGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/location/all [get]
func (h *Handler) LocationGetAll(c *gin.Context) {
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

	req := cp.LocationGetAllReq{
		CourierId: id,
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

	res, err := h.srvs.Location.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get locations", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// LocationUpdate handles updating an existing location.
// @Summary Update location
// @Description Update an existing location
// @Tags location
// @Accept json
// @Produce json
// @Param id query string false "Id"
// @Param location query string false "Location"
// @Success 200 {object} string "Location updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Location not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/location/{id} [put]
func (h *Handler) LocationUpdate(c *gin.Context) {
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

	location := cp.LocationUpdateReq{
		Id:       c.Query("id"),
		Location: c.Query("location"),
	}

	_, err := h.srvs.Location.Update(context.Background(), &location)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update location", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Location updated"})
}

// LocationDelete handles deleting a location by ID.
// @Summary Delete location
// @Description Delete a location by ID
// @Tags location
// @Accept json
// @Produce json
// @Param id path string true "Location ID"
// @Success 200 {object} string "Location deleted"
// @Failure 400 {object} string "Invalid media ID"
// @Failure 404 {object} string "Location not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /courier/location/{id} [delete]
func (h *Handler) LocationDelete(c *gin.Context) {
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

	id := &cp.ById{Id: c.Param("id")}
	_, err := h.srvs.Location.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete location", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Location deleted"})
}
