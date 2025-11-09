# Book Management System API - Curl Commands (Refactored)

**Base URL:** http://localhost:8080

## API Response Format

All API responses now follow a consistent format:

```json
{
  "success": true,
  "message": "Operation completed successfully",
  "data": { ... } // or [] for arrays
}
```

Error responses:

```json
{
  "success": false,
  "message": "Error description",
  "error": "Detailed error message"
}
```

## Books API Endpoints

### 1. Get All Books (Basic Search)

```bash
# Get all books without filters
curl -X GET "http://localhost:8080/books" \
  -H "Content-Type: application/json"

# Get all books with search parameter (basic contains search)
curl -X GET "http://localhost:8080/books?search=harry" \
  -H "Content-Type: application/json"

# Get all books with category filter
curl -X GET "http://localhost:8080/books?category=fiction" \
  -H "Content-Type: application/json"

# Get all books with both search and category
curl -X GET "http://localhost:8080/books?search=harry&category=fiction" \
  -H "Content-Type: application/json"
```

### 2. Advanced Search

```bash
# Advanced search with relevance scoring
curl -X GET "http://localhost:8080/books/search?query=harry%20potter&search_type=contains&sort_by=relevance&limit=10" \
  -H "Content-Type: application/json"

# Exact title search
curl -X GET "http://localhost:8080/books/search?query=The%20Great%20Gatsby&search_type=exact&sort_by=title" \
  -H "Content-Type: application/json"

# Search by author with fuzzy matching
curl -X GET "http://localhost:8080/books/search?query=rowling&search_type=fuzzy&sort_by=relevance&limit=5" \
  -H "Content-Type: application/json"

# Search books starting with specific text
curl -X GET "http://localhost:8080/books/search?query=Harry&search_type=starts_with&sort_by=title&sort_order=ASC" \
  -H "Content-Type: application/json"

# Complex search with multiple filters
curl -X GET "http://localhost:8080/books/search?query=fantasy&category=Fantasy&author=j&search_type=contains&sort_by=relevance&limit=20&offset=0" \
  -H "Content-Type: application/json"

# Search by category only
curl -X GET "http://localhost:8080/books/search?category=Science%20Fiction&sort_by=title&sort_order=ASC" \
  -H "Content-Type: application/json"

# Search by author only
curl -X GET "http://localhost:8080/books/search?author=Asimov&sort_by=title&sort_order=ASC" \
  -H "Content-Type: application/json"
```

### 3. Search Suggestions

```bash
# Get search suggestions for "harry"
curl -X GET "http://localhost:8080/books/suggestions?query=harry&limit=5" \
  -H "Content-Type: application/json"

# Get search suggestions for "lord"
curl -X GET "http://localhost:8080/books/suggestions?query=lord&limit=10" \
  -H "Content-Type: application/json"
```

### 4. Get Book by ID

```bash
# Get book with ID 1
curl -X GET "http://localhost:8080/books/1" \
  -H "Content-Type: application/json"

# Get book with ID 5
curl -X GET "http://localhost:8080/books/5" \
  -H "Content-Type: application/json"
```

### 5. Create New Book

```bash
# Create a new book
curl -X POST "http://localhost:8080/books" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Great Gatsby",
    "author": "F. Scott Fitzgerald",
    "category": "Classic"
  }'

# Create another book example
curl -X POST "http://localhost:8080/books" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Harry Potter and the Sorcerer''s Stone",
    "author": "J.K. Rowling",
    "category": "Fantasy"
  }'

# Create a science fiction book
curl -X POST "http://localhost:8080/books" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Dune",
    "author": "Frank Herbert",
    "category": "Science Fiction"
  }'
```

### 6. Update Book

```bash
# Update book with ID 1
curl -X PUT "http://localhost:8080/books/1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "The Great Gatsby (Updated Edition)",
    "author": "F. Scott Fitzgerald",
    "category": "Classic Literature"
  }'

# Update book with ID 2
curl -X PUT "http://localhost:8080/books/2" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "Harry Potter and the Philosopher''s Stone",
    "author": "J.K. Rowling",
    "category": "Fantasy Adventure"
  }'
```

### 7. Delete Book

```bash
# Delete book with ID 1
curl -X DELETE "http://localhost:8080/books/1" \
  -H "Content-Type: application/json"

# Delete book with ID 3
curl -X DELETE "http://localhost:8080/books/3" \
  -H "Content-Type: application/json"
```

## Favorites API Endpoints

### 8. Get All Favorites

```bash
# Get all user's favorite books (user ID is hardcoded to 1 in the backend)
curl -X GET "http://localhost:8080/favorites" \
  -H "Content-Type: application/json"
```

### 9. Get Favorite by ID

```bash
# Get favorite with ID 1
curl -X GET "http://localhost:8080/favorites/1" \
  -H "Content-Type: application/json"

# Get favorite with ID 5
curl -X GET "http://localhost:8080/favorites/5" \
  -H "Content-Type: application/json"
```

