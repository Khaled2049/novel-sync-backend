package postgres

import (
	"context"
	"errors"
	"fmt"
	"log"

	"github.com/google/uuid"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
	"github.com/khaled2049/server/internal/domain"
	"github.com/khaled2049/server/internal/repository"
)

// postgresChapterRepository implements the repository.ChapterRepository interface.
type postgresChapterRepository struct {
	pool *pgxpool.Pool
}

// NewChapterRepository creates a new instance of postgresChapterRepository.
func NewChapterRepository(pool *pgxpool.Pool) repository.ChapterRepository {
	return &postgresChapterRepository{pool: pool}
}

// Create saves a new chapter to the storage.
func (r *postgresChapterRepository) Create(ctx context.Context, chapter *domain.Chapter) (*domain.Chapter, error) {

	query := `
		INSERT INTO chapters (
			novel_id, title, content, status, order_index, word_count,
			last_edited_by_user_id
		) VALUES (
			$1, $2, $3, $4, $5, $6, $7
		)
		RETURNING id, created_at, updated_at;`

	// Calculate word count if not provided
	wordCount := chapter.WordCount
	if wordCount == 0 && chapter.Content != "" {
		// Simple word count calculation
		// In a real implementation, you might want a more sophisticated approach
		wordCount = countWords(chapter.Content)
	}

	// Default status if not provided
	status := chapter.Status
	if status == "" {
		status = "draft" // Using your enum default
	}

	// Generate a new ID if not provided
	if chapter.ID == "" {
		fmt.Println("Chapter ID is empty, creating a new chapter")
	}

	err := r.pool.QueryRow(ctx, query,
		chapter.NovelID, chapter.Title, chapter.Content,
		status, chapter.OrderIndex, wordCount, chapter.LastEditedByUserID,
	).Scan(&chapter.ID, &chapter.CreatedAt, &chapter.UpdatedAt)

	if err != nil {
		log.Printf("Error creating chapter: %v", err)
		return nil, fmt.Errorf("failed to create chapter: %w", err)
	}

	// Update the domain object with calculated values
	chapter.WordCount = wordCount
	chapter.Status = status

	return chapter, nil
}

// GetByID retrieves a chapter by its ID.
func (r *postgresChapterRepository) GetByID(ctx context.Context, id uuid.UUID) (*domain.Chapter, error) {
	query := `
		SELECT id, novel_id, title, content, status, order_index, word_count,
			last_edited_by_user_id, created_at, updated_at, published_at
		FROM chapters
		WHERE id = $1;`

	chapter := &domain.Chapter{}
	err := r.pool.QueryRow(ctx, query, id).Scan(
		&chapter.ID, &chapter.NovelID, &chapter.Title, &chapter.Content,
		&chapter.Status, &chapter.OrderIndex, &chapter.WordCount,
		&chapter.LastEditedByUserID, &chapter.CreatedAt, &chapter.UpdatedAt, &chapter.PublishedAt,
	)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return nil, errors.New("chapter not found")
		}
		log.Printf("Error scanning chapter by ID %s: %v", id, err)
		return nil, fmt.Errorf("failed to find chapter by ID: %w", err)
	}

	return chapter, nil
}

// Update updates an existing chapter in the storage.
func (r *postgresChapterRepository) Update(ctx context.Context, chapter *domain.Chapter) error {
	query := `
		UPDATE chapters
		SET title = $2,
			content = $3,
			status = $4,
			order_index = $5,
			word_count = $6,
			last_edited_by_user_id = $7,
			updated_at = NOW()
		WHERE id = $1
		RETURNING updated_at;`

	// Calculate word count if content is provided
	wordCount := chapter.WordCount
	if chapter.Content != "" {
		wordCount = countWords(chapter.Content)
	}

	err := r.pool.QueryRow(ctx, query,
		chapter.ID, chapter.Title, chapter.Content, chapter.Status,
		chapter.OrderIndex, wordCount, chapter.LastEditedByUserID,
	).Scan(&chapter.UpdatedAt)

	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return errors.New("chapter not found")
		}
		log.Printf("Error updating chapter with ID %s: %v", chapter.ID, err)
		return fmt.Errorf("failed to update chapter: %w", err)
	}

	// Update the word count in the domain object
	chapter.WordCount = wordCount

	return nil
}

// Delete removes a chapter from the storage by its ID.
func (r *postgresChapterRepository) Delete(ctx context.Context, id uuid.UUID) error {
	query := `
		DELETE FROM chapters
		WHERE id = $1;`

	result, err := r.pool.Exec(ctx, query, id)
	if err != nil {
		log.Printf("Error deleting chapter with ID %s: %v", id, err)
		return fmt.Errorf("failed to delete chapter: %w", err)
	}

	rowsAffected := result.RowsAffected()

	if rowsAffected == 0 {
		return errors.New("chapter not found")
	}

	return nil
}

// ListByNovelID retrieves all chapters belonging to a specific novel.
func (r *postgresChapterRepository) ListByNovelID(ctx context.Context, novelID uuid.UUID) ([]*domain.Chapter, error) {
	query := `
		SELECT id, novel_id, title, content, status, order_index, word_count,
			last_edited_by_user_id, created_at, updated_at, published_at
		FROM chapters
		WHERE novel_id = $1
		ORDER BY order_index ASC;`

	var chapters []*domain.Chapter
	rows, err := r.pool.Query(ctx, query, novelID)
	if err != nil {
		log.Printf("Error querying chapters by novel ID %s: %v", novelID, err)
		return nil, fmt.Errorf("failed to find chapters by novel ID: %w", err)
	}
	defer rows.Close()

	for rows.Next() {
		chapter := &domain.Chapter{}
		err := rows.Scan(
			&chapter.ID, &chapter.NovelID, &chapter.Title, &chapter.Content,
			&chapter.Status, &chapter.OrderIndex, &chapter.WordCount,
			&chapter.LastEditedByUserID, &chapter.CreatedAt, &chapter.UpdatedAt, &chapter.PublishedAt,
		)
		if err != nil {
			log.Printf("Error scanning chapter row: %v", err)
			return nil, fmt.Errorf("failed to scan chapter: %w", err)
		}
		chapters = append(chapters, chapter)
	}

	if err := rows.Err(); err != nil {
		log.Printf("Error iterating over chapter rows: %v", err)
		return nil, fmt.Errorf("failed to iterate chapter rows: %w", err)
	}

	return chapters, nil
}

// Helper function to count words in a string
func countWords(text string) int {
	if text == "" {
		return 0
	}

	// This is a simple word count implementation
	// In a production system, you might want something more sophisticated
	wordCount := 0
	inWord := false

	for _, r := range text {
		isSpace := r == ' ' || r == '\n' || r == '\t' || r == '\r'

		if !isSpace && !inWord {
			inWord = true
			wordCount++
		} else if isSpace {
			inWord = false
		}
	}

	return wordCount
}
