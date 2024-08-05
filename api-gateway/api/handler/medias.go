package handler

// import (
// 	"context"
// 	"fmt"
// 	cp "gateway/genproto"
// 	"log"
// 	"net/http"
// 	"strconv"

// 	"github.com/gin-gonic/gin"
// )

// // MediaCreate handles the creation of a new media.
// // @Summary Create media
// // @Description Create a new media
// // @Tags media
// // @Accept json
// // @Produce json
// // @Param media body cp.MediasCreateReq true "Media data"
// // @Success 200 {object} string "Media created"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/media [post]
// func (h *Handler) MediaCreate(c *gin.Context) {
// 	var req cp.MediasCreateReq

// 	if err := c.BindJSON(&req); err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request payload"})
// 		return
// 	}

// 	_, err := h.srvs.Media.Create(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
// 		log.Println("error: ", err)
// 		return
// 	}

// 	c.JSON(http.StatusOK, gin.H{"message": "Media created"})
// }

// // MediaGetById handles the get a media.
// // @Summary Get media
// // @Description Get a media
// // @Tags media
// // @Accept json
// // @Produce json
// // @Param id path string true "Media ID"
// // @Success 200 {object} cp.MediasGetByIdRes
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/media/{id} [get]
// func (h *Handler) MediaGetById(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	res, err := h.srvs.Media.GetById(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get media", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, res)
// }

// // MediaGetAll handles getting all medias.
// // @Summary Get all medias
// // @Description Get all medias
// // @Tags media
// // @Accept json
// // @Produce json
// // @Param memory_id query string false "MemoryId"
// // @Param page query integer false "Page"
// // @Success 200 {object} cp.MediasGetAllRes
// // @Failure 400 {object} string "Invalid parameters"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/media/all [get]
// func (h *Handler) MediaGetAll(c *gin.Context) {
// 	req := cp.MediasGetAllReq{
// 		MemoryId: c.Query("memory_id"),
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

// 	res, err := h.srvs.Media.GetAll(context.Background(), &req)

// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't get medias", "details": err.Error()})
// 		return
// 	}
// 	fmt.Println(res)
// 	c.JSON(http.StatusOK, res)
// }

// // MediaUpdate handles updating an existing media.
// // @Summary Update media
// // @Description Update an existing media
// // @Tags media
// // @Accept json
// // @Produce json
// // @Param id query string false "Id"
// // @Param type query string false "Type"
// // @Param url query string false "Url"
// // @Success 200 {object} string "Media updated"
// // @Failure 400 {object} string "Invalid request payload"
// // @Failure 404 {object} string "Media not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/media/{id} [put]
// func (h *Handler) MediaUpdate(c *gin.Context) {
// 	media := cp.MediasUpdateReq{
// 		Id:          c.Query("id"),
// 		Type:       c.Query("type"),
// 		Url: c.Query("url"),
// 	}

// 	_, err := h.srvs.Media.Update(context.Background(), &media)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't update media", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Media updated"})
// }

// // MediaDelete handles deleting a media by ID.
// // @Summary Delete media
// // @Description Delete a media by ID
// // @Tags media
// // @Accept json
// // @Produce json
// // @Param id path string true "Media ID"
// // @Success 200 {object} string "Media deleted"
// // @Failure 400 {object} string "Invalid media ID"
// // @Failure 404 {object} string "Media not found"
// // @Failure 500 {object} string "Server error"
// // @Security BearerAuth
// // @Router /memory/{id}/media/{id} [delete]
// func (h *Handler) MediaDelete(c *gin.Context) {
// 	id := &cp.ById{Id: c.Param("id")}
// 	_, err := h.srvs.Media.Delete(context.Background(), id)
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Couldn't delete media", "details": err.Error()})
// 		return
// 	}
// 	c.JSON(http.StatusOK, gin.H{"message": "Media deleted"})
// }
