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

// // MilestoneCreate handles the creation of a new milestone.
// // @Summary Create milestone
// // @Description Create a new milestone
// // @Tags milestone
// // @Accept json
// // @Produce json
// // @Param memory body cp.MilestonesCreateReq true "Milestone data"
// // @Success 200 {object} string "Milestone created"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/milestone [post]
// func (h *Handler) MilestoneCreate(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	id := claims.(jwt.MapClaims)["user_id"].(string)
// 	var req cp.MilestonesCreateReq

// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	req.UserId = id

// 	_, err := h.srvs.Milestone.Create(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		log.Println("error: ", err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Milestone created"})
// }

// // MilestoneGetById handles the get a milestone.
// // @Summary Get milestone
// // @Description Get a milestone
// // @Tags milestone
// // @Accept json
// // @Produce json
// // @Param id path string true "Milestone ID"
// // @Success 200 {object} cp.MilestonesGetByIdRes
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/milestone/{id} [get]
// func (h *Handler) MilestoneGetById(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	res, err := h.srvs.Milestone.GetById(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get milestone", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// // MilestoneGetAll handles getting all milestones.
// // @Summary Get all milestones
// // @Description Get all milestones
// // @Tags milestone
// // @Accept json
// // @Produce json
// // @Param page query integer false "Page"
// // @Success 200 {object} cp.MilestonesGetAllRes
// // @Failure 400 {object} string "Invalid parameters"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/milestone/all [get]
// func (h *Handler) MilestoneGetAll(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	userId := claims.(jwt.MapClaims)["user_id"].(string)
// 	req := cp.MilestonesGetAllReq{
// 		UserId: userId,
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

// 	res, err := h.srvs.Milestone.GetAll(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get milestones", "details": err.Error()})
// 		return
// 	}
// 	fmt.Println(res)
// 	c.JSON(http.StatusOK, res)
// }

// // MilestoneUpdate handles updating an existing milestone.
// // @Summary Update milestone
// // @Description Update an existing milestone
// // @Tags milestone
// // @Accept json
// // @Produce json
// // @Param id query string false "Id"
// // @Param title query string false "Title"
// // @Param date query string false "Date"
// // @Success 200 {object} string "Milestone updated"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 404 {object} string "Milestone not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/milestone/{id} [put]
// func (h *Handler) MilestoneUpdate(c *gin.Context) {
// 	memory := cp.MilestonesUpdateReq{
// 		Id:    c.Query("id"),
// 		Title: c.Query("title"),
// 		Date:  c.Query("date"),
// 	}

// 	_, err := h.srvs.Milestone.Update(context.Background(), &memory)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update milestone", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Milestone updated"})
// }

// // MilestoneDelete handles deleting a milestone by ID.
// // @Summary Delete milestone
// // @Description Delete a milestone by ID
// // @Tags milestone
// // @Accept json
// // @Produce json
// // @Param id path string true "Milestone ID"
// // @Success 200 {object} string "Milestone deleted"
// // @Failure 400 {object} string "Invalid milestone ID"
// // @Failure 404 {object} string "Milestone not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/milestone/{id} [delete]
// func (h *Handler) MilestoneDelete(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	_, err := h.srvs.Milestone.Delete(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete milestone", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Milestone deleted"})
// }
