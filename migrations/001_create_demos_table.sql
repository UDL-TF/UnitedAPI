-- Migration: Create demos table
-- Date: 2026-02-28
-- Description: Creates the demos table for storing demo file metadata

CREATE TABLE IF NOT EXISTS demos (
    id SERIAL PRIMARY KEY,
    is_tournament BOOLEAN NOT NULL DEFAULT FALSE,
    tournament_id INTEGER,
    match_id INTEGER,
    round_id INTEGER,
    raw_demo_name VARCHAR(255) NOT NULL,
    file_size BIGINT NOT NULL,
    blue_team_name VARCHAR(100) NOT NULL,
    red_team_name VARCHAR(100) NOT NULL,
    object_name VARCHAR(500) NOT NULL UNIQUE,
    content_type VARCHAR(100),
    is_compressed BOOLEAN DEFAULT TRUE,
    compressed_size BIGINT,
    uploaded_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    
    -- Constraints
    CONSTRAINT chk_tournament_data CHECK (
        (is_tournament = FALSE) OR 
        (is_tournament = TRUE AND tournament_id IS NOT NULL AND match_id IS NOT NULL AND round_id IS NOT NULL)
    )
);

-- Create indexes for better query performance
CREATE INDEX IF NOT EXISTS idx_demos_match_round ON demos (match_id, round_id);
CREATE INDEX IF NOT EXISTS idx_demos_tournament ON demos (tournament_id);
CREATE INDEX IF NOT EXISTS idx_demos_tournament_match_round ON demos (tournament_id, match_id, round_id);
CREATE INDEX IF NOT EXISTS idx_demos_uploaded_at ON demos (uploaded_at);
CREATE INDEX IF NOT EXISTS idx_demos_object_name ON demos (object_name);