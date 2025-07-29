-- Step 1: Add isbn column
ALTER TABLE books ADD COLUMN IF NOT EXISTS isbn TEXT;

-- Step 2: Make isbn UNIQUE (optional: defer until data is clean)
ALTER TABLE books ADD CONSTRAINT books_isbn_unique UNIQUE (isbn);