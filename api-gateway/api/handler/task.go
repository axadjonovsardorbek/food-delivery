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

// TaskCreate handles the creation of a new task.
// @Summary Create task
// @Description Create a new task
// @Tags task
// @Accept json
// @Produce json
// @Param media body cp.TaskCreateReq true "Task data"
// @Success 200 {object} string "Task created"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /admin/task [post]
func (h *Handler) TaskCreate(c *gin.Context) {
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

	var req cp.TaskCreateReq

	if err := c.BindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
		return
	}

	_, err := h.srvs.Task.Create(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	_, err = h.srvs.Notification.Create(context.Background(), &cp.NotificationCreateReq{
		UserId:  req.AssignedTo,
		Message: "Your cart has been emptied",
	})

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		log.Println("error: ", err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "Task created"})
}

// TaskGetById handles the get a task.
// @Summary Get task
// @Description Get a task
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} cp.TaskGetByIdRes
// @Failure 400 {object} string "Invalid request payload"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /admin/task/{id} [get]
func (h *Handler) TaskGetById(c *gin.Context) {
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

	id := &cp.ById{Id: c.Param("id")}
	res, err := h.srvs.Task.GetById(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get task", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, res)
}

// TaskGetAll handles getting all tasks.
// @Summary Get all tasks
// @Description Get all tasks
// @Tags task
// @Accept json
// @Produce json
// @Param status query string false "Status"
// @Param page query integer false "Page"
// @Success 200 {object} cp.TaskGetAllRes
// @Failure 400 {object} string "Invalid parameters"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /admin/task/all [get]
func (h *Handler) TaskGetAll(c *gin.Context) {
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

	req := cp.TaskGetAllReq{
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

	res, err := h.srvs.Task.GetAll(context.Background(), &req)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get tasks", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, res)
}

// TaskUpdate handles updating an existing task.
// @Summary Update task
// @Description Update an existing task
// @Tags task
// @Accept json
// @Produce json
// @Param id query string false "Id"
// @Param status query string false "Status"
// @Success 200 {object} string "Task updated"
// @Failure 400 {object} string "Invalid request payload"
// @Failure 404 {object} string "Task not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /admin/task/{id} [put]
func (h *Handler) TaskUpdate(c *gin.Context) {
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

	task := cp.TaskUpdateReq{
		Id:     c.Query("id"),
		Status: c.Query("status"),
	}

	_, err := h.srvs.Task.Update(context.Background(), &task)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update task", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task updated"})
}

// TaskDelete handles deleting a task by ID.
// @Summary Delete task
// @Description Delete a task by ID
// @Tags task
// @Accept json
// @Produce json
// @Param id path string true "Task ID"
// @Success 200 {object} string "Task deleted"
// @Failure 400 {object} string "Invalid media ID"
// @Failure 404 {object} string "Task not found"
// @Failure 500 {object} string "Server error"
// @Security BearerAuth
// @Router /admin/task/{id} [delete]
func (h *Handler) TaskDelete(c *gin.Context) {
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

	id := &cp.ById{Id: c.Param("id")}
	_, err := h.srvs.Task.Delete(context.Background(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete task", "details": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Task deleted"})
}
