package repository

import (
	"bms-go/internal/model"
	"strings"

	"gorm.io/gorm"
)

type BookRepository struct {
	db *gorm.DB
}

func NewBookRepository(db *gorm.DB) *BookRepository {
	return &BookRepository{db: db}
}

func (r *BookRepository) FindAll(search, category string) ([]model.Book, error) {
	var books []model.Book
	query := r.db.Where("deleted_at IS NULL")

	if search != "" {
		// Enhanced search with multiple strategies
		search = strings.TrimSpace(search)
		
		// Strategy 1: Exact title match (highest priority)
		// Strategy 2: Title starts with search term
		// Strategy 3: Title contains search term
		// Strategy 4: Author contains search term
		// Strategy 5: Category contains search term
		
		query = query.Where(`
			(title = ?) OR
			(title LIKE ?) OR
			(title LIKE ?) OR
			(author LIKE ?) OR
			(category LIKE ?)
		`,
			search,                           // Exact match
			search+"%",                        // Starts with
			"%"+search+"%",                     // Contains
			"%"+search+"%",                     // Author contains
			"%"+search+"%")                     // Category contains
	}

	if category != "" {
		query = query.Where("category = ?", strings.TrimSpace(category))
	}

	// Order by relevance for search results
	if search != "" {
		// Use raw SQL for complex ordering with parameters
		query = query.Raw(`
			SELECT * FROM books
			WHERE deleted_at IS NULL AND (
				(title = ?) OR
				(title LIKE ?) OR
				(title LIKE ?) OR
				(author LIKE ?) OR
				(category LIKE ?)
			)
			ORDER BY
				CASE
					WHEN title = ? THEN 1
					WHEN title LIKE ? THEN 2
					WHEN title LIKE ? THEN 3
					WHEN author LIKE ? THEN 4
					ELSE 5
				END,
				title ASC
		`, search, search+"%", "%"+search+"%", "%"+search+"%", "%"+search+"%",
			search, search+"%", "%"+search+"%", "%"+search+"%")
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// AdvancedSearch implements sophisticated search with multiple criteria
func (r *BookRepository) AdvancedSearch(params AdvancedSearchParams) ([]model.Book, error) {
	var books []model.Book
	query := r.db.Where("deleted_at IS NULL")

	// Apply search term with advanced matching
	if params.Query != "" {
		queryTerm := strings.TrimSpace(params.Query)
		
		// Build search conditions based on search type
		switch params.SearchType {
		case "exact":
			query = query.Where("title = ? OR author = ?", queryTerm, queryTerm)
		case "starts_with":
			query = query.Where("title LIKE ? OR author LIKE ?", queryTerm+"%", queryTerm+"%")
		case "fuzzy":
			// Implement fuzzy search using multiple LIKE patterns
			fuzzyPatterns := r.generateFuzzyPatterns(queryTerm)
			searchConditions := []string{}
			searchArgs := []interface{}{}
			
			for _, pattern := range fuzzyPatterns {
				searchConditions = append(searchConditions, "title LIKE ?")
				searchArgs = append(searchArgs, pattern)
				searchConditions = append(searchConditions, "author LIKE ?")
				searchArgs = append(searchArgs, pattern)
			}
			
			query = query.Where(strings.Join(searchConditions, " OR "), searchArgs...)
		default: // "contains"
			query = query.Where(`
				(title LIKE ?) OR
				(author LIKE ?) OR
				(category LIKE ?)
			`, "%"+queryTerm+"%", "%"+queryTerm+"%", "%"+queryTerm+"%")
		}
	}

	// Apply category filter
	if params.Category != "" {
		query = query.Where("category = ?", strings.TrimSpace(params.Category))
	}

	// Apply author filter
	if params.Author != "" {
		query = query.Where("author LIKE ?", "%"+strings.TrimSpace(params.Author)+"%")
	}

	// Apply sorting
	switch params.SortBy {
	case "title":
		query = query.Order("title " + params.SortOrder)
	case "author":
		query = query.Order("author " + params.SortOrder)
	case "category":
		query = query.Order("category " + params.SortOrder)
	case "created_at":
		query = query.Order("created_at " + params.SortOrder)
	case "relevance":
		if params.Query != "" {
			queryTerm := strings.TrimSpace(params.Query)
			// Use raw SQL for complex ordering with parameters
			query = query.Raw(`
				SELECT * FROM books
				WHERE deleted_at IS NULL
				ORDER BY
					CASE
						WHEN title = ? THEN 1
						WHEN title LIKE ? THEN 2
						WHEN title LIKE ? THEN 3
						WHEN author LIKE ? THEN 4
						ELSE 5
					END,
					title ASC
			`, queryTerm, queryTerm+"%", "%"+queryTerm+"%", "%"+queryTerm+"%")
		} else {
			query = query.Order("title ASC")
		}
	default:
		query = query.Order("title ASC")
	}

	// Apply pagination
	if params.Limit > 0 {
		query = query.Limit(params.Limit)
	}
	if params.Offset > 0 {
		query = query.Offset(params.Offset)
	}

	if err := query.Find(&books).Error; err != nil {
		return nil, err
	}
	return books, nil
}

// generateFuzzyPatterns creates patterns for fuzzy search
func (r *BookRepository) generateFuzzyPatterns(term string) []string {
	patterns := []string{}
	
	// Original pattern
	patterns = append(patterns, "%"+term+"%")
	
	// Split term into words and search for individual words
	words := strings.Fields(term)
	for _, word := range words {
		if len(word) > 2 { // Only consider words longer than 2 characters
			patterns = append(patterns, "%"+word+"%")
		}
	}
	
	// Common misspellings and variations (can be extended)
	variations := r.generateVariations(term)
	for _, variation := range variations {
		patterns = append(patterns, "%"+variation+"%")
	}
	
	return patterns
}

// generateVariations creates common variations of search terms
func (r *BookRepository) generateVariations(term string) []string {
	variations := []string{}
	term = strings.ToLower(term)
	
	// Common substitutions
	substitutions := map[string][]string{
		"har": {"harr", "harry"},
		"pot": {"pott", "potter"},
		"lord": {"lords"},
		"ring": {"rings", "ring"},
		"game": {"gaming"},
		"throne": {"thrones"},
	}
	
	for key, subs := range substitutions {
		if strings.Contains(term, key) {
			for _, sub := range subs {
				variations = append(variations, strings.Replace(term, key, sub, -1))
			}
		}
	}
	
	return variations
}

// GetSearchSuggestions provides search suggestions based on existing books
func (r *BookRepository) GetSearchSuggestions(query string, limit int) ([]string, error) {
	var suggestions []string
	
	if query == "" {
		return suggestions, nil
	}
	
	query = strings.TrimSpace(query)
	
	// Get unique titles and authors that match the query
	var results []struct {
		Suggestion string
	}
	
	err := r.db.Raw(`
		SELECT DISTINCT title as suggestion FROM books
		WHERE deleted_at IS NULL AND (title LIKE ? OR author LIKE ?)
		UNION
		SELECT DISTINCT author as suggestion FROM books
		WHERE deleted_at IS NULL AND (title LIKE ? OR author LIKE ?)
		LIMIT ?
	`, "%"+query+"%", "%"+query+"%", "%"+query+"%", "%"+query+"%", limit).Scan(&results).Error
	
	if err != nil {
		return nil, err
	}
	
	for _, result := range results {
		suggestions = append(suggestions, result.Suggestion)
	}
	
	return suggestions, nil
}

// AdvancedSearchParams represents parameters for advanced search
type AdvancedSearchParams struct {
	Query      string `json:"query"`
	Category   string `json:"category"`
	Author     string `json:"author"`
	SearchType string `json:"search_type"` // exact, starts_with, contains, fuzzy
	SortBy     string `json:"sort_by"`     // title, author, category, created_at, relevance
	SortOrder  string `json:"sort_order"`  // ASC, DESC
	Limit      int    `json:"limit"`
	Offset     int    `json:"offset"`
}

func (r *BookRepository) FindByID(id uint) (*model.Book, error) {
	var book model.Book
	if err := r.db.Where("deleted_at IS NULL").First(&book, id).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

// Exists checks if a book exists by ID
func (r *BookRepository) Exists(id uint) (bool, error) {
	var count int64
	err := r.db.Model(&model.Book{}).Where("id = ? AND deleted_at IS NULL", id).Count(&count).Error
	if err != nil {
		return false, err
	}
	return count > 0, nil
}

// FindByTitle checks if a book exists by title (for duplicate checking)
func (r *BookRepository) FindByTitle(title string) (*model.Book, error) {
	var book model.Book
	if err := r.db.Where("title = ? AND deleted_at IS NULL", title).First(&book).Error; err != nil {
		return nil, err
	}
	return &book, nil
}

func (r *BookRepository) Create(book *model.Book) error {
	return r.db.Create(book).Error
}

func (r *BookRepository) Update(book *model.Book) error {
	return r.db.Model(&model.Book{}).Where("id = ? AND deleted_at IS NULL", book.ID).Updates(book).Error
}

func (r *BookRepository) Delete(id uint) error {
	return r.db.Delete(&model.Book{}, id).Error
}

func (r *BookRepository) HardDelete(id uint) error {
	return r.db.Unscoped().Delete(&model.Book{}, id).Error
}

func (r *BookRepository) GetCount() (int64, error) {
	var count int64
	err := r.db.Model(&model.Book{}).Where("deleted_at IS NULL").Count(&count).Error
	return count, err
}