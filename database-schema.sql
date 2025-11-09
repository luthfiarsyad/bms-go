-- Book Management System Database Schema
-- MySQL Database Structure

-- Create database if not exists
CREATE DATABASE IF NOT EXISTS bms_go CHARACTER SET utf8mb4 COLLATE utf8mb4_unicode_ci;

-- Use the database
USE bms_go;

-- Books table
-- Based on model.Book struct with GORM Model fields
CREATE TABLE IF NOT EXISTS books (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    title VARCHAR(255) NOT NULL,
    author VARCHAR(255) NOT NULL,
    category VARCHAR(255) NOT NULL,
    INDEX idx_books_deleted_at (deleted_at)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Favorites table
-- Based on model.Favorite struct with GORM Model fields
CREATE TABLE IF NOT EXISTS favorites (
    id BIGINT UNSIGNED AUTO_INCREMENT PRIMARY KEY,
    created_at DATETIME(3) NULL,
    updated_at DATETIME(3) NULL,
    deleted_at DATETIME(3) NULL,
    user_id BIGINT UNSIGNED NOT NULL,
    book_id BIGINT UNSIGNED NOT NULL,
    INDEX idx_favorites_deleted_at (deleted_at),
    INDEX idx_favorites_user_id (user_id),
    INDEX idx_favorites_book_id (book_id),
    UNIQUE KEY unique_user_book (user_id, book_id)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4 COLLATE=utf8mb4_unicode_ci;

-- Add foreign key constraint for favorites.book_id referencing books.id
ALTER TABLE favorites 
ADD CONSTRAINT fk_favorites_book_id 
FOREIGN KEY (book_id) REFERENCES books(id) 
ON DELETE CASCADE 
ON UPDATE CASCADE;

-- Insert sample data for testing
-- Sample Books
INSERT INTO books (title, author, category, created_at, updated_at) VALUES
('The Great Gatsby', 'F. Scott Fitzgerald', 'Classic', NOW(), NOW()),
('Harry Potter and the Sorcerer''s Stone', 'J.K. Rowling', 'Fantasy', NOW(), NOW()),
('Dune', 'Frank Herbert', 'Science Fiction', NOW(), NOW()),
('1984', 'George Orwell', 'Dystopian', NOW(), NOW()),
('To Kill a Mockingbird', 'Harper Lee', 'Classic', NOW(), NOW()),
('The Hobbit', 'J.R.R. Tolkien', 'Fantasy', NOW(), NOW()),
('Pride and Prejudice', 'Jane Austen', 'Romance', NOW(), NOW()),
('The Catcher in the Rye', 'J.D. Salinger', 'Classic', NOW(), NOW()),
('Brave New World', 'Aldous Huxley', 'Science Fiction', NOW(), NOW()),
('The Lord of the Rings', 'J.R.R. Tolkien', 'Fantasy', NOW(), NOW());

-- Sample Favorites (assuming user_id = 1 as hardcoded in the application)
INSERT INTO favorites (user_id, book_id, created_at, updated_at) VALUES
(1, 1, NOW(), NOW()),  -- User 1 favorites "The Great Gatsby"
(1, 3, NOW(), NOW()),  -- User 1 favorites "Dune"
(1, 4, NOW(), NOW()),  -- User 1 favorites "1984"
(1, 7, NOW(), NOW());  -- User 1 favorites "Pride and Prejudice"

-- Create a view for books with favorites count (useful for reporting)
CREATE OR REPLACE VIEW books_with_favorite_count AS
SELECT 
    b.id,
    b.title,
    b.author,
    b.category,
    b.created_at,
    b.updated_at,
    COUNT(f.id) as favorite_count
FROM books b
LEFT JOIN favorites f ON b.id = f.book_id AND f.deleted_at IS NULL
WHERE b.deleted_at IS NULL
GROUP BY b.id, b.title, b.author, b.category, b.created_at, b.updated_at;

-- Create a view for user favorites with book details
CREATE OR REPLACE VIEW user_favorites_details AS
SELECT 
    f.id as favorite_id,
    f.user_id,
    f.book_id,
    f.created_at as favorited_at,
    b.title,
    b.author,
    b.category
FROM favorites f
INNER JOIN books b ON f.book_id = b.id
WHERE f.deleted_at IS NULL AND b.deleted_at IS NULL;


-- Create indexes for better performance
CREATE INDEX idx_books_title ON books(title);
CREATE INDEX idx_books_author ON books(author);
CREATE INDEX idx_books_category ON books(category);
CREATE INDEX idx_books_title_author ON books(title, author);

-- Create full-text search index (optional, for advanced search)
-- Note: Requires MyISAM or InnoDB with full-text support
-- ALTER TABLE books ADD FULLTEXT(title, author, category);

-- Grant permissions (adjust as needed for your setup)
-- GRANT ALL PRIVILEGES ON bms_go.* TO 'bms_user'@'localhost' IDENTIFIED BY 'your_password';
-- FLUSH PRIVILEGES;

-- Show table structure
-- DESCRIBE books;
-- DESCRIBE favorites;

-- Sample queries for testing
/*
-- Get all books
SELECT * FROM books WHERE deleted_at IS NULL;

-- Get book by ID
SELECT * FROM books WHERE id = 1 AND deleted_at IS NULL;

-- Search books by title
SELECT * FROM books WHERE title LIKE '%Harry%' AND deleted_at IS NULL;

-- Get books by category
SELECT * FROM books WHERE category = 'Fantasy' AND deleted_at IS NULL;

-- Get all favorites for user 1
SELECT f.*, b.title, b.author, b.category 
FROM favorites f 
INNER JOIN books b ON f.book_id = b.id 
WHERE f.user_id = 1 AND f.deleted_at IS NULL AND b.deleted_at IS NULL;

-- Get books with favorite count
SELECT * FROM books_with_favorite_count;

-- Get user favorites with details
SELECT * FROM user_favorites_details WHERE user_id = 1;
*/