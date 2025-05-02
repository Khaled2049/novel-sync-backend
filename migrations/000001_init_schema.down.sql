-- File: migrations/000001_init_schema.down.sql

-- Drop triggers (reverse order of creation doesn't strictly matter here, but good practice)
DROP TRIGGER IF EXISTS update_user_private_notes_updated_at ON user_private_notes;
DROP TRIGGER IF EXISTS update_comments_updated_at ON comments;
DROP TRIGGER IF EXISTS update_ai_suggestions_updated_at ON ai_suggestions;
DROP TRIGGER IF EXISTS update_timeline_events_updated_at ON timeline_events;
DROP TRIGGER IF EXISTS update_character_relationships_updated_at ON character_relationships;
DROP TRIGGER IF EXISTS update_notes_updated_at ON notes;
DROP TRIGGER IF EXISTS update_places_updated_at ON places;
DROP TRIGGER IF EXISTS update_characters_updated_at ON characters;
DROP TRIGGER IF EXISTS update_chapters_updated_at ON chapters;
DROP TRIGGER IF EXISTS update_novels_updated_at ON novels;
DROP TRIGGER IF EXISTS update_users_updated_at ON users;

-- Drop the function
DROP FUNCTION IF EXISTS update_updated_at_column();

-- Drop tables in reverse order of creation (or based on FK constraints)
-- Consider dependencies: Drop linking tables first, then main tables.
DROP TABLE IF EXISTS timeline_event_links;
DROP TABLE IF EXISTS ai_suggestions;
DROP TABLE IF EXISTS user_private_notes;
DROP TABLE IF EXISTS comments;
DROP TABLE IF EXISTS character_relationships;
DROP TABLE IF EXISTS chapter_places;
DROP TABLE IF EXISTS chapter_characters;
DROP TABLE IF EXISTS novel_collaborators;
DROP TABLE IF EXISTS notes;
DROP TABLE IF EXISTS places;
DROP TABLE IF EXISTS characters;
DROP TABLE IF EXISTS chapter_revisions;
DROP TABLE IF EXISTS chapters;
DROP TABLE IF EXISTS timeline_events; -- Moved here as comments/notes might reference it
DROP TABLE IF EXISTS novels;
DROP TABLE IF EXISTS users;

-- Drop Enum types (reverse order of creation recommended)
DROP TYPE IF EXISTS ai_suggestion_type;
DROP TYPE IF EXISTS content_source;
DROP TYPE IF EXISTS ai_suggestion_status;
DROP TYPE IF EXISTS note_type;
DROP TYPE IF EXISTS chapter_status;
DROP TYPE IF EXISTS collaboration_role;
DROP TYPE IF EXISTS novel_visibility;
DROP TYPE IF EXISTS user_role;

-- Drop Extensions (if they are exclusively used by this schema and safe to drop)
-- Be cautious dropping extensions if other schemas might use them. Often omitted from 'down'.
-- DROP EXTENSION IF EXISTS "pg_trgm";
-- DROP EXTENSION IF EXISTS "uuid-ossp";