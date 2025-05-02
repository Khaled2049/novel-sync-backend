-- Enable UUID extension if not already enabled
CREATE EXTENSION IF NOT EXISTS "uuid-ossp";
CREATE EXTENSION IF NOT EXISTS "pg_trgm";

-- Enum Types (Keeping most from your definition)
CREATE TYPE user_role AS ENUM ('admin', 'user'); -- Global role, might simplify
CREATE TYPE novel_visibility AS ENUM ('private', 'invite_only', 'public');
CREATE TYPE collaboration_role AS ENUM ('owner', 'editor', 'viewer', 'commenter'); -- Added commenter
CREATE TYPE chapter_status AS ENUM ('draft', 'published', 'archived');
CREATE TYPE note_type AS ENUM ('general', 'character_bio', 'place_description', 'plot_point', 'research', 'world_rule', 'item', 'magic_system', 'species', 'organization'); -- Expanded note types for worldbuilding
CREATE TYPE ai_suggestion_status AS ENUM ('pending', 'generating', 'generated', 'failed', 'accepted', 'rejected', 'edited'); -- Added 'failed'
CREATE TYPE content_source AS ENUM ('user', 'ai');
CREATE TYPE ai_suggestion_type AS ENUM ('continuation', 'dialogue', 'description', 'character_idea', 'place_idea', 'plot_point', 'summary', 'rewrite', 'brainstorm'); -- Expanded types

