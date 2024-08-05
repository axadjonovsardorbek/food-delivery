package handler

// import (
// 	"context"
// 	cp "gateway/genproto"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// 	"github.com/golang-jwt/jwt"
// )

// // EventCreate handles the creation of a new event.
// // @Summary Create event
// // @Description Create a new event
// // @Tags historical
// // @Accept json
// // @Produce json
// // @Param event body cp.HistoricalEventsCreateReq true "Event data"
// // @Success 200 {object} string "Event created"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/historical [post]
// func (h *Handler) HistoricalEventCreate(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	id := claims.(jwt.MapClaims)["user_id"].(string)

// 	var req cp.HistoricalEventsCreateReq
// 	var event cp.HistoricalEventsRes

// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	event.Category = req.Category
// 	event.Date = req.Date
// 	event.Description = req.Description
// 	event.UserId = id
// 	event.Title = req.Title

// 	_, err := h.srvs.HistoricalEvent.Create(context.Background(), &event)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		log.Println("error: ", err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Event created"})
// }

// // EventGetById handles the get a event.
// // @Summary Get event
// // @Description Get a event
// // @Tags historical
// // @Accept json
// // @Produce json
// // @Param id path string true "Event ID"
// // @Success 200 {object} cp.HistoricalEventsGetByIdRes
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/historical/{id} [get]
// func (h *Handler) HistoricalEventGetById(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	res, err := h.srvs.HistoricalEvent.GetById(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get event", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// // EventGetAll handles getting all event.
// // @Summary Get all event
// // @Description Get all event
// // @Tags historical
// // @Accept json
// // @Produce json
// // @Param date query string false "Date"
// // @Param category query string false "Category"
// // @Param page query integer false "Page"
// // @Success 200 {object} cp.HistoricalEventsGetAllRes
// // @Failure 400 {object} string "Invalid parameters"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/historical/all [get]
// func (h *Handler) HistoricalEventGetAll(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	userId := claims.(jwt.MapClaims)["user_id"].(string)
// 	req := cp.HistoricalEventsGetAllReq{
// 		UserId:   userId,
// 		Date:     c.Query("date"),
// 		Category: c.Query("category"),
// 		Filter:   &cp.Filter{},
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

// 	res, err := h.srvs.HistoricalEvent.GetAll(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get events", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, res)
// }

// // EventUpdate handles updating an existing event.
// // @Summary Update event
// // @Description Update an existing event
// // @Tags historical
// // @Accept json
// // @Produce json
// // @Param id query string false "Id"
// // @Param title query string false "Title"
// // @Param description query string false "Description"
// // @Success 200 {object} string "Event updated"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 404 {object} string "Event not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/historical/{id} [put]
// func (h *Handler) HistoricalEventUpdate(c *gin.Context) {
// 	media := cp.HistoricalEventsUpdateReq{
// 		Id:          c.Query("id"),
// 		Title:       c.Query("title"),
// 		Description: c.Query("description"),
// 	}

// 	_, err := h.srvs.HistoricalEvent.Update(context.Background(), &media)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update event", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Event updated"})
// }

// // EventDelete handles deleting a event by ID.
// // @Summary Delete event
// // @Description Delete a event by ID
// // @Tags historical
// // @Accept json
// // @Produce json
// // @Param id path string true "Event ID"
// // @Success 200 {object} string "Event deleted"
// // @Failure 400 {object} string "Invalid event ID"
// // @Failure 404 {object} string "Event not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/historical/{id} [delete]
// func (h *Handler) HistoricalEventDelete(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	_, err := h.srvs.HistoricalEvent.Delete(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete event", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Event deleted"})
// }

// // EventGetAll handles getting all event.
// // @Summary Get all event
// // @Description Get all event
// // @Tags context
// // @Accept json
// // @Produce json
// // @Param date query string false "Date"
// // @Success 200 {object} cp.ContextRes
// // @Failure 400 {object} string "Invalid parameters"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /timeline/context/{date} [get]
// func (h *Handler) Context(c *gin.Context) {
// 	claims, exists := c.Get("claims")
// 	if !exists {
// 		c.JSON(http.StatusUnauthorized, gin.H{"error": "Unauthorized"})
// 		return
// 	}

// 	userId := claims.(jwt.MapClaims)["user_id"].(string)
// 	req := cp.ContextReq{
// 		UserId: userId,
// 		Date:   c.Query("date"),
// 	}

// 	res, err := h.srvs.HistoricalEvent.Context(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get events", "details": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, res)
// }
