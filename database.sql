-- This is the SQL script that will be used to initialize the database schema.
-- We will evaluate you based on how well you design your database.
-- 1. How you design the tables.
-- 2. How you choose the data types and keys.
-- 3. How you name the fields.
-- In this assignment we will use PostgreSQL as the database.

-- This is test table. Remove this table and replace with your own tables. 
CREATE TABLE IF NOT EXISTS estates (
    id VARCHAR(36) NOT NULL,
    width BIGINT NOT NULL,
    length BIGINT NOT NULL,
    count BIGINT NOT NULL DEFAULT 0,
    max BIGINT NOT NULL DEFAULT 0,
    min BIGINT NOT NULL DEFAULT 0,
    drone_distance BIGINT NOT NULL DEFAULT 0,
    median DOUBLE PRECISION NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY(id)
);

CREATE TABLE IF NOT EXISTS estate_trees (
    id VARCHAR(36) NOT NULL,
    estate_id VARCHAR(36) NOT NULL,
    x BIGINT NOT NULL,
    y BIGINT NOT NULL,
    height BIGINT NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    
    PRIMARY KEY(id),
    FOREIGN KEY (estate_id) REFERENCES estates(id)
);

CREATE INDEX IF NOT EXISTS idx_estate_estate_trees ON estate_trees(estate_id);