-- ##################################
-- ######## User Management #########
-- ##################################
CREATE TABLE users (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    -- Link to Firebase Authentication User ID
    firebase_uid VARCHAR(128) UNIQUE NOT NULL,
    -- Email is essential for invites and identification
    email VARCHAR(255) UNIQUE NOT NULL,
    -- No password_hash needed - handled by Firebase Auth
    full_name VARCHAR(255),
    bio TEXT,
    profile_picture_url TEXT,
    -- Global role, primarily for platform administration
    role user_role NOT NULL DEFAULT 'user',
    is_active BOOLEAN NOT NULL DEFAULT true,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Index for finding users by Firebase UID and email
CREATE INDEX idx_users_firebase_uid ON users(firebase_uid);
CREATE INDEX idx_users_email ON users(email);

-- ##################################
-- ######## Core Novel Entities ######
-- ##################################
CREATE TABLE novels (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    title TEXT NOT NULL,
    logline TEXT,
    description TEXT,
    genre VARCHAR(100),
    visibility novel_visibility NOT NULL DEFAULT 'private',
    owner_user_id UUID NOT NULL REFERENCES users(id) ON DELETE RESTRICT,
    cover_image_url TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_novels_owner_user_id ON novels(owner_user_id);
CREATE INDEX idx_novels_title ON novels(title); -- Consider GIN/pg_trgm for search

CREATE TABLE chapters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    -- 'content' stores the current working version
    content TEXT,
    status chapter_status NOT NULL DEFAULT 'draft',
    order_index INTEGER NOT NULL,
    word_count INTEGER NOT NULL DEFAULT 0,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    last_edited_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    published_at TIMESTAMPTZ,
    UNIQUE (novel_id, order_index)
);
CREATE INDEX idx_chapters_novel_id_order ON chapters(novel_id, order_index);
CREATE INDEX idx_chapters_last_edited_by ON chapters(last_edited_by_user_id);

CREATE TABLE chapter_revisions (
    id BIGSERIAL PRIMARY KEY,
    chapter_id UUID NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    content TEXT NOT NULL,
    source content_source NOT NULL DEFAULT 'user',
    edited_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    edited_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    revision_notes TEXT,
    word_count INTEGER NOT NULL DEFAULT 0
);
CREATE INDEX idx_chapter_revisions_chapter_id_edited_at ON chapter_revisions(chapter_id, edited_at DESC); -- Ordering history
CREATE INDEX idx_chapter_revisions_edited_by ON chapter_revisions(edited_by_user_id);

-- ##################################
-- ###### World Building Entities ####
-- ##################################
-- Characters, Places, Notes remain largely the same as your proposal
CREATE TABLE characters (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    backstory TEXT,
    motivations TEXT,
    physical_description TEXT,
    image_url TEXT,
    source content_source NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_characters_novel_id ON characters(novel_id);
CREATE INDEX idx_characters_name_trgm ON characters USING gin (name gin_trgm_ops); -- Requires pg_trgm
CREATE INDEX idx_characters_created_by ON characters(created_by_user_id);

CREATE TABLE places (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    name TEXT NOT NULL,
    description TEXT,
    location_details TEXT,
    atmosphere TEXT,
    image_url TEXT,
    source content_source NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL
);
CREATE INDEX idx_places_novel_id ON places(novel_id);
CREATE INDEX idx_places_name_trgm ON places USING gin (name gin_trgm_ops); -- Requires pg_trgm
CREATE INDEX idx_places_created_by ON places(created_by_user_id);

-- Shared Notes (Worldbuilding, Plot, Research)
CREATE TABLE notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    title TEXT,
    content TEXT NOT NULL,
    note_type note_type NOT NULL DEFAULT 'general',
    source content_source NOT NULL DEFAULT 'user',
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    -- Optional links to specific entities
    linked_chapter_id UUID REFERENCES chapters(id) ON DELETE SET NULL,
    linked_character_id UUID REFERENCES characters(id) ON DELETE SET NULL,
    linked_place_id UUID REFERENCES places(id) ON DELETE SET NULL
-- Add other links as needed
);
-- Indexes remain similar, consider GIN index on content for full-text search
CREATE INDEX idx_notes_novel_id ON notes(novel_id);
CREATE INDEX idx_notes_type ON notes(novel_id, note_type);
CREATE INDEX idx_notes_created_by ON notes(created_by_user_id);
CREATE INDEX idx_notes_linked_chapter ON notes(linked_chapter_id) WHERE linked_chapter_id IS NOT NULL;
CREATE INDEX idx_notes_linked_character ON notes(linked_character_id) WHERE linked_character_id IS NOT NULL;
CREATE INDEX idx_notes_linked_place ON notes(linked_place_id) WHERE linked_place_id IS NOT NULL;
-- CREATE INDEX idx_notes_content_fts ON notes USING gin(to_tsvector('english', content)); -- Optional Full Text Search

-- ##################################
-- #### Relationships & Linking #####
-- ##################################
-- Novel Collaborators remains the same
CREATE TABLE novel_collaborators (
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    role collaboration_role NOT NULL DEFAULT 'viewer',
    invited_at TIMESTAMPTZ DEFAULT NOW(),
    joined_at TIMESTAMPTZ,
    PRIMARY KEY (novel_id, user_id)
);
CREATE INDEX idx_novel_collaborators_user_id ON novel_collaborators(user_id);

-- Chapter Characters & Places remain the same
CREATE TABLE chapter_characters (
    chapter_id UUID NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    character_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    appearance_details TEXT,
    PRIMARY KEY (chapter_id, character_id)
);
CREATE INDEX idx_chapter_characters_character_id ON chapter_characters(character_id);

CREATE TABLE chapter_places (
    chapter_id UUID NOT NULL REFERENCES chapters(id) ON DELETE CASCADE,
    place_id UUID NOT NULL REFERENCES places(id) ON DELETE CASCADE,
    scene_details TEXT,
    PRIMARY KEY (chapter_id, place_id)
);
CREATE INDEX idx_chapter_places_place_id ON chapter_places(place_id);

-- Character Relationships remains the same
CREATE TABLE character_relationships (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    character1_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    character2_id UUID NOT NULL REFERENCES characters(id) ON DELETE CASCADE,
    relationship_type TEXT NOT NULL,
    description TEXT,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    CHECK (character1_id <> character2_id),
    -- Choose directional or non-directional unique constraint
    UNIQUE (novel_id, character1_id, character2_id, relationship_type) -- Directional
    -- UNIQUE (novel_id, LEAST(character1_id, character2_id), GREATEST(character1_id, character2_id), relationship_type) -- Non-directional
);
CREATE INDEX idx_character_relationships_novel_id ON character_relationships(novel_id);
CREATE INDEX idx_character_relationships_char1 ON character_relationships(character1_id);
CREATE INDEX idx_character_relationships_char2 ON character_relationships(character2_id);
CREATE INDEX idx_character_relationships_created_by ON character_relationships(created_by_user_id);


-- ##################################
-- ####### Additional Features ######
-- ##################################
-- Timeline Events and Links remain the same as your proposal
CREATE TABLE timeline_events (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    title TEXT NOT NULL,
    description TEXT,
    event_order INTEGER NOT NULL,
    event_date_text VARCHAR(100),
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    created_by_user_id UUID REFERENCES users(id) ON DELETE SET NULL,
    UNIQUE(novel_id, event_order)
);
CREATE INDEX idx_timeline_events_novel_id_order ON timeline_events(novel_id, event_order);
CREATE INDEX idx_timeline_events_created_by ON timeline_events(created_by_user_id);


-- ##################################
-- ####### Collaboration Features ###
-- ##################################
-- **NEW:** Comments Table (Similar to original proposal, adapted to this schema)
CREATE TABLE comments (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE, -- For context/permissions
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE, -- Author of comment
    -- Target of the comment
    target_chapter_id UUID REFERENCES chapters(id) ON DELETE SET NULL,
    target_character_id UUID REFERENCES characters(id) ON DELETE SET NULL,
    target_place_id UUID REFERENCES places(id) ON DELETE SET NULL,
    target_note_id UUID REFERENCES notes(id) ON DELETE SET NULL,
    target_timeline_event_id UUID REFERENCES timeline_events(id) ON DELETE SET NULL, -- Link to timeline if added
    -- Add target_content_block_id if block-level granularity is added later
    -- Threading
    parent_comment_id UUID REFERENCES comments(id) ON DELETE CASCADE,
    -- Content & Status
    content TEXT NOT NULL,
    status VARCHAR(50) DEFAULT 'open', -- e.g., 'open', 'resolved'
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
    -- Constraint to ensure it links to ONE specific target or is a reply (parent_comment_id is NOT NULL)
    -- CHECK constraint can be complex; might rely on application logic or a trigger
);
-- Indexes for comments
CREATE INDEX idx_comments_novel_id ON comments(novel_id);
CREATE INDEX idx_comments_user_id ON comments(user_id);
CREATE INDEX idx_comments_parent_comment_id ON comments(parent_comment_id) WHERE parent_comment_id IS NOT NULL;
-- Indexes for finding comments by target efficiently
CREATE INDEX idx_comments_target_chapter ON comments(target_chapter_id) WHERE target_chapter_id IS NOT NULL;
CREATE INDEX idx_comments_target_character ON comments(target_character_id) WHERE target_character_id IS NOT NULL;
CREATE INDEX idx_comments_target_place ON comments(target_place_id) WHERE target_place_id IS NOT NULL;
CREATE INDEX idx_comments_target_note ON comments(target_note_id) WHERE target_note_id IS NOT NULL;
CREATE INDEX idx_comments_target_timeline ON comments(target_timeline_event_id) WHERE target_timeline_event_id IS NOT NULL;

-- **NEW:** User Private Notes (Separate from shared novel Notes)
CREATE TABLE user_private_notes (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    -- Link to a novel for context, but access is strictly by user_id
    novel_id UUID REFERENCES novels(id) ON DELETE SET NULL,
    title TEXT,
    content TEXT NOT NULL,
    -- Optional: Could link polymorphically or keep simple
    -- target_type VARCHAR(50), -- e.g., 'novel', 'chapter', 'character'
    -- target_id UUID,
    created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
CREATE INDEX idx_user_private_notes_user_id ON user_private_notes(user_id);
-- Add index on (user_id, novel_id) if filtering by novel context is common
CREATE INDEX idx_user_private_notes_user_novel ON user_private_notes(user_id, novel_id) WHERE novel_id IS NOT NULL;



CREATE TABLE timeline_event_links (
    event_id UUID NOT NULL REFERENCES timeline_events(id) ON DELETE CASCADE,
    entity_type VARCHAR(50) NOT NULL,
    entity_id UUID NOT NULL,
    PRIMARY KEY (event_id, entity_type, entity_id)
);
CREATE INDEX idx_timeline_event_links_entity ON timeline_event_links(entity_type, entity_id);

-- ##################################
-- ######## AI Integration ##########
-- ##################################
-- AI Suggestions table remains the same as your proposal
CREATE TABLE ai_suggestions (
    id UUID PRIMARY KEY DEFAULT uuid_generate_v4(),
    novel_id UUID NOT NULL REFERENCES novels(id) ON DELETE CASCADE,
    user_id UUID NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    suggestion_type ai_suggestion_type NOT NULL,
    status ai_suggestion_status NOT NULL DEFAULT 'pending',
    context_chapter_id UUID REFERENCES chapters(id) ON DELETE SET NULL,
    context_character_id UUID REFERENCES characters(id) ON DELETE SET NULL,
    context_place_id UUID REFERENCES places(id) ON DELETE SET NULL,
    context_note_id UUID REFERENCES notes(id) ON DELETE SET NULL,
    -- Context snippet provided (optional, could be large)
    prompt_context TEXT,
    prompt_instructions TEXT,
    generated_content TEXT,
    model_used VARCHAR(100),
    generation_metadata JSONB,
    user_feedback SMALLINT,
    user_notes TEXT,
    requested_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    generated_at TIMESTAMPTZ,
    updated_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
);
-- Indexes remain the same
CREATE INDEX idx_ai_suggestions_novel_id ON ai_suggestions(novel_id);
CREATE INDEX idx_ai_suggestions_user_id ON ai_suggestions(user_id);
CREATE INDEX idx_ai_suggestions_status ON ai_suggestions(novel_id, status);
CREATE INDEX idx_ai_suggestions_context_chapter ON ai_suggestions(context_chapter_id) WHERE context_chapter_id IS NOT NULL;
CREATE INDEX idx_ai_suggestions_context_character ON ai_suggestions(context_character_id) WHERE context_character_id IS NOT NULL;
CREATE INDEX idx_ai_suggestions_context_place ON ai_suggestions(context_place_id) WHERE context_place_id IS NOT NULL;


-- ##################################
-- ######## Utility Functions #######
-- ##################################
-- Trigger function to automatically update 'updated_at' timestamp
CREATE OR REPLACE FUNCTION update_updated_at_column()
RETURNS TRIGGER AS $$
BEGIN
NEW.updated_at = NOW();
RETURN NEW;
END;
$$ language 'plpgsql';

-- Apply the trigger to all relevant tables
CREATE TRIGGER update_users_updated_at BEFORE UPDATE ON users FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_novels_updated_at BEFORE UPDATE ON novels FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_chapters_updated_at BEFORE UPDATE ON chapters FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_characters_updated_at BEFORE UPDATE ON characters FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_places_updated_at BEFORE UPDATE ON places FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_notes_updated_at BEFORE UPDATE ON notes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_character_relationships_updated_at BEFORE UPDATE ON character_relationships FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_timeline_events_updated_at BEFORE UPDATE ON timeline_events FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_ai_suggestions_updated_at BEFORE UPDATE ON ai_suggestions FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_comments_updated_at BEFORE UPDATE ON comments FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();
CREATE TRIGGER update_user_private_notes_updated_at BEFORE UPDATE ON user_private_notes FOR EACH ROW EXECUTE FUNCTION update_updated_at_column();