-- Drop the trigger if it exists
DROP TRIGGER IF EXISTS set_updated_at ON todos;
-- Drop the function if it exists
DROP FUNCTION IF EXISTS update_modified_column();
-- Drop the todos table
DROP TABLE IF EXISTS todos;