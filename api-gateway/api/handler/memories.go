package handler

// import (
// 	"context"
// 	"fmt"
// 	cp "gateway/genproto"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// )

// // MemoryCreate handles the creation of a new memory.
// // @Summary Create memory
// // @Description Create a new memory
// // @Tags memory
// // @Accept json
// // @Produce json
// // @Param memory body cp.MemoriesCreateReq true "Memory data"
// // @Success 200 {object} string "Memory created"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory [post]
// func (h *Handler) MemoryCreate(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	id := claims.(jwt.MapClaims)["user_id"].(string)

// 	var req cp.MemoriesCreateReq

// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	req.UserId = id

// 	_, err := h.srvs.Memory.Create(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		log.Println("error: ", err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Memory created"})
// }

// // MemoryGetById handles the get a memory.
// // @Summary Get memory
// // @Description Get a memory
// // @Tags memory
// // @Accept json
// // @Produce json
// // @Param id path string true "Memory ID"
// // @Success 200 {object} cp.MemoriesGetByIdRes
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id} [get]
// func (h *Handler) MemoryGetById(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	res, err := h.srvs.Memory.GetById(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get memory", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// // MemoryGetAll handles getting all memories.
// // @Summary Get all memories
// // @Description Get all memories
// // @Tags memory
// // @Accept json
// // @Produce json
// // @Param page query integer false "Page"
// // @Success 200 {object} cp.MemoriesGetAllRes
// // @Failure 400 {object} string "Invalid parameters"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/all [get]
// func (h *Handler) MemoryGetAll(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	id := claims.(jwt.MapClaims)["user_id"].(string)

// 	req := cp.MemoriesGetAllReq{
// 		UserId: id,
// 		Filter: &cp.Filter{},
// 	}

// 	pageStr := c.Query("page")
// 	var page int
// 	var err error
// 	if pageStr == "" {
// 		page = 0
// 	} else {
// 		page, err = strconv.Atoi(pageStr)
// 		if err != nil {
// 			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid page parameter"})
// 			return
// 		}
// 	}

// 	filter := cp.Filter{
// 		Page: int32(page),
// 	}

// 	req.Filter.Page = filter.Page

// 	res, err := h.srvs.Memory.GetAll(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get memories", "details": err.Error()})
// 		return
// 	}
// 	fmt.Println(res)
// 	c.JSON(http.StatusOK, res)
// }

// // MemoryUpdate handles updating an existing memory.
// // @Summary Update memory
// // @Description Update an existing memory
// // @Tags memory
// // @Accept json
// // @Produce json
// // @Param id query string false "Id"
// // @Param title query string false "Title"
// // @Param description query string false "Description"
// // @Param privacy query string false "Privacy"
// // @Success 200 {object} string "Memory updated"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 404 {object} string "Memory not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id} [put]
// func (h *Handler) MemoryUpdate(c *gin.Context) {
// 	memory := cp.MemoriesUpdateReq{
// 		Id:          c.Query("id"),
// 		Title:       c.Query("title"),
// 		Description: c.Query("description"),
// 		Privacy:     c.Query("privacy"),
// 	}

// 	_, err := h.srvs.Memory.Update(context.Background(), &memory)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update memory", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Memory updated"})
// }

// // MemoryDelete handles deleting a memory by ID.
// // @Summary Delete memory
// // @Description Delete a memory by ID
// // @Tags memory
// // @Accept json
// // @Produce json
// // @Param id path string true "Memory ID"
// // @Success 200 {object} string "Memory deleted"
// // @Failure 400 {object} string "Invalid memory ID"
// // @Failure 404 {object} string "Memory not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id} [delete]
// func (h *Handler) MemoryDelete(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	_, err := h.srvs.Memory.Delete(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete memory", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Memory deleted"})
// }
