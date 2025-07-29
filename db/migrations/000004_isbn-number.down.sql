-- Step 1: Drop the UNIQUE constraint on isbn
ALTER TABLE books DROP CONSTRAINT IF EXISTS books_isbn_unique;

-- Step 2: Drop the isbn column
ALTER TABLE books DROP COLUMN IF EXISTS isbn;