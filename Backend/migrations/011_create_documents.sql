-- +migrate Up
-- Create documents table for RAG system
CREATE TABLE IF NOT EXISTS documents (
    id SERIAL PRIMARY KEY,
    content TEXT NOT NULL,
    embedding vector(1536),
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP
);

-- Create index for vector similarity search (HNSW works with any number of rows)
CREATE INDEX ON documents USING hnsw (embedding vector_cosine_ops);

-- +migrate Down
DROP INDEX IF EXISTS documents_embedding_idx;
DROP TABLE IF EXISTS documents;
