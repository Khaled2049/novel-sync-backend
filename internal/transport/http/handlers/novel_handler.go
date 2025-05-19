// handlers/novel_handler.go
package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
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

		// Routes for chapters specifically related to a novel
		novelChaptersGroup := novelGroup.Group("/:novelID/chapters")
		{
			novelChaptersGroup.POST("", h.AddChapterToNovelHandler)
			// Add other novel-specific chapter routes here if needed, e.g., ListChaptersByNovelID
		}

		// novelGroup.PUT("/:id", h.UpdateNovel) // Service method not implemented
		// novelGroup.DELETE("/:id", h.DeleteNovel) // Service method not implemented
	}

	// Group for operations on chapters directly if chapter IDs are globally unique
	// If chapter IDs are only unique within a novel, these routes might need to be /novels/:novelID/chapters/:chapterID/...
	chapterGroup := router.Group("/chapters")
	{
		chapterGroup.PUT("/:chapterID/autosave", h.AutosaveChapterHandler)
		chapterGroup.POST("/:chapterID/save-revision", h.SaveChapterWithRevisionHandler)
		// Add other chapter-specific routes here, e.g., GetChapterByID, UpdateChapter
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

	// TODO: Extract OwnerUserID from authenticated user context
	// For now, assuming it's set in the request or handled by service default
	// if novel.OwnerUserID == "" {
	//     userID, exists := c.Get("userID") // Example: Get userID from auth middleware
	//     if exists {
	//         novel.OwnerUserID = userID.(string)
	//     } else {
	//         c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID not found"})
	//         return
	//     }
	// }

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

	novel, err := h.novelService.GetNovelByID(ctx, novelID)
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

	// TODO: Validate OwnerUserID in req.NovelData or set it from authenticated user
	// if req.NovelData.OwnerUserID == "" {
	//     userID, exists := c.Get("userID") // Example: Get userID from auth middleware
	//     if exists {
	//         req.NovelData.OwnerUserID = userID.(string)
	//     } else {
	//         c.JSON(http.StatusUnauthorized, gin.H{"error": "Owner User ID is required"})
	//         return
	//     }
	// }

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

// AddChapterToNovelRequest defines the payload for adding a chapter.
// The service expects a full domain.Chapter object.

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
		NovelID:         novelID, // Will be overridden by service, but good to have
		Title:           reqChapter.Title,
		Content:         reqChapter.Content,
		Status:          domain.ChapterStatusDraft, // Default status
		LastEditedByUserID: lastEditedByUserID,
		// OrderIndex and WordCount will be handled by the service/repository
	}


	createdChapter, err := h.novelService.AddChapterToNovel(ctx, novelID, chapter)
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


// AutosaveChapterHandler handles autosaving chapter content.
func (h *NovelHandler) AutosaveChapterHandler(c *gin.Context) {
	ctx := c.Request.Context()
	chapterID := c.Param("chapterID")

	var req request.AutosaveChapterRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for autosave", "details": err.Error()})
		return
	}
    
	// TODO: Extract UserID from authenticated user context
	// userIDAuth, exists := c.Get("userID")
	// if !exists || userIDAuth.(string) != req.UserID { // Optional: verify req.UserID matches auth user
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID mismatch or not authenticated"})
	//     return
	// }

	err := h.novelService.AutosaveChapter(ctx, chapterID, req.Content, req.UserID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to autosave chapter", "details": err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}

// SaveChapterWithRevisionRequest defines the payload for saving a chapter with a revision.
type SaveChapterWithRevisionRequest struct {
	NewContent string `json:"new_content"`
	UserID     string `json:"user_id" binding:"required"` // Or get from auth context
	Notes      string `json:"notes"`
}

// SaveChapterWithRevisionHandler handles saving chapter content and creating a revision.
func (h *NovelHandler) SaveChapterWithRevisionHandler(c *gin.Context) {
	ctx := c.Request.Context()
	chapterID := c.Param("chapterID")

	var req SaveChapterWithRevisionRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input for saving revision", "details": err.Error()})
		return
	}

	// TODO: Extract UserID from authenticated user context
	// userIDAuth, exists := c.Get("userID")
	// if !exists || userIDAuth.(string) != req.UserID { // Optional: verify req.UserID matches auth user
	//     c.JSON(http.StatusUnauthorized, gin.H{"error": "User ID mismatch or not authenticated"})
	//     return
	// }

	err := h.novelService.SaveChapterWithRevision(ctx, chapterID, req.NewContent, req.UserID, req.Notes)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save chapter with revision", "details": err.Error()})
		return
	}

	c.Status(http.StatusOK) // Or http.StatusNoContent if no body is returned by design
}