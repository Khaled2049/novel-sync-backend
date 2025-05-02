// handlers/novel_handler.go
package handlers

import (
	"context"
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/service"
)

type NovelHandler struct {
	novelService *service.NovelService
}

func NewNovelHandler(novelService *service.NovelService) *NovelHandler {
	return &NovelHandler{
		novelService: novelService,
	}
}

func (h *NovelHandler) RegisterRoutes(router *gin.Engine) {
	novelGroup := router.Group("/novels")
	{
		novelGroup.GET("", h.GetAllNovels)
		novelGroup.POST("", h.CreateNovel)
		// novelGroup.GET("/:id", h.GetNovelByID)
		// novelGroup.PUT("/:id", h.UpdateNovel)
		// novelGroup.DELETE("/:id", h.DeleteNovel)
	}
}

func (h *NovelHandler) GetAllNovels(c *gin.Context) {
	ctx := context.Background()
	novels, err := h.novelService.GetAllNovels(ctx)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to fetch novels"})
		return
	}

	c.JSON(200, novels)
}

func (h *NovelHandler) CreateNovel(c *gin.Context) {
	fmt.Println("Creating novel...")
	ctx := context.Background()
	var novel domain.Novel
	if err := c.ShouldBindJSON(&novel); err != nil {
		c.JSON(400, gin.H{"error": "Invalid input"})
		return
	}

	createdNovel, err := h.novelService.CreateNovel(ctx, &novel)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to create novel"})
		return
	}

	c.JSON(201, createdNovel)
}