### 10. Add Book to Favorites

```bash
# Add book with ID 1 to favorites
curl -X POST "http://localhost:8080/favorites" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 1
  }'

# Add book with ID 2 to favorites
curl -X POST "http://localhost:8080/favorites" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 2
  }'

# Add book with ID 5 to favorites
curl -X POST "http://localhost:8080/favorites" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 5
  }'
```

### 11. Remove Favorite

```bash
# Remove favorite with ID 1
curl -X DELETE "http://localhost:8080/favorites/1" \
  -H "Content-Type: application/json"

# Remove favorite with ID 3
curl -X DELETE "http://localhost:8080/favorites/3" \
  -H "Content-Type: application/json"
```

## Additional Endpoints

### 12. Swagger Documentation

```bash
# Access Swagger UI
curl -X GET "http://localhost:8080/swagger/index.html"

# Access Swagger JSON
curl -X GET "http://localhost:8080/swagger/doc.json"
```

## Usage Notes

1. **Server Setup**: Make sure the server is running on `http://localhost:8080` before executing these commands.

2. **User Context**: The favorites API currently uses a hardcoded user ID (1) as seen in the handler code.

3. **Response Codes**:

   - `200 OK`: Successful GET, PUT, DELETE requests
   - `201 Created`: Successful POST requests
   - `400 Bad Request`: Invalid request body or validation errors
   - `404 Not Found`: Resource not found
   - `409 Conflict`: Duplicate resource (e.g., book already in favorites)
   - `500 Internal Server Error`: Server error

4. **JSON Validation**: All POST and PUT requests require valid JSON with required fields:

   - For books: `title`, `author`, `category` are required (max 255 characters each)
   - For favorites: `book_id` is required and must be a positive integer

5. **Query Parameters**:

   **Basic Search (GET /books):**

   - `search`: Search in book titles/authors (trimmed whitespace)
   - `category`: Filter by book category (trimmed whitespace)

   **Advanced Search (GET /books/search):**

   - `query`: Search query string
   - `category`: Filter by category
   - `author`: Filter by author (partial match)
   - `search_type`: Search strategy - `exact`, `starts_with`, `contains`, `fuzzy`
   - `sort_by`: Sort field - `title`, `author`, `category`, `created_at`, `relevance`
   - `sort_order`: Sort order - `ASC`, `DESC`
   - `limit`: Maximum results (max 100, default 20)
   - `offset`: Pagination offset (default 0)

   **Search Suggestions (GET /books/suggestions):**

   - `query`: Search query for suggestions (required)
   - `limit`: Maximum suggestions (max 20, default 10)

6. **Search Features**:

   - **Relevance Scoring**: Results are ranked by relevance (exact matches first, then starts with, then contains)
   - **Multiple Search Strategies**:
     - `exact`: Exact title or author match
     - `starts_with`: Title/author starts with query
     - `contains`: Title/author contains query (default)
     - `fuzzy`: Advanced fuzzy matching with variations and misspellings
   - **Fuzzy Search**: Automatically handles common misspellings and variations
   - **Search Suggestions**: Provides autocomplete suggestions based on existing titles/authors
   - **Multi-field Search**: Searches across title, author, and category fields

7. **Validation Features**:

   - Duplicate book title checking
   - Duplicate favorite checking
   - Input validation and sanitization
   - Proper error messages
   - Search parameter validation

8. **New Features**:
   - Consistent API response format
   - Better error handling
   - Input validation
   - Complete CRUD operations for favorites
   - Soft delete support for books
   - Advanced search with relevance scoring
   - Search suggestions/autocomplete
   - Multiple search strategies
   - Pagination support

## Example Workflow

```bash
# 1. Create a new book
curl -X POST "http://localhost:8080/books" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "1984",
    "author": "George Orwell",
    "category": "Dystopian"
  }'

# 2. Get all books to see the new book
curl -X GET "http://localhost:8080/books" \
  -H "Content-Type: application/json"

# 3. Add the book to favorites (assuming it got ID 1)
curl -X POST "http://localhost:8080/favorites" \
  -H "Content-Type: application/json" \
  -d '{
    "book_id": 1
  }'

# 4. Get all favorites to verify
curl -X GET "http://localhost:8080/favorites" \
  -H "Content-Type: application/json"

# 5. Get specific favorite by ID
curl -X GET "http://localhost:8080/favorites/1" \
  -H "Content-Type: application/json"

# 6. Update the book
curl -X PUT "http://localhost:8080/books/1" \
  -H "Content-Type: application/json" \
  -d '{
    "title": "1984: Special Edition",
    "author": "George Orwell",
    "category": "Classic Dystopian"
  }'

# 7. Remove from favorites
curl -X DELETE "http://localhost:8080/favorites/1" \
  -H "Content-Type: application/json"

# 8. Delete the book
curl -X DELETE "http://localhost:8080/books/1" \
  -H "Content-Type: application/json"
```
