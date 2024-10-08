CREATE TABLE IF NOT EXISTS posts (
    id SERIAL PRIMARY KEY,
    title VARCHAR(70) NOT NULL,
    content TEXT NULL,
    tags text[] NOT NULL,
    version INT NOT NULL DEFAULT 0,
    created_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP,
    updated_at TIMESTAMP WITH TIME ZONE NOT NULL DEFAULT CURRENT_TIMESTAMP
);

CREATE TRIGGER posts_moddatetime
BEFORE UPDATE ON posts
FOR EACH ROW
EXECUTE PROCEDURE moddatetime(updated_at);