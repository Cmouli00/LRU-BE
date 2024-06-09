package controller

import (
	"fmt"
	"net/http"

	"lru/resources"

	"lru/service"

	"github.com/gin-gonic/gin"
)

// CacheHandler represents the handler for cache operations
type CacheHandler struct {
	cache *service.LRUCache
}

// NewCacheHandler creates a new instance of CacheHandler
func NewCacheHandler(cache *service.LRUCache) *CacheHandler {
	return &CacheHandler{
		cache: cache,
	}
}

func (h *CacheHandler) Get(c *gin.Context) {

	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, resources.ErrorResponse{Message: "Key is required"})
		return
	}

	value, found := h.cache.Get(key)
	if !found {
		c.JSON(http.StatusOK, resources.GetResponse{Found: found})
		return
	}

	c.JSON(http.StatusOK, resources.GetResponse{Value: fmt.Sprint(value), Found: found})

}

func (h *CacheHandler) GetAll(c *gin.Context) {

	value := h.cache.GetAll()

	c.JSON(http.StatusOK, value)

}

func (h *CacheHandler) Set(c *gin.Context) {

	var req resources.SetRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusBadRequest, resources.ErrorResponse{Message: "Invalid request payload"})
		return
	}
	if req.Expiration < 0 {
		c.JSON(http.StatusBadRequest, resources.ErrorResponse{Message: "Invalid expiration payload"})
		return
	}

	h.cache.Set(req.Key, req.Value, req.Expiration)

	c.JSON(http.StatusOK, "success")

}

func (h *CacheHandler) Delete(c *gin.Context) {

	key := c.Param("key")
	if key == "" {
		c.JSON(http.StatusBadRequest, resources.ErrorResponse{Message: "Key is required"})
		return
	}

	h.cache.Delete(key)

	c.JSON(http.StatusOK, "success")

}
