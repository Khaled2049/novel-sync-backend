// handlers/novel_handler.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/service"
	"github.com/khaled2049/server/internal/transport/http/request"
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
		novelGroup.POST("", h.CreateNovelHandler)
		novelGroup.GET("", h.GetAllNovelsHandler)
		novelGroup.GET("/:novelID", h.GetNovelByIDHandler)
		novelGroup.POST("/with-first-chapter", h.CreateNovelWithFirstChapterHandler)
		novelGroup.POST("/:novelID/characters", h.CreateCharacterForNovelHandler)

		// Routes for chapters specifically related to a novel
		novelChaptersGroup := novelGroup.Group("/:novelID/chapters")
		{
			novelChaptersGroup.POST("", h.AddChapterToNovelHandler)

		}

	}

}

// GetAllNovelsHandler handles fetching all novels.
func (h *NovelHandler) GetAllNovelsHandler(c *gin.Context) {
	ctx := c.Request.Context() // Use request context
	novels, err := h.novelService.GetAllNovels(ctx)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch novels", "details": err.Error()})
		return
	}

	c.JSON(http.StatusOK, novels)
}

// CreateNovelHandler handles the creation of a new novel.
func (h *NovelHandler) CreateNovelHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var novel domain.Novel
	if err := c.ShouldBindJSON(&novel); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	createdNovel, err := h.novelService.CreateNovel(ctx, &novel)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create novel", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, createdNovel)
}

// GetNovelByIDHandler handles fetching a single novel by its ID.
func (h *NovelHandler) GetNovelByIDHandler(c *gin.Context) {
	ctx := c.Request.Context()
	novelID := c.Param("novelID")

	// convert novelID from string to uuid.UUID
	parsedNovelID, err := uuid.Parse(novelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid novel ID", "details": err.Error()})
		return
	}

	novel, err := h.novelService.GetNovelByID(ctx, parsedNovelID)
	if err != nil {
		// TODO: Differentiate between not found and other errors
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch novel", "details": err.Error()})
		return
	}
	if novel == nil { // Should be handled by error in a real repo (e.g., sql.ErrNoRows)
		c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found"})
		return
	}

	c.JSON(http.StatusOK, novel)
}

// CreateNovelWithFirstChapterHandler handles creating a novel along with its first chapter.
func (h *NovelHandler) CreateNovelWithFirstChapterHandler(c *gin.Context) {
	ctx := c.Request.Context()
	var req request.CreateNovelWithFirstChapterRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input", "details": err.Error()})
		return
	}

	novel, chapter, err := h.novelService.CreateNovelWithFirstChapter(
		ctx,
		&req.NovelData,
		req.ChapterTitle,
		req.InitialContent,
		req.NovelData.OwnerUserID, // Service expects userID, which is novel's owner
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create novel with first chapter", "details": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"novel":   novel,
		"chapter": chapter,
	})
}

// AddChapterToNovelHandler handles adding a new chapter to an existing novel.
func (h *NovelHandler) AddChapterToNovelHandler(c *gin.Context) {
	ctx := c.Request.Context()
	novelID := c.Param("novelID")

	var reqChapter request.AddChapterToNovelRequest
	if err := c.ShouldBindJSON(&reqChapter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for chapter", "details": err.Error()})
		return
	}

	// TODO: Extract LastEditedByUserID from authenticated user context if not in request
	// userID, exists := c.Get("userID")
	// if !exists {
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
	//     return
	// }
	// lastEditedByUserID := userID.(string)
	lastEditedByUserID := reqChapter.LastEditedByUserID // Assuming it's passed in request for now

	fmt.Println("DEBUG: lastEditedByUserID = ", lastEditedByUserID)

	chapter := &domain.Chapter{
		NovelID:            novelID, // Will be overridden by service, but good to have
		Title:              reqChapter.Title,
		Content:            reqChapter.Content,
		Status:             domain.ChapterStatusDraft, // Default status
		LastEditedByUserID: lastEditedByUserID,
		// OrderIndex and WordCount will be handled by the service/repository
	}

	// convert novelID from string to uuid.UUID
	parsedNovelID, err := uuid.Parse(novelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid novel ID", "details": err.Error()})
		return
	}

	createdChapter, err := h.novelService.AddChapterToNovel(ctx, parsedNovelID, chapter)
	if err != nil {
		// Check if it's a "novel not found" type of error
		if err.Error() == fmt.Sprintf("novel not found: %s", novelID) { // This check is brittle; better to use custom errors
			c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found", "details": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add chapter to novel", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createdChapter)
}

func (h *NovelHandler) CreateCharacterForNovelHandler(c *gin.Context) {
	ctx := c.Request.Context()
	novelID := c.Param("novelID")

	var reqCharacter request.CreateCharacterRequest
	if err := c.ShouldBindJSON(&reqCharacter); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for character", "details": err.Error()})
		return
	}

	// convert novelID from string to uuid.UUID
	parsedNovelID, err := uuid.Parse(novelID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid novel ID", "details": err.Error()})
		return
	}

	character := &domain.Character{
		NovelID:             parsedNovelID,
		Name:                reqCharacter.Name,
		Description:         reqCharacter.Description,
		Backstory:           reqCharacter.Backstory,
		Motivations:         reqCharacter.Motivations,
		PhysicalDescription: reqCharacter.PhysicalDescription,
		ImageURL:            reqCharacter.ImageURL,
	}

	createdCharacter, err := h.novelService.CreateCharacter(ctx, parsedNovelID, character)
	if err != nil {
		// Check if it's a "novel not found" type of error
		if err.Error() == fmt.Sprintf("novel not found: %s", novelID) { // This check is brittle; better to use custom errors
			c.JSON(http.StatusNotFound, gin.H{"error": "Novel not found", "details": err.Error()})
		} else {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to add character to novel", "details": err.Error()})
		}
		return
	}

	c.JSON(http.StatusCreated, createdCharacter)
}
