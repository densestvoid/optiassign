-- OptiAssign Database Schema
-- Initial migration for core entities

-- Users table for Google SSO authentication
CREATE TABLE users (
    id SERIAL PRIMARY KEY,
    google_id VARCHAR(255) UNIQUE NOT NULL,
    email VARCHAR(255) UNIQUE NOT NULL,
    name VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Groups table for assignment containers
CREATE TABLE groups (
    id SERIAL PRIMARY KEY,
    owner_user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    status VARCHAR(50) DEFAULT 'draft' CHECK (status IN ('draft', 'prioritizing', 'completed')),
    rule VARCHAR(50) DEFAULT 'equal' CHECK (rule IN ('equal', 'exhaust')),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Items table for assignment items
CREATE TABLE items (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    name VARCHAR(255) NOT NULL,
    description TEXT,
    is_assigned BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Participants table linking users to groups
CREATE TABLE participants (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    user_id INTEGER NOT NULL REFERENCES users(id) ON DELETE CASCADE,
    token VARCHAR(255) UNIQUE NOT NULL,
    has_submitted BOOLEAN DEFAULT FALSE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, user_id)
);

-- Prioritization table for user rankings
CREATE TABLE prioritizations (
    id SERIAL PRIMARY KEY,
    participant_id INTEGER NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    rank INTEGER NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(participant_id, item_id),
    UNIQUE(participant_id, rank)
);

-- Assignments table for final results
CREATE TABLE assignments (
    id SERIAL PRIMARY KEY,
    group_id INTEGER NOT NULL REFERENCES groups(id) ON DELETE CASCADE,
    participant_id INTEGER NOT NULL REFERENCES participants(id) ON DELETE CASCADE,
    item_id INTEGER NOT NULL REFERENCES items(id) ON DELETE CASCADE,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP,
    UNIQUE(group_id, item_id)
);

-- Indexes for performance
CREATE INDEX idx_groups_owner ON groups(owner_user_id);
CREATE INDEX idx_items_group ON items(group_id);
CREATE INDEX idx_participants_group ON participants(group_id);
CREATE INDEX idx_participants_user ON participants(user_id);
CREATE INDEX idx_prioritizations_participant ON prioritizations(participant_id);
CREATE INDEX idx_prioritizations_item ON prioritizations(item_id);
CREATE INDEX idx_assignments_group ON assignments(group_id);
CREATE INDEX idx_assignments_participant ON assignments(participant_id);