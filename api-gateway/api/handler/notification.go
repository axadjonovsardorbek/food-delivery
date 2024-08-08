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

// NotificationCreate handles the creation of a new notification.
// @Summary Create notification
// @Description Create a new notification
// @Tags notification
// @Accept json
// @Produce json
// @Param memory body cp.NotificationCreateReq true "Notification data"
// @Success 200 {object} string "Notification created"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /notification [post]
func (h *Handler) NotificationCreate(c *gin.Context) {
	var req cp.NotificationCreateReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	_, err := h.srvs.Notification.Create(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Notification created"})
}

// NotificationGetById handles the get a notification.
// @Summary Get notification
// @Description Get a notification
// @Tags notification
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} cp.NotificationGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /notification/{id} [get]
func (h *Handler) NotificationGetById(c *gin.Context) {
	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.Notification.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get notification", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// NotificationGetAll handles getting all notification.
// @Summary Get all notification
// @Description Get all notification
// @Tags notification
// @Accept json
// @Produce json
// @Param page query integer false "Page"
// @Success 200 {object} cp.NotificationGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /notification/all [get]
func (h *Handler) NotificationGetAll(c *gin.Context) {
	claims, exists := c.Get("claims")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
		return
	}

	userId := claims.(jwt.MapClaims)["user_id"].(string)

	req := cp.NotificationGetAllReq{
		UserId: userId,
		Filter: &cp.Filter{},
	}

	pageStr := c.Query("page")
	var page int
	var err error
	if pageStr == "" {
		page = 0
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

	res, err := h.srvs.Notification.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get notifications", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// NotificationUpdate handles updating an existing notification.
// @Summary Update notification
// @Description Update an existing notification
// @Tags notification
// @Accept json
// @Produce json
// @Param id query string false "Id"
// @Param is_read query string false "IsRead"
// @Success 200 {object} string "Notification updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Notification not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /notification/{id} [put]
func (h *Handler) NotificationUpdate(c *gin.Context) {
	memory := cp.NotificationUpdateReq{
		Id:     c.Query("id"),
		IsRead: c.Query("is_read"),
	}

	_, err := h.srvs.Notification.Update(context.Background(), &memory)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update notification", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification updated"})
}

// NotificationDelete handles deleting a notification by ID.
// @Summary Delete notification
// @Description Delete a notification by ID
// @Tags notification
// @Accept json
// @Produce json
// @Param id path string true "Notification ID"
// @Success 200 {object} string "Notification deleted"
// @Failure 400 {object} string "Invalid notification ID"
// @Failure 404 {object} string "Notification not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /notification/{id} [delete]
func (h *Handler) NotificationDelete(c *gin.Context) {
	id := &cp.ById{Id: c.Param("id")}
	_, err := h.srvs.Notification.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete notification", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification deleted"})
}
