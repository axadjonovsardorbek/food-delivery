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

// // SharedMemoryCreate handles the creation of a new memory.
// // @Summary Shared memory
// // @Description Shared a new memory
// // @Tags shared
// // @Accept json
// // @Produce json
// // @Param memory body cp.SharedMemoriesCreateReq true "SharedMemory data"
// // @Success 200 {object} string "SharedMemory created"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/shared [post]
// func (h *Handler) SharedMemoryCreate(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	id := claims.(jwt.MapClaims)["user_id"].(string)

// 	var req cp.SharedMemoriesCreateReq

// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	req.SharedId = id

// 	_, err := h.srvs.SharedMemory.Create(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		log.Println("error: ", err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "SharedMemory created"})
// }

// // SharedMemoryGetById handles the get a memory.
// // @Summary Get memory
// // @Description Get a memory
// // @Tags shared
// // @Accept json
// // @Produce json
// // @Param id path string true "SharedMemory ID"
// // @Success 200 {object} cp.SharedMemoriesGetByIdRes
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/shared/{id} [get]
// func (h *Handler) SharedMemoryGetById(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	res, err := h.srvs.SharedMemory.GetById(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get shared memory", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// // SharedMemoryGetAll handles getting all memories.
// // @Summary Get all memories
// // @Description Get all memories
// // @Tags shared
// // @Accept json
// // @Produce json
// // @Param shared_id query string false "SharedId"
// // @Param recipient_id query string false "RecipientId"
// // @Param page query integer false "Page"
// // @Success 200 {object} cp.SharedMemoriesGetAllRes
// // @Failure 400 {object} string "Invalid parameters"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/shared/all [get]
// func (h *Handler) SharedMemoryGetAll(c *gin.Context) {
// 	req := cp.SharedMemoriesGetAllReq{
// 		SharedId:    c.Query("shared_id"),
// 		RecipientId: c.Query("recipient_id"),
// 		Filter:      &cp.Filter{},
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

// 	res, err := h.srvs.SharedMemory.GetAll(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get shared memories", "details": err.Error()})
// 		return
// 	}
// 	fmt.Println(res)
// 	c.JSON(http.StatusOK, res)
// }

// // SharedMemoryUpdate handles updating an existing memory.
// // @Summary Update shared memory
// // @Description Update an existing shared memory
// // @Tags shared
// // @Accept json
// // @Produce json
// // @Param id query string false "Id"
// // @Param message query string false "Message"
// // @Success 200 {object} string "Memory updated"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 404 {object} string "Memory not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/shared/{id} [put]
// func (h *Handler) SharedMemoryUpdate(c *gin.Context) {
// 	memory := cp.SharedMemoriesUpdateReq{
// 		Id:      c.Query("id"),
// 		Message: c.Query("message"),
// 	}

// 	_, err := h.srvs.SharedMemory.Update(context.Background(), &memory)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update memory", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Memory updated"})
// }

// // SharedMemoryDelete handles deleting a memory by ID.
// // @Summary Delete memory
// // @Description Delete a memory by ID
// // @Tags shared
// // @Accept json
// // @Produce json
// // @Param id path string true "Memory ID"
// // @Success 200 {object} string "Memory deleted"
// // @Failure 400 {object} string "Invalid memory ID"
// // @Failure 404 {object} string "Memory not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/shared/{id} [delete]
// func (h *Handler) SharedMemoryDelete(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	_, err := h.srvs.SharedMemory.Delete(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete memory", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Memory deleted"})
// }